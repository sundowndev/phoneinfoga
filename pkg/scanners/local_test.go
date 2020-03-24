package scanners

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/pkg/utils"
)

func TestLocalScan(t *testing.T) {
	assert := assert.New(t)

	t.Run("should scan number", func(t *testing.T) {
		result := localScanCLI(utils.LoggerService, "+1 718-521-2994")

		expectedResult := &Number{
			RawLocal:      "7185212994",
			Local:         "(718) 521-2994",
			E164:          "+17185212994",
			International: "17185212994",
			CountryCode:   1,
			Country:       "US",
			Carrier:       "",
		}

		assert.Equal(result, expectedResult, "they should be equal")
	})

	t.Run("should fail and return error", func(t *testing.T) {
		_, err := LocalScan("this is not a phone number")

		assert.Equal(err.Error(), "the phone number supplied is not a number", "they should be equal")
	})
}
