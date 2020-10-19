package config

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	t.Run("Version", func(t *testing.T) {
		t.Run("version should be use semver format", func(t *testing.T) {
			matched, err := regexp.MatchString(`^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$`, Version)

			assert.Equal(matched, true, "should be equal")
			assert.Equal(err, nil, "should be equal")
		})
	})
}
