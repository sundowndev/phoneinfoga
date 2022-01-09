package remote

import (
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

func NewOVHScanner(s suppliers.OVHSupplierInterface) *ovhScanner {
	return &ovhScanner{client: s}
}

func (s *ovhScanner) Identifier() string {
	return OVH
}

func (s *ovhScanner) ShouldRun() bool {
	return true
}

func (s *ovhScanner) Scan(n *number.Number) (interface{}, error) {
	res, err := s.client.Search(*n)
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
