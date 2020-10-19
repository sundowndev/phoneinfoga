package config

import (
	"testing"

	assertion "github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assertion.New(t)

	t.Run("Version", func(t *testing.T) {
		t.Run("version should be unknown by default", func(t *testing.T) {
			assert.Equal("unknown", Version)
			assert.Equal("unknown", Commit)
		})
	})
}
