// DIRECTORY: ~/phoneinfoga/lib/remote/library.go

package remote

import (
	"context"
	"fmt"
	"log"
	"os"
	"plugin"

	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/scanner"
)

// ScannerOptions defines all available configuration options for remote scanners.
type ScannerOptions struct {
	// API keys for external services
	NumverifyAPIKey string
}

// Library holds the list of scanners to be executed
type Library struct {
	scanners []scanner.Scanner
	filter   Filter
}

// NewLibrary returns a new library instance initialized with the given filter.
func NewLibrary(filter Filter) *Library {
	return &Library{filter: filter}
}

// InitScanners initializes the remote scanner library with default and optional custom scanners.
func InitScanners(library *Library) {
	// Options structure for scanners (used to pass API keys)
	opts := &ScannerOptions{
		NumverifyAPIKey: os.Getenv("NUMVERIFY_API_KEY"),
	}
	
	// Add all default and custom remote scanners
	for _, s := range RemoteScanners(opts) {
		library.Add(s)
	}
}

// RemoteScanners returns a slice of default scanner implementations.
func RemoteScanners(opts *ScannerOptions) []scanner.Scanner {
	// This list contains all external data gathering sources.
	return []scanner.Scanner{
		NewNumverifyScanner(opts),
		NewGoogleSearch(opts), // Assuming this file also exists
		NewOVHScanner(opts),   // Assuming this file also exists
	}
}

// Add adds a scanner to the library if it is not filtered.
func (l *Library) Add(s scanner.Scanner) {
	if l.filter.IsFiltered(s.Name()) {
		return
	}
	l.scanners = append(l.scanners, s)
}

// ... (Rest of the original file functions like AddPlugin and Scan, if they exist)
