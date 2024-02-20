package remote

import (
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"testing"
)

func TestLocalScanner_Metadata(t *testing.T) {
	scanner := NewLocalScanner()
	assert.Equal(t, Local, scanner.Name())
	assert.NotEmpty(t, scanner.Description())
}

func TestLocalScanner(t *testing.T) {
	testcases := []struct {
		name       string
		number     *number.Number
		expected   map[string]interface{}
		wantErrors map[string]error
	}{
		{
			name: "successful scan",
			number: func() *number.Number {
				n, _ := number.NewNumber("15556661212")
				return n
			}(),
			expected: map[string]interface{}{
				"local": LocalScannerResponse{
					RawLocal:      "5556661212",
					Local:         "(555) 666-1212",
					E164:          "+15556661212",
					International: "15556661212",
					CountryCode:   1,
				},
			},
			wantErrors: map[string]error{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewLocalScanner()
			remote := NewLibrary(filter.NewEngine())
			remote.AddScanner(scanner)

			if scanner.DryRun(*tt.number, ScannerOptions{}) != nil {
				t.Fatal("DryRun() should return nil")
			}

			got, errs := remote.Scan(tt.number, ScannerOptions{})
			if len(tt.wantErrors) > 0 {
				assert.Equal(t, tt.wantErrors, errs)
			} else {
				assert.Len(t, errs, 0)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}
