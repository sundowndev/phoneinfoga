package scanners

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/sundowndev/phoneinfoga/pkg/utils"
)

// Numverify REST API response
type Numverify struct {
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
func NumverifyScan(number *Number) (res *Numverify, err error) {
	if err != nil {
		utils.LoggerService.Errorln("The number is not valid")
		os.Exit(0)
	}

	response, _, errs := gorequest.New().Get("http://numverify.com/").End()
	if errs != nil {
		log.Fatal(errs)
	}
	defer response.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	secret, _ := doc.Find("[name=\"scl_request_secret\"]").Attr("value")

	// Then fetch infos
	safeNumber := number.International
	apiKey := md5.Sum([]byte(safeNumber + secret))

	url := fmt.Sprintf("https://numverify.com/php_helper_scripts/phone_api.php?secret_key=%s&number=%s", hex.EncodeToString(apiKey[:]), safeNumber)

	// Build the request
	response2, _, errs := gorequest.New().Get(url).End()
	if errs != nil {
		log.Fatal(errs)
	}
	defer response2.Body.Close()

	// Fill the response with the data from the JSON
	var result Numverify

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(response2.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	res = &Numverify{
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
