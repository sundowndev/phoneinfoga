package remote

import (
	"fmt"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

const OVH = "ovh"

type ovhScanner struct {
	client suppliers.OVHSupplierInterface
}

// OVHScannerResponse is the OVH scanner response
type OVHScannerResponse struct {
	Found       bool   `json:"found" console:"Found"`
	NumberRange string `json:"number_range,omitempty" console:"Number range,omitempty"`
	City        string `json:"city,omitempty" console:"City,omitempty"`
	ZipCode     string `json:"zip_code,omitempty" console:"Zip code,omitempty"`
}

func NewOVHScanner(s suppliers.OVHSupplierInterface) Scanner {
	return &ovhScanner{client: s}
}

func (s *ovhScanner) Name() string {
	return OVH
}

func (s *ovhScanner) Description() string {
	return "Search a phone number through the OVH Telecom REST API."
}

func (s *ovhScanner) DryRun(n number.Number, _ ScannerOptions) error {
	if !s.isSupported(n.CountryCode) {
		return fmt.Errorf("country code %d is not supported", n.CountryCode)
	}
	return nil
}

func (s *ovhScanner) Run(n number.Number, _ ScannerOptions) (interface{}, error) {
	res, err := s.client.Search(n)
	if err != nil {
		return nil, err
	}

	data := OVHScannerResponse{
		Found:       res.Found,
		NumberRange: res.NumberRange,
		City:        res.City,
		ZipCode:     res.ZipCode,
	}

	return data, nil
}

func (s *ovhScanner) supportedCountryCodes() []int32 {
	// See https://api.ovh.com/console/#/telephony/number/detailedZones#GET
	return []int32{33, 32, 44, 34, 41}
}

func (s *ovhScanner) isSupported(code int32) bool {
	supported := false
	for _, c := range s.supportedCountryCodes() {
		if code == c {
			supported = true
		}
	}
	return supported
}
