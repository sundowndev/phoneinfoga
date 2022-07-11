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

// NewScanner creates a new instance of this scanner.
// The name of the function MUST be NewScanner.
func NewScanner() remote.Scanner {
	return &customScanner{}
}

// Identifier returns the unique identifier this
// scanner should be associated to.
// Please keep in mind this value could be used for
// automation and so must remain simple.
func (s *customScanner) Identifier() string {
	return "customscanner"
}

// ShouldRun returns a boolean indicating whether
// this scanner should be used or not.
// This can be useful to check for authentication and
// avoid running the scanner when it just can't work.
func (s *customScanner) ShouldRun() bool {
	return true
}

// Scan does the actual scan of the phone number.
// Note this function will be executed in a goroutine.
func (s *customScanner) Scan(n *number.Number) (interface{}, error) {
	data := customScannerResponse{
		Valid:  true,
		Info:   "This number is known for scams!",
		Hidden: "This will not appear in the output",
	}
	return data, nil
}
