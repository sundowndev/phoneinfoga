package build

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuild(t *testing.T) {
	t.Run("version and commit default values", func(t *testing.T) {
		assert.Equal(t, "dev", Version)
		assert.Equal(t, "dev", Commit)
	})
}
