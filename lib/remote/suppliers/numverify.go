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

type NumverifyErrorResponse struct {
	Message string `json:"message"`
}

// NumverifyValidateResponse REST API response
type NumverifyValidateResponse struct {
	Valid               bool   `json:"valid"`
	Number              string `json:"number"`
	LocalFormat         string `json:"local_format"`
	InternationalFormat string `json:"international_format"`
	CountryPrefix       string `json:"country_prefix"`
	CountryCode         string `json:"country_code"`
	CountryName         string `json:"country_name"`
	Location            string `json:"location"`
	Carrier             string `json:"carrier"`
	LineType            string `json:"line_type"`
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

func (s *NumverifySupplier) Validate(internationalNumber string) (res *NumverifyValidateResponse, err error) {
	logrus.
		WithField("number", internationalNumber).
		Debug("Running validate operation through Numverify API")

	url := fmt.Sprintf("https://api.apilayer.com/number_verification/validate?number=%s", internationalNumber)

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

	if response.StatusCode >= 400 {
		errorResponse := NumverifyErrorResponse{}
		if err := json.NewDecoder(response.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse.Message)
	}

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
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
