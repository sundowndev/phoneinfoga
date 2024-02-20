package remote

import (
	"errors"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

const Numverify = "numverify"

type numverifyScanner struct {
	client suppliers.NumverifySupplierInterface
}

type NumverifyScannerResponse struct {
	Valid               bool   `json:"valid" console:"Valid"`
	Number              string `json:"number" console:"Number,omitempty"`
	LocalFormat         string `json:"local_format" console:"Local format,omitempty"`
	InternationalFormat string `json:"international_format" console:"International format,omitempty"`
	CountryPrefix       string `json:"country_prefix" console:"Country prefix,omitempty"`
	CountryCode         string `json:"country_code" console:"Country code,omitempty"`
	CountryName         string `json:"country_name" console:"Country name,omitempty"`
	Location            string `json:"location" console:"Location,omitempty"`
	Carrier             string `json:"carrier" console:"Carrier,omitempty"`
	LineType            string `json:"line_type" console:"Line type,omitempty"`
}

func NewNumverifyScanner(s suppliers.NumverifySupplierInterface) Scanner {
	return &numverifyScanner{client: s}
}

func (s *numverifyScanner) Name() string {
	return Numverify
}

func (s *numverifyScanner) Description() string {
	return "Request info about a given phone number through the Numverify API."
}

func (s *numverifyScanner) DryRun(_ number.Number, opts ScannerOptions) error {
	if opts.GetStringEnv("NUMVERIFY_API_KEY") != "" {
		return nil
	}
	return errors.New("API key is not defined")
}

func (s *numverifyScanner) Run(n number.Number, opts ScannerOptions) (interface{}, error) {
	apiKey := opts.GetStringEnv("NUMVERIFY_API_KEY")

	res, err := s.client.Request().SetApiKey(apiKey).ValidateNumber(n.International)
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
