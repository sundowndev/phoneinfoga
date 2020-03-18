package scanners

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/parnurzeal/gorequest"
)

type ovhAPIResponseNumber struct {
	MatchingCriteria interface{} `json:"matchingCriteria"`
	City             string      `json:"city"`
	ZneList          []string    `json:"zneList"`
	// type string
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
	response, _, errs := gorequest.New().Get(url).End()
	if errs != nil {
		log.Fatal(errs)
	}
	defer response.Body.Close()

	// Fill the response with the data from the JSON
	var result []ovhAPIResponseNumber

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	askedNumber := number.RawLocal[0:6] + "xxxx"
	var foundNumber ovhAPIResponseNumber

	for _, n := range result {
		if n.Number == askedNumber {
			foundNumber = n
		}
	}

	res = &OVHScannerResponse{
		Found:       foundNumber.Number != "",
		NumberRange: foundNumber.Number,
		City:        foundNumber.City,
		ZipCode:     foundNumber.ZipCode,
	}

	return res, nil
}
