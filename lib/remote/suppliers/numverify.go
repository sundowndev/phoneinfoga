package suppliers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type NumverifySupplierInterface interface {
	Request() NumverifySupplierRequestInterface
}

type NumverifySupplierRequestInterface interface {
	SetApiKey(string) NumverifySupplierRequestInterface
	ValidateNumber(string) (*NumverifyValidateResponse, error)
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
	Uri string
}

func NewNumverifySupplier() *NumverifySupplier {
	return &NumverifySupplier{
		Uri: "https://api.apilayer.com",
	}
}

type NumverifyRequest struct {
	apiKey string
	uri    string
}

func (s *NumverifySupplier) Request() NumverifySupplierRequestInterface {
	return &NumverifyRequest{uri: s.Uri}
}

func (r *NumverifyRequest) SetApiKey(k string) NumverifySupplierRequestInterface {
	r.apiKey = k
	return r
}

func (r *NumverifyRequest) ValidateNumber(internationalNumber string) (res *NumverifyValidateResponse, err error) {
	logrus.
		WithField("number", internationalNumber).
		Debug("Running validate operation through Numverify API")

	res, statusCode, err := r.validateViaMarketplace(internationalNumber)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}

	// A 401 here means the key isn't an apilayer.com marketplace key - most free
	// keys from numverify.com's own signup are the older style, which only this
	// legacy endpoint accepts (and only over plain HTTP - HTTPS is a paid feature
	// on that plan). Only fall back on 401 so real marketplace-key errors
	// (rate limit, bad number, etc.) still surface normally.
	if statusCode == http.StatusUnauthorized {
		logrus.Debug("Numverify marketplace auth failed, retrying against legacy numverify.com endpoint")
		return r.validateViaLegacyEndpoint(internationalNumber)
	}

	return nil, errors.New("numverify request failed")
}

func (r *NumverifyRequest) validateViaMarketplace(internationalNumber string) (res *NumverifyValidateResponse, statusCode int, err error) {
	url := fmt.Sprintf("%s/number_verification/validate?number=%s", r.uri, internationalNumber)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Apikey", r.apiKey)

	response, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, response.StatusCode, nil
		}
		errorResponse := NumverifyErrorResponse{}
		if err := json.NewDecoder(response.Body).Decode(&errorResponse); err != nil {
			return nil, response.StatusCode, err
		}
		return nil, response.StatusCode, errors.New(errorResponse.Message)
	}

	var result NumverifyValidateResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, response.StatusCode, err
	}

	return &result, response.StatusCode, nil
}

// validateViaLegacyEndpoint talks to the original numverify.com free-tier API:
// plain HTTP only, auth via access_key query param instead of a header.
func (r *NumverifyRequest) validateViaLegacyEndpoint(internationalNumber string) (res *NumverifyValidateResponse, err error) {
	url := fmt.Sprintf("http://apilayer.net/api/validate?access_key=%s&number=%s", r.apiKey, internationalNumber)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result NumverifyValidateResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
