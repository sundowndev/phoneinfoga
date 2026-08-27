package remote

import (
	"errors"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

const Veriphone = "veriphone"

type veriphoneScanner struct {
	client suppliers.VeriphoneSupplierInterface
}

type VeriphoneScannerResponse struct {
	PhoneValid          bool   `json:"phone_valid" console:"Valid"`
	PhoneType           string `json:"phone_type" console:"Phone type,omitempty"`
	PhoneRegion         string `json:"phone_region" console:"Phone region,omitempty"`
	Country             string `json:"country" console:"Country,omitempty"`
	CountryCode         string `json:"country_code" console:"Country code,omitempty"`
	InternationalNumber string `json:"international_number" console:"International number,omitempty"`
	Carrier             string `json:"carrier" console:"Carrier,omitempty"`
	SpoofRisk           string `json:"spoof_risk" console:"Spoof risk,omitempty"`
}

func NewVeriphoneScanner(s suppliers.VeriphoneSupplierInterface) Scanner {
	return &veriphoneScanner{client: s}
}

func (s *veriphoneScanner) Name() string {
	return Veriphone
}

func (s *veriphoneScanner) Description() string {
	return "Cross-check carrier/line-type via Veriphone, a second independent source from Numverify (free, needs VERIPHONE_API_KEY)."
}

func (s *veriphoneScanner) DryRun(_ number.Number, opts ScannerOptions) error {
	if opts.GetStringEnv("VERIPHONE_API_KEY") != "" {
		return nil
	}
	return errors.New("API key is not defined")
}

func (s *veriphoneScanner) Run(n number.Number, opts ScannerOptions) (interface{}, error) {
	apiKey := opts.GetStringEnv("VERIPHONE_API_KEY")

	res, err := s.client.Verify(apiKey, n.International)
	if err != nil {
		return nil, err
	}

	return VeriphoneScannerResponse{
		PhoneValid:          res.PhoneValid,
		PhoneType:           res.PhoneType,
		PhoneRegion:         res.PhoneRegion,
		Country:             res.Country,
		CountryCode:         res.CountryCode,
		InternationalNumber: res.InternationalNumber,
		Carrier:             res.Carrier,
		SpoofRisk:           assessSpoofRisk(res.Carrier, res.PhoneType),
	}, nil
}
