package remote

import (
	"errors"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
	"strings"
)

const Numverify = "numverify"

type numverifyScanner struct {
	client suppliers.NumverifySupplierInterface
}

type NumverifyScannerResponse struct {
	Valid               bool   `json:"valid" console:"Valid"`
	Number              string `json:"number" console:"Number,omitempty"`
	LocalFormat         string `json:"local_format" console:"Local format,omitempty"`
	InternationalFormat string `json:"international_format" console:"International format,omitempty"`
	CountryPrefix       string `json:"country_prefix" console:"Country prefix,omitempty"`
	CountryCode         string `json:"country_code" console:"Country code,omitempty"`
	CountryName         string `json:"country_name" console:"Country name,omitempty"`
	Location            string `json:"location" console:"Location,omitempty"`
	Carrier             string `json:"carrier" console:"Carrier,omitempty"`
	LineType            string `json:"line_type" console:"Line type,omitempty"`
	SpoofRisk           string `json:"spoof_risk" console:"Spoof risk,omitempty"`
}

// knownVoipCarriers is a non-exhaustive list of wholesale/CLEC carriers commonly
// underlying VOIP and app-based numbers (Google Voice, Skype, TextNow, RingCentral,
// etc. are all provisioned through carriers like these rather than showing up by name).
var knownVoipCarriers = []string{
	"bandwidth", "level 3", "level3", "lumen", "twilio", "onvoy", "inteliquent",
	"neutral tandem", "peerless", "flowroute", "telnyx", "sinch", "vonage",
	"ringcentral", "8x8", "zipwhip", "voip innovations", "voxbone", "skype",
	"broadvoice", "callcentric", "vitelity",
}

// assessSpoofRisk gives a rough, honest signal for "is the caller ID here likely
// spoofable" based on carrier/line-type data alone. It can NOT determine whether a
// call actually was spoofed, or where a call really originated - no phone-number
// lookup can do that without carrier-level call records. This is a proxy, not a verdict.
func assessSpoofRisk(carrier, lineType string) string {
	lowerCarrier := strings.ToLower(carrier)
	lowerLineType := strings.ToLower(lineType)

	if lowerLineType == "voip" {
		return "elevated - line type reported as VOIP"
	}

	for _, known := range knownVoipCarriers {
		if strings.Contains(lowerCarrier, known) {
			return "elevated - carrier is a VOIP/wholesale provider (" + carrier + ")"
		}
	}

	if lowerLineType == "landline" && carrier == "" {
		return "uncertain - \"landline\" with no carrier attribution; free-tier databases often mislabel ported/VOIP numbers this way"
	}

	if lowerLineType == "mobile" && carrier != "" {
		return "low - mobile line with a named carrier"
	}

	if lowerLineType == "landline" && carrier != "" {
		return "low - landline with a named carrier"
	}

	return "unknown - not enough data"
}

func NewNumverifyScanner(s suppliers.NumverifySupplierInterface) Scanner {
	return &numverifyScanner{client: s}
}

func (s *numverifyScanner) Name() string {
	return Numverify
}

func (s *numverifyScanner) Description() string {
	return "Request info about a given phone number through the Numverify API."
}

func (s *numverifyScanner) DryRun(_ number.Number, opts ScannerOptions) error {
	if opts.GetStringEnv("NUMVERIFY_API_KEY") != "" {
		return nil
	}
	return errors.New("API key is not defined")
}

func (s *numverifyScanner) Run(n number.Number, opts ScannerOptions) (interface{}, error) {
	apiKey := opts.GetStringEnv("NUMVERIFY_API_KEY")

	res, err := s.client.Request().SetApiKey(apiKey).ValidateNumber(n.International)
	if err != nil {
		return nil, err
	}

	data := NumverifyScannerResponse{
		Valid:               res.Valid,
		Number:              res.Number,
		LocalFormat:         res.LocalFormat,
		InternationalFormat: res.InternationalFormat,
		CountryPrefix:       res.CountryPrefix,
		CountryCode:         res.CountryCode,
		CountryName:         res.CountryName,
		Location:            res.Location,
		Carrier:             res.Carrier,
		LineType:            res.LineType,
		SpoofRisk:           assessSpoofRisk(res.Carrier, res.LineType),
	}

	return data, nil
}
