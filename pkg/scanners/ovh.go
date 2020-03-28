package scanners

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// OVHAPIResponseNumber is a type that describes an OVH number range
type OVHAPIResponseNumber struct {
	MatchingCriteria    interface{} `json:"matchingCriteria"`
	City                string      `json:"city"`
	ZneList             []string    `json:"zneList"`
	InternationalNumber string      `json:"internationalNumber"`
	Country             string      `json:"country"`
	AskedCity           interface{} `json:"askedCity"`
	ZipCode             string      `json:"zipCode"`
	Number              string      `json:"number"`
	Prefix              int         `json:"prefix"`
}

// OVHScannerResponse is the OVH scanner response
type OVHScannerResponse struct {
	Found       bool   `json:"found"`
	NumberRange string `json:"numberRange"`
	City        string `json:"city"`
	ZipCode     string `json:"zipCode"`
}

// OVHScan fetches OVH's REST API
func OVHScan(number *Number) (res *OVHScannerResponse, err error) {
	countryCode := strings.ToLower(number.Country)
	url := fmt.Sprintf("https://api.ovh.com/1.0/telephony/number/detailedZones?country=%s", countryCode)

	// Build the request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Fill the response with the data from the JSON
	var results []OVHAPIResponseNumber

	// Use json.Decode for reading streams of JSON data
	json.NewDecoder(response.Body).Decode(&results)

	var foundNumber OVHAPIResponseNumber

	rt := reflect.TypeOf(results)
	if rt.Kind() == reflect.Slice && len(number.RawLocal) > 6 {
		askedNumber := number.RawLocal[0:6] + "xxxx"

		for _, result := range results {
			if result.Number == askedNumber {
				foundNumber = result
			}
		}
	}

	res = &OVHScannerResponse{
		Found:       len(foundNumber.Number) > 0,
		NumberRange: foundNumber.Number,
		City:        foundNumber.City,
		ZipCode:     foundNumber.ZipCode,
	}

	return res, nil
}
