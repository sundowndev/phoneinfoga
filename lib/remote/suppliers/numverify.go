package suppliers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type NumverifySupplierInterface interface {
	IsAvailable() bool
	Validate(string) (*NumverifyScannerResponse, error)
}

type numverifyError struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

// NumverifyScannerResponse REST API response
type NumverifyScannerResponse struct {
	Valid               bool           `json:"valid"`
	Number              string         `json:"number"`
	LocalFormat         string         `json:"local_format"`
	InternationalFormat string         `json:"international_format"`
	CountryPrefix       string         `json:"country_prefix"`
	CountryCode         string         `json:"country_code"`
	CountryName         string         `json:"country_name"`
	Location            string         `json:"location"`
	Carrier             string         `json:"carrier"`
	LineType            string         `json:"line_type"`
	Error               numverifyError `json:"error"`
}

type NumverifySupplier struct {
	ApiKey string
}

func NewNumverifySupplier() *NumverifySupplier {
	return &NumverifySupplier{
		ApiKey: os.Getenv("NUMVERIFY_API_KEY"),
	}
}

func (s *NumverifySupplier) IsAvailable() bool {
	return s.ApiKey != ""
}

func (s *NumverifySupplier) Validate(internationalNumber string) (res *NumverifyScannerResponse, err error) {
	// Build the request
	response, err := http.Get(fmt.Sprintf("http://apilayer.net/api/validate?access_key=%s&number=%s", s.ApiKey, internationalNumber))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Fill the response with the data from the JSON
	var result NumverifyScannerResponse

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Error.Info) > 0 {
		return nil, fmt.Errorf("%s", result.Error.Info)
	}

	res = &NumverifyScannerResponse{
		Valid:               result.Valid,
		Number:              result.Number,
		LocalFormat:         result.LocalFormat,
		InternationalFormat: result.InternationalFormat,
		CountryPrefix:       result.CountryPrefix,
		CountryCode:         result.CountryCode,
		CountryName:         result.CountryName,
		Location:            result.Location,
		Carrier:             result.Carrier,
		LineType:            result.LineType,
	}

	return res, nil
}
