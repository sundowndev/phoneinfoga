package remote

import (
	"errors"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

const Nomorobo = "nomorobo"

type nomoroboScanner struct {
	client suppliers.NomoroboSupplierInterface
}

type NomoroboScannerResponse struct {
	Description string `json:"description" console:"Description,omitempty"`
}

func NewNomoroboScanner(s suppliers.NomoroboSupplierInterface) Scanner {
	return &nomoroboScanner{client: s}
}

func (s *nomoroboScanner) Name() string {
	return Nomorobo
}

func (s *nomoroboScanner) Description() string {
	return "Look up a US/Canada number's caller classification on Nomorobo (free, no API key)."
}

func (s *nomoroboScanner) DryRun(n number.Number, _ ScannerOptions) error {
	if n.CountryCode != 1 || len(n.RawLocal) != 10 {
		return errors.New("only US/Canada numbers are supported")
	}
	return nil
}

func (s *nomoroboScanner) Run(n number.Number, _ ScannerOptions) (interface{}, error) {
	res, err := s.client.Lookup(n.RawLocal)
	if err != nil {
		return nil, err
	}

	return NomoroboScannerResponse{Description: res.Description}, nil
}
