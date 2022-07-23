package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"testing"
)

func TestCustomScanner(t *testing.T) {
	testcases := []struct {
		name      string
		number    *number.Number
		expected  customScannerResponse
		wantError string
	}{
		{
			name: "test successful scan",
			number: func() *number.Number {
				n, _ := number.NewNumber("15556661212")
				return n
			}(),
			expected: customScannerResponse{
				Valid:  true,
				Info:   "This number is known for scams!",
				Hidden: "This will not appear in the output",
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			scanner := &customScanner{}

			if !scanner.ShouldRun(*tt.number) {
				t.Fatal("ShouldRun() should be truthy")
			}

			got, err := scanner.Scan(*tt.number)
			if tt.wantError != "" {
				assert.EqualError(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}
