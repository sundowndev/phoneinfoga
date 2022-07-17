package build

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuild(t *testing.T) {
	t.Run("version and commit default values", func(t *testing.T) {
		assert.Equal(t, "dev", Version)
		assert.Equal(t, "dev", Commit)
		assert.Equal(t, false, IsRelease())
		assert.Equal(t, "dev-dev", String())
	})

	t.Run("version and commit default values", func(t *testing.T) {
		Version = "v2.4.4"
		Commit = "0ba854f"
		assert.Equal(t, true, IsRelease())
		assert.Equal(t, "v2.4.4-0ba854f", String())
	})
}
