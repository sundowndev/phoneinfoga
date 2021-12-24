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
	allData := make(map[string]interface{}, 0)
	errors := make(map[string]error, 0)

	for name, s := range r.scanners {
		data, err := s.Scan(n)
		if err != nil {
			errors[name] = err
			continue
		}
		allData[name] = data
	}
	return allData, errors
}
