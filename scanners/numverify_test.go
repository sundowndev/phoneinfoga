package scanners

import (
	"testing"

	assertion "github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
	"gopkg.in/sundowndev/phoneinfoga.v2/utils"
)

func TestNumverifyScanner(t *testing.T) {
	assert := assertion.New(t)

	t.Run("should succeed", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		expectedResult := NumverifyScannerResponse{
			Valid:               true,
			Number:              "79516566591",
			LocalFormat:         "9516566591",
			InternationalFormat: "+79516566591",
			CountryPrefix:       "+7",
			CountryCode:         "RU",
			CountryName:         "Russian Federation",
			Location:            "Saint Petersburg and Leningrad Oblast",
			Carrier:             "OJSC St. Petersburg Telecom (OJSC Tele2-Saint-Petersburg)",
			LineType:            "mobile",
		}

		number, _ := LocalScan("+79516566591")

		gock.New("http://numverify.com").
			Get("/").
			Reply(200).BodyString(`<html><body><input type="hidden" name="scl_request_secret" value="secret"/></body></html>`)

		gock.New("https://numverify.com").
			Get("/php_helper_scripts/phone_api.php").
			MatchParam("secret_key", "5ad5554ac240e4d3d31107941b35a5eb").
			MatchParam("number", number.International).
			Reply(200).
			JSON(expectedResult)

		result, err := numverifyScanCLI(utils.LoggerService, number)

		assert.Nil(err, "they should be equal")
		assert.Equal(result, &expectedResult, "they should be equal")

		assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
	})

	t.Run("should return invalid number", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		expectedResult := NumverifyScannerResponse{
			Valid:               false,
			Number:              "",
			LocalFormat:         "",
			InternationalFormat: "",
			CountryPrefix:       "",
			CountryCode:         "",
			CountryName:         "",
			Location:            "",
			Carrier:             "",
			LineType:            "",
		}

		number, _ := LocalScan("+123456789")

		gock.New("http://numverify.com").
			Get("/").
			Reply(200).BodyString(`<html><body><input type="hidden" name="scl_request_secret" value="secret"/></body></html>`)

		gock.New("https://numverify.com").
			Get("/php_helper_scripts/phone_api.php").
			MatchParam("secret_key", "7ccde16e862dfe7681297713e9e9cadb").
			MatchParam("number", number.International).
			Reply(200).
			JSON(expectedResult)

		result, err := numverifyScanCLI(utils.LoggerService, number)

		assert.Nil(err, "they should be equal")
		assert.Equal(result, &expectedResult, "they should be equal")

		assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
	})

	t.Run("should return empty response and handle error properly", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		number, _ := LocalScan("+123456789")

		gock.New("http://numverify.com").
			Get("/").
			Reply(200).BodyString(`<html><body><input type="hidden" name="scl_request_secret" value="secret"/></body></html>`)

		gock.New("https://numverify.com").
			Get("/php_helper_scripts/phone_api.php").
			MatchParam("secret_key", "7ccde16e862dfe7681297713e9e9cadb").
			MatchParam("number", number.International).
			Reply(200)

		result, err := numverifyScanCLI(utils.LoggerService, number)

		assert.EqualError(err, "EOF", "they should be equal")
		assert.Nil(result, "they should be equal")

		assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
	})
}
