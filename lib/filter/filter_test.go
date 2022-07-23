package filter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterEngine(t *testing.T) {
	testcases := []struct {
		name     string
		rules    []string
		expected map[string]bool
	}{
		{
			name:  "test googlesearch is ignored",
			rules: []string{"googlesearch"},
			expected: map[string]bool{
				"googlesearch": true,
				"numverify":    false,
			},
		},
		{
			name:  "test none is ignored",
			rules: []string{},
			expected: map[string]bool{
				"googlesearch": false,
				"numverify":    false,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEngine()
			e.AddRule(tt.rules...)
			for r, isIgnored := range tt.expected {
				assert.Equal(t, isIgnored, e.Match(r))
			}
		})
	}
}
