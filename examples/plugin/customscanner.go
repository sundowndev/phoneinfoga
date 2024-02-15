package main

import (
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
)

type customScanner struct{}

type customScannerResponse struct {
	Valid  bool   `json:"valid" console:"Valid"`
	Info   string `json:"info" console:"Info"`
	Hidden string `json:"-" console:"-"`
}

// Name returns the unique name this scanner.
func (s *customScanner) Name() string {
	return "customscanner"
}

// Description returns a short description for this scanner.
func (s *customScanner) Description() string {
	return "This is a dummy scanner"
}

// DryRun returns an error indicating whether
// this scanner can be used with the given number.
// This can be useful to check for authentication or
// country code support for example, and avoid running
// the scanner when it just can't work.
func (s *customScanner) DryRun(n number.Number, opts remote.ScannerOptions) error {
	return nil
}

// Run does the actual scan of the phone number.
// Note this function will be executed in a goroutine.
func (s *customScanner) Run(n number.Number, opts remote.ScannerOptions) (interface{}, error) {
	data := customScannerResponse{
		Valid:  true,
		Info:   "This number is known for scams!",
		Hidden: "This will not appear in the output",
	}
	return data, nil
}

func init() {
	remote.RegisterPlugin(&customScanner{})
}
