package remote

import (
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"testing"
)

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
			remote := NewLibrary()
			remote.AddScanner(scanner)

			if !scanner.ShouldRun() {
				t.Fatal("ShouldRun() should be truthy")
			}

			got, errs := remote.Scan(tt.number)
			if len(tt.wantErrors) > 0 {
				assert.Equal(t, tt.wantErrors, errs)
			} else {
				assert.Len(t, errs, 0)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}
