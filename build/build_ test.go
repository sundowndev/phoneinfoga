package build

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("Version", func(t *testing.T) {
		t.Run("version should be unknown by default", func(t *testing.T) {
			assert.Equal(t, "unknown", Version)
			assert.Equal(t, "unknown", Commit)
		})
	})
}
