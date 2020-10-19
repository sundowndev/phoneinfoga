package scanners

import (
	"log"
	"os"

	"gopkg.in/sundowndev/phoneinfoga.v2/utils"
)

// Number is a phone number
type Number struct {
	RawLocal      string `json:"rawLocal"`
	Local         string `json:"local"`
	E164          string `json:"E164"`
	International string `json:"international"`
	CountryCode   int32  `json:"countryCode"`
	Country       string `json:"country"`
	Carrier       string `json:"carrier"`
}

func localScanCLI(l *utils.Logger, number string) (*Number, error) {
	l.Infoln("Running local scan...")

	scan, err := LocalScan(number)

	if err != nil {
		l.Errorln("An error occurred:", err.Error())
		return nil, err
	}

	l.Successln("Local format:", scan.Local)
	l.Successln("E164 format:", scan.E164)
	l.Successln("International format:", scan.International)
	l.Successf("Country found: +%v (%v)", scan.CountryCode, scan.Country)
	l.Successln("Carrier:", scan.Carrier)

	return scan, nil
}

func numverifyScanCLI(l *utils.Logger, number *Number) (*NumverifyScannerResponse, error) {
	l.Infoln("Running Numverify.com scan...")

	scan, err := NumverifyScan(number)

	if err != nil {
		l.Errorln("An error occurred:", err.Error())
		l.Errorln("It may be Numverify not finding information about that number, or the number being invalid.")
		return nil, err
	}

	l.Successf(`Valid: %v`, scan.Valid)
	l.Successln("Number:", scan.Number)
	l.Successln("Local format:", scan.LocalFormat)
	l.Successln("International format:", scan.InternationalFormat)
	l.Successf("Country code: %v (%v)", scan.CountryCode, scan.CountryPrefix)
	l.Successln("Country:", scan.CountryName)
	l.Successln("Location:", scan.Location)
	l.Successln("Carrier:", scan.Carrier)
	l.Successln("Line type:", scan.LineType)

	return scan, nil
}

func googlesearchScanCLI(l *utils.Logger, number *Number, formats ...string) GoogleSearchResponse {
	l.Infoln("Generating Google search dork requests...")

	scan := GoogleSearchScan(number, formats...)

	l.Infoln("Social media footprints")
	for _, dork := range scan.SocialMedia {
		l.Successf(`Link: %v`, dork.URL)
	}

	l.Infoln("Individual footprints")
	for _, dork := range scan.Individuals {
		l.Successf(`Link: %v`, dork.URL)
	}

	l.Infoln("Reputation footprints")
	for _, dork := range scan.Reputation {
		l.Successf(`Link: %v`, dork.URL)
	}

	l.Infoln("Temporary number providers footprints")
	for _, dork := range scan.DisposableProviders {
		l.Successf(`Link: %v`, dork.URL)
	}

	return scan
}

func ovhScanCLI(l *utils.Logger, number *Number) (*OVHScannerResponse, error) {
	l.Infoln("Running OVH API scan...")

	scan, err := OVHScan(number)

	if err != nil {
		l.Errorln("An error occurred:", err.Error())
		return nil, err
	}

	l.Successf(`Found: %v`, scan.Found)
	l.Successf(`Number range: %v`, scan.NumberRange)
	l.Successln("City:", scan.City)
	l.Successln("Zip code:", scan.ZipCode)

	return scan, nil
}

// ScanCLI Run scans with CLI output
func ScanCLI(number string) {
	num, err := localScanCLI(utils.LoggerService, number)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	numverifyScanCLI(utils.LoggerService, num)
	googlesearchScanCLI(utils.LoggerService, num)
	ovhScanCLI(utils.LoggerService, num)
}
