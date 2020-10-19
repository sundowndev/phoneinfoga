package scanners

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// NumverifyScannerResponse REST API response
type NumverifyScannerResponse struct {
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

// NumverifyScan fetches Numverify's API
func NumverifyScan(number *Number) (res *NumverifyScannerResponse, err error) {
	html, err := http.Get("http://numverify.com/")
	if err != nil {
		return nil, err
	}
	defer html.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(html.Body)
	if err != nil {
		return nil, err
	}

	secret, _ := doc.Find("[name=\"scl_request_secret\"]").Attr("value")

	// Then fetch REST API
	safeNumber := number.International
	apiKey := md5.Sum([]byte(safeNumber + secret))

	url := fmt.Sprintf("https://numverify.com/php_helper_scripts/phone_api.php?secret_key=%s&number=%s", hex.EncodeToString(apiKey[:]), safeNumber)

	// Build the request
	response, err := http.Get(url)
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
