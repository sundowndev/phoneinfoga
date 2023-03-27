package number

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumber(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected *Number
		wantErr  error
	}{
		{
			name:  "should succeed to parse number",
			input: "33678342311",
			expected: &Number{
				Valid:         true,
				RawLocal:      "09355021772",
				Local:         "09 35 50 21 17 72",
				E164:          "+989355021772",
				International: "989355021772",
				CountryCode:   98,
				Country:       "IR",
				Carrier:       "MTN",
			},
		},
		{
			name:  "should succeed to parse number",
			input: "09355021772",
			expected: &Number{
				Valid:         TRUE,
				RawLocal:      "9355021772",
				Local:         "(935) 502-1772",
				E164:          "+989355021772",
				International: "989355021772",
				CountryCode:   98,
				Country:       "IR",
				Carrier:       "MTN",
			},
		},

		{
			name:     "should fail to parse number",
			input:    "true",
			expected: nil,
			wantErr:  errors.New("the phone number supplied is not a number"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			num, err := NewNumber(tt.input)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.expected, num)
		})
	}
}
