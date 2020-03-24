package scanners

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestOVHScanner(t *testing.T) {
	assert := assert.New(t)

	t.Run("OVHScan", func(t *testing.T) {
		t.Run("should find number on OVH", func(t *testing.T) {
			defer gock.Off() // Flush pending mocks after test execution

			gock.New("https://api.ovh.com").
				Get("/1.0/telephony/number/detailedZones").
				Reply(200).
				JSON([]ovhAPIResponseNumber{
					ovhAPIResponseNumber{
						ZneList:             []string{},
						MatchingCriteria:    "",
						Prefix:              33,
						InternationalNumber: "003336517xxxx",
						Country:             "fr",
						ZipCode:             "",
						Number:              "036517xxxx",
						City:                "Abbeville",
						AskedCity:           "",
					},
				})

			number, _ := LocalScan("+33 0365179268")

			result, err := OVHScan(number)

			assert.Equal(err, nil, "should not be errored")
			assert.Equal(result, &OVHScannerResponse{
				Found:       true,
				NumberRange: "036517xxxx",
				City:        "Abbeville",
				ZipCode:     "",
			}, "they should be equal")

			assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
		})

		t.Run("should not find number on OVH", func(t *testing.T) {
			defer gock.Off() // Flush pending mocks after test execution

			gock.New("https://api.ovh.com").
				Get("/1.0/telephony/number/detailedZones").
				Reply(200).
				JSON([]ovhAPIResponseNumber{
					ovhAPIResponseNumber{
						ZneList:             []string{},
						MatchingCriteria:    "",
						Prefix:              33,
						InternationalNumber: "003336517xxxx",
						Country:             "fr",
						ZipCode:             "",
						Number:              "036517xxxx",
						City:                "Abbeville",
						AskedCity:           "",
					},
				})

			number, _ := LocalScan("+1 555-444-1212")

			result, err := OVHScan(number)

			assert.Equal(err, nil, "should not be errored")
			assert.Equal(result, &OVHScannerResponse{
				Found:       false,
				NumberRange: "",
				City:        "",
				ZipCode:     "",
			}, "they should be equal")

			assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
		})
	})
}
