package suppliers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type NumverifySupplierInterface interface {
	IsAvailable() bool
	Validate(string) (*NumverifyValidateResponse, error)
}

type numverifyError struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

// NumverifyValidateResponse REST API response
type NumverifyValidateResponse struct {
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
	ApiKey    string
	EnableSSL string
}

func NewNumverifySupplier() *NumverifySupplier {
	return &NumverifySupplier{
		ApiKey:    os.Getenv("NUMVERIFY_API_KEY"),
		EnableSSL: os.Getenv("NUMVERIFY_ENABLE_SSL"),
	}
}

func (s *NumverifySupplier) IsAvailable() bool {
	return s.ApiKey != ""
}

func (s *NumverifySupplier) Validate(internationalNumber string) (res *NumverifyValidateResponse, err error) {
	scheme := "http"

	if s.EnableSSL != "" {
		scheme = "https"
	}

	logrus.
		WithField("number", internationalNumber).
		WithField("scheme", scheme).
		Debug("Running validate operation through Numverify API")

	url := fmt.Sprintf("%s://api.apilayer.com/number_verification/validate?number=%s", scheme, internationalNumber)

	// Build the request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Apikey", s.ApiKey)

	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Fill the response with the data from the JSON
	var result NumverifyValidateResponse

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Error.Info) > 0 {
		return nil, errors.New(result.Error.Info)
	}

	res = &NumverifyValidateResponse{
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
