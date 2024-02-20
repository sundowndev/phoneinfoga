package remote

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"sync"
)

var mu sync.Locker = &sync.RWMutex{}
var plugins []Scanner

type Library struct {
	m        *sync.RWMutex
	scanners []Scanner
	results  map[string]interface{}
	errors   map[string]error
	filter   filter.Filter
}

func NewLibrary(filterEngine filter.Filter) *Library {
	return &Library{
		m:        &sync.RWMutex{},
		scanners: []Scanner{},
		results:  map[string]interface{}{},
		errors:   map[string]error{},
		filter:   filterEngine,
	}
}

func (r *Library) LoadPlugins() {
	for _, s := range plugins {
		r.AddScanner(s)
	}
}

func (r *Library) AddScanner(s Scanner) {
	if r.filter.Match(s.Name()) {
		logrus.WithField("scanner", s.Name()).Debug("Scanner was ignored by filter")
		return
	}
	r.scanners = append(r.scanners, s)
}

func (r *Library) addResult(k string, v interface{}) {
	r.m.Lock()
	defer r.m.Unlock()
	r.results[k] = v
}

func (r *Library) addError(k string, err error) {
	r.m.Lock()
	defer r.m.Unlock()
	r.errors[k] = err
}

func (r *Library) Scan(n *number.Number, opts ScannerOptions) (map[string]interface{}, map[string]error) {
	var wg sync.WaitGroup

	for _, s := range r.scanners {
		wg.Add(1)
		go func(s Scanner) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					logrus.WithField("scanner", s.Name()).WithField("error", err).Debug("Scanner panicked")
					r.addError(s.Name(), errors.New("panic occurred while running scan, see debug logs"))
				}
			}()

			if err := s.DryRun(*n, opts); err != nil {
				logrus.
					WithField("scanner", s.Name()).
					WithField("reason", err.Error()).
					Debug("Scanner was ignored because it should not run")
				return
			}

			data, err := s.Run(*n, opts)
			if err != nil {
				r.addError(s.Name(), err)
				return
			}
			if data != nil {
				r.addResult(s.Name(), data)
			}
		}(s)
	}

	wg.Wait()

	return r.results, r.errors
}

func (r *Library) GetAllScanners() []Scanner {
	return r.scanners
}

func (r *Library) GetScanner(name string) Scanner {
	r.m.RLock()
	defer r.m.RUnlock()
	for _, s := range r.scanners {
		if s.Name() == name {
			return s
		}
	}
	return nil
}

func RegisterPlugin(s Scanner) {
	mu.Lock()
	defer mu.Unlock()
	plugins = append(plugins, s)
}
