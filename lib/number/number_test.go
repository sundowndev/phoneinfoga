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
				RawLocal:      "0678342311",
				Local:         "06 78 34 23 11",
				E164:          "+33678342311",
				International: "33678342311",
				CountryCode:   33,
				Country:       "FR",
				Carrier:       "",
			},
		},
		{
			name:  "should succeed to parse number",
			input: "15552221212",
			expected: &Number{
				Valid:         false,
				RawLocal:      "5552221212",
				Local:         "(555) 222-1212",
				E164:          "+15552221212",
				International: "15552221212",
				CountryCode:   1,
				Country:       "",
				Carrier:       "",
			},
		},

		{
			name:     "should fail to parse number",
			input:    "wrong",
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
