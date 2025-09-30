// DIRECTORY: ~/phoneinfoga/lib/remote/custom_scanner.go

package remote

import (
	"context"
	"fmt"

	"github.com/sundowndev/phoneinfoga/v2/lib/number"
)

// CustomScanner defines the structure for your new scanner.
type **CustomScanner** struct {
	Options *ScannerOptions
}

// NewCustomScanner is the constructor function required by the main application.
func NewCustomScanner(options *ScannerOptions) Scanner {
	return &CustomScanner{Options: options}
}

// Name returns the unique identifier for the scanner.
func (s *CustomScanner) Name() string {
	return "**CustomScanner**" 
}

// IsSlow indicates if the scanner performs network lookups.
func (s *CustomScanner) IsSlow() bool {
	return true
}

// Scan contains the core logic for checking the number.
func (s *CustomScanner) Scan(ctx context.Context, number *number.Number) (*ScanResult, error) {
	
	// Initialize the result object
	result := &ScanResult{}
	
	// --- YOUR CUSTOM LOGIC GOES HERE ---
	// We use "general" type here to force the output to be displayed in the report.
	
	result.Footprints = append(result.Footprints, number.Footprint{
		Type:    "general", // <--- CHANGE IS HERE! This type is almost always displayed.
		Title:   "Custom Scanner Confirmation",
		Message: fmt.Sprintf("SUCCESS: The %s ran successfully for number %s", s.Name(), number.String()),
		Source:  "CustomScanner V1.0",
	})
	
	// --- END CUSTOM LOGIC ---
	
	return result, nil
}
