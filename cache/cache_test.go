package cache

import (
	"gopkg.in/sundowndev/phoneinfoga.v2/scanners"
	"testing"

	assertion "github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	assert := assertion.New(t)

	t.Run("database was already initialized", func(t *testing.T) {
		c := &Cache{Directory: "."}
		err := c.Init()
		defer c.Close()

		assert.NoError(err)

		err = c.Init()

		assert.EqualError(err, "database was already initialized")
	})

	t.Run("check if number is cached", func(t *testing.T) {

	})

	t.Run("add number to cache", func(t *testing.T) {
		c := &Cache{Directory: "."}
		err := c.Init()
		defer c.Close()

		assert.NoError(err)

		n, err := scanners.LocalScan("7185212994")
		assert.NoError(err)

		results := &scanners.ScanResult{
			Local: n,
		}

		err = c.CacheResults(results)
		assert.NoError(err)

		cachedResults, err := c.GetResults(n)
		assert.NoError(err)

		assert.Equal("7185212994", cachedResults.Local.International)
	})
}
