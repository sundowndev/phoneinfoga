package config

import (
	"testing"

	"github.com/blang/semver/v4"
	assertion "github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assertion.New(t)

	t.Run("Version", func(t *testing.T) {
		t.Run("version should be use semver format", func(t *testing.T) {
			v, err := semver.Make(Version)
			assert.Nil(err)

			err = v.Validate()
			assert.Nil(err)
		})
	})
}
