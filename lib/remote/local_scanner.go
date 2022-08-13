package remote

import (
	"github.com/nyaruka/phonenumbers"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
)

const Local = "local"

type localScanner struct{}

type LocalScannerResponse struct {
	RawLocal      string `json:"raw_local,omitempty" console:"Raw local,omitempty"`
	Local         string `json:"local,omitempty" console:"Local,omitempty"`
	E164          string `json:"e164,omitempty" console:"E164,omitempty"`
	International string `json:"international,omitempty" console:"International,omitempty"`
	CountryCode   int32  `json:"country_code,omitempty" console:"Country code,omitempty"`
	Country       string `json:"country,omitempty" console:"Country,omitempty"`
	Carrier       string `json:"carrier,omitempty" console:"Carrier,omitempty"`
	Valid         bool   `json:"valid" console:"Is valid"`
}

func NewLocalScanner() *localScanner {
	return &localScanner{}
}

func (s *localScanner) Name() string {
	return Local
}

func (s *localScanner) ShouldRun(_ number.Number) bool {
	return true
}

func (s *localScanner) Scan(n number.Number) (interface{}, error) {
	num, err := phonenumbers.Parse(n.E164, n.Country)
	if err != nil {
		return nil, err
	}

	data := LocalScannerResponse{
		Valid:         phonenumbers.IsValidNumber(num),
		RawLocal:      n.RawLocal,
		Local:         n.Local,
		E164:          n.E164,
		International: n.International,
		CountryCode:   n.CountryCode,
		Country:       n.Country,
		Carrier:       n.Carrier,
	}
	return data, nil
}
