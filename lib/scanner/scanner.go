// DIRECTORY: ~/phoneinfoga/lib/scanner/scanner.go

package scanner

import (
	"context"
	"github.com/sundowndev/phoneinfoga/v2/lib/number" 
)

// ScanResult contains the result of a single scan operation.
type ScanResult struct {
	Footprints []number.Footprint
}

// Scanner is the interface that all scanners must implement.
type Scanner interface {
	Name() string
	IsSlow() bool
	Scan(ctx context.Context, number *number.Number) (*ScanResult, error)
}
