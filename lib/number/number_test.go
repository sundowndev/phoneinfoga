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
				CustomFormats: []string{},
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
				CustomFormats: []string{},
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

func TestFormatTemplate(t *testing.T) {
	cases := []struct {
		name     string
		number   string
		format   string
		expected []string
		wantErr  error
	}{
		{
			name:     "should succeed to format with template",
			number:   "+15552221212",
			format:   "xxx-xxx-xxxx",
			expected: []string{"555-222-1212"},
		},
		{
			name:     "should succeed to format with template",
			number:   "+33678342211",
			format:   "0x-Xx-xx-xX-xx",
			expected: []string{"06-78-34-22-11"},
		},
		{
			name:     "should succeed to format with template",
			number:   "+911401871759",
			format:   "+{{.CountryCode}} xxxx-xxxxxx",
			expected: []string{"+91 1401-871759"},
		},
		{
			name:     "should fail to format with template",
			number:   "+911401871759",
			format:   "+{{.DummyVar}} xxxx-xxxxxx",
			expected: []string{"+{{.DummyVar}} 1401-871759"},
		},
		{
			name:     "should fail to format with template",
			number:   "+333333",
			format:   "+{{.DummyVar}} xxxx-xxxxxx",
			expected: []string{"+{{.DummyVar}} 3333-xxxxxx"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			n, err := NewNumber(tt.number, tt.format)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.expected, n.CustomFormats)
		})
	}
}
