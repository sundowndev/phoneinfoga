package suppliers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"net/http"
	"reflect"
	"strings"
)

type OVHSupplierInterface interface {
	Search(number.Number) (*OVHScannerResponse, error)
}

// OVHAPIResponseNumber is a type that describes an OVH number range
type OVHAPIResponseNumber struct {
	MatchingCriteria    interface{} `json:"matchingCriteria"`
	City                string      `json:"city"`
	ZneList             []string    `json:"zne-list"`
	InternationalNumber string      `json:"internationalNumber"`
	Country             string      `json:"country"`
	AskedCity           interface{} `json:"askedCity"`
	ZipCode             string      `json:"zipCode"`
	Number              string      `json:"number"`
	Prefix              int         `json:"prefix"`
}

type OVHAPIErrorResponse struct {
	Message string `json:"message"`
}

// OVHScannerResponse is the OVH scanner response
type OVHScannerResponse struct {
	Found       bool
	NumberRange string
	City        string
	ZipCode     string
}

type OVHSupplier struct{}

func NewOVHSupplier() *OVHSupplier {
	return &OVHSupplier{}
}

func (s *OVHSupplier) Search(num number.Number) (*OVHScannerResponse, error) {
	countryCode := strings.ToLower(num.Country)

	if countryCode == "" {
		return nil, fmt.Errorf("country code +%d wasn't recognized", num.CountryCode)
	}

	// Build the request
	response, err := http.Get(fmt.Sprintf("https://api.ovh.com/1.0/telephony/number/detailedZones?country=%s", countryCode))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		var result OVHAPIErrorResponse
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(result.Message)
	}

	// Fill the response with the data from the JSON
	var results []OVHAPIResponseNumber

	// Use json.Decode for reading streams of JSON data
	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	var foundNumber OVHAPIResponseNumber

	rt := reflect.TypeOf(results)
	if rt.Kind() == reflect.Slice && len(num.RawLocal) > 6 {
		askedNumber := num.RawLocal[0:6] + "xxxx"

		for _, result := range results {
			if result.Number == askedNumber {
				foundNumber = result
			}
		}
	}

	return &OVHScannerResponse{
		Found:       len(foundNumber.Number) > 0,
		NumberRange: foundNumber.Number,
		City:        foundNumber.City,
		ZipCode:     foundNumber.ZipCode,
	}, nil
}
