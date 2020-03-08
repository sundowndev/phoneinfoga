package scanners

import (
	"os"

	"github.com/sundowndev/phoneinfoga/pkg/utils"
)

// Number is a phone number
type Number struct {
	Local         string `json:"local"`
	E164          string `json:"E164"`
	International string `json:"international"`
	CountryCode   int32  `json:"countryCode"`
	Country       string `json:"country"`
	Carrier       string `json:"carrier"`
}

func localScanCLI(number string) *Number {
	utils.LoggerService.Infoln("Running local scan...")

	scan, err := LocalScan(number)

	if err != nil {
		utils.LoggerService.Errorln("The number is not valid")
		os.Exit(0)
	}

	utils.LoggerService.Successln("Local format:", scan.Local)
	utils.LoggerService.Successln("E164 format:", scan.E164)
	utils.LoggerService.Successln("International format:", scan.International)
	utils.LoggerService.Successf("Country found: +%v (%v)", scan.CountryCode, scan.Country)
	utils.LoggerService.Successln("Carrier:", scan.Carrier)

	return scan
}

func numverifyScanCLI(number *Number) {
	utils.LoggerService.Infoln("Running Numverify.com scan...")

	scan, err := NumverifyScan(number)

	if err != nil {
		utils.LoggerService.Errorln("The number is not valid")
		os.Exit(0)
	}

	utils.LoggerService.Successf(`Valid: %v`, scan.Valid)
	utils.LoggerService.Successln("Number:", scan.Number)
	utils.LoggerService.Successln("Local format:", scan.LocalFormat)
	utils.LoggerService.Successln("International format:", scan.InternationalFormat)
	utils.LoggerService.Successf("Country code: %v (%v)", scan.CountryCode, scan.CountryPrefix)
	utils.LoggerService.Successln("Country:", scan.CountryName)
	utils.LoggerService.Successln("Location:", scan.Location)
	utils.LoggerService.Successln("Carrier:", scan.Carrier)
	utils.LoggerService.Successln("Line type:", scan.LineType)
}

func googlesearchScanCLI(number *Number) {
	utils.LoggerService.Infoln("Generating Google search dork requests...")

	scan := GoogleSearchScan(number)

	utils.LoggerService.Infoln("Social media footprints")
	for _, dork := range scan.SocialMedia {
		utils.LoggerService.Successf(`Link: %v`, dork.URL)
	}

	utils.LoggerService.Infoln("Individual footprints")
	for _, dork := range scan.Individuals {
		utils.LoggerService.Successf(`Link: %v`, dork.URL)
	}

	utils.LoggerService.Infoln("Reputation footprints")
	for _, dork := range scan.Reputation {
		utils.LoggerService.Successf(`Link: %v`, dork.URL)
	}

	utils.LoggerService.Infoln("Temporary number providers footprints")
	for _, dork := range scan.DisposableProviders {
		utils.LoggerService.Successf(`Link: %v`, dork.URL)
	}
}

// ScanCLI Run scans with CLI output
func ScanCLI(number string) {
	num := localScanCLI(number)

	numverifyScanCLI(num)
	googlesearchScanCLI(num)
}
