package remote

import (
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

const Numverify = "numverify"

type numverifyScanner struct {
	client suppliers.NumverifySupplierInterface
}

func NewNumverifyScanner(s suppliers.NumverifySupplierInterface) *numverifyScanner {
	return &numverifyScanner{client: s}
}

func (s *numverifyScanner) Identifier() string {
	return Numverify
}

func (s *numverifyScanner) ShouldRun() bool {
	return s.client.IsAvailable()
}

func (s *numverifyScanner) Scan(n *number.Number) (ScannerResult, error) {
	data := make(ScannerResult, 0)

	res, err := s.client.ScanNumber(n.International)
	if err != nil {
		return data, err
	}

	data["valid"] = res.Valid
	data["number"] = res.Number
	data["local_format"] = res.LocalFormat
	data["international_format"] = res.InternationalFormat
	data["country_prefix"] = res.CountryPrefix
	data["country_code"] = res.CountryCode
	data["country_name"] = res.CountryName
	data["location"] = res.Location
	data["carrier"] = res.Carrier
	data["line_type"] = res.LineType

	return data, nil
}
