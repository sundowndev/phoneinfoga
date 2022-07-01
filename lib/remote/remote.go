package remote

import (
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"sync"
)

type Library struct {
	m        *sync.RWMutex
	scanners []Scanner
	results  map[string]interface{}
	errors   map[string]error
}

type Scanner interface {
	Scan(*number.Number) (interface{}, error)
	ShouldRun() bool
	Identifier() string
}

func NewLibrary() *Library {
	return &Library{
		m:        &sync.RWMutex{},
		scanners: []Scanner{},
		results:  map[string]interface{}{},
		errors:   map[string]error{},
	}
}

func (r *Library) AddScanner(s Scanner) {
	if !s.ShouldRun() {
		return
	}
	r.scanners = append(r.scanners, s)
}

func (r *Library) AddResult(k string, v interface{}) {
	r.m.Lock()
	defer r.m.Unlock()
	r.results[k] = v
}

func (r *Library) AddError(k string, err error) {
	r.m.Lock()
	defer r.m.Unlock()
	r.errors[k] = err
}

func (r *Library) Scan(n *number.Number) (map[string]interface{}, map[string]error) {
	var wg sync.WaitGroup

	for _, s := range r.scanners {
		wg.Add(1)

		go func(s Scanner) {
			defer wg.Done()
			data, err := s.Scan(n)
			if err != nil {
				r.AddError(s.Identifier(), err)
				return
			}
			if data != nil {
				r.AddResult(s.Identifier(), data)
			}
		}(s)
	}

	wg.Wait()

	return r.results, r.errors
}
