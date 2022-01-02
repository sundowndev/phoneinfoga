package remote

import (
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

const Numverify = "numverify"

type numverifyScanner struct {
	client suppliers.NumverifySupplierInterface
}

type NumverifyScannerResponse struct {
	Valid               bool   `json:"valid" console:"Valid"`
	Number              string `json:"number" console:"Number"`
	LocalFormat         string `json:"local_format" console:"Local format"`
	InternationalFormat string `json:"international_format" console:"International format"`
	CountryPrefix       string `json:"country_prefix" console:"Country prefix"`
	CountryCode         string `json:"country_code" console:"Country code"`
	CountryName         string `json:"country_name" console:"Country name"`
	Location            string `json:"location" console:"Location"`
	Carrier             string `json:"carrier" console:"Carrier"`
	LineType            string `json:"line_type" console:"Line type"`
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

func (s *numverifyScanner) Scan(n *number.Number) (interface{}, error) {
	res, err := s.client.Validate(n.International)
	if err != nil {
		return nil, err
	}

	data := NumverifyScannerResponse{
		Valid:               res.Valid,
		Number:              res.Number,
		LocalFormat:         res.LocalFormat,
		InternationalFormat: res.InternationalFormat,
		CountryPrefix:       res.CountryPrefix,
		CountryCode:         res.CountryCode,
		CountryName:         res.CountryName,
		Location:            res.Location,
		Carrier:             res.Carrier,
		LineType:            res.LineType,
	}

	return data, nil
}
