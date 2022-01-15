package remote

import (
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
)

type Library struct {
	scanners map[string]Scanner
}

func NewLibrary() *Library {
	return &Library{
		scanners: map[string]Scanner{},
	}
}

func (r *Library) AddScanner(s Scanner) {
	if !s.ShouldRun() {
		return
	}
	r.scanners[s.Identifier()] = s
}

func (r *Library) Scan(n *number.Number) (map[string]interface{}, map[string]error) {
	allData := make(map[string]interface{})
	errors := make(map[string]error)

	for name, s := range r.scanners {
		data, err := s.Scan(n)
		if err != nil {
			errors[name] = err
			continue
		}
		if data != nil {
			allData[name] = data
		}
	}
	return allData, errors
}
