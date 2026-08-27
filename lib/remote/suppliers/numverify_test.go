package suppliers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/url"
	"os"
	"testing"
)

func TestNumverifySupplierSuccessCustomApiKey(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "11115551212"
	apikey := "5ad5554ac240e4d3d31107941b35a5eb"

	expectedResult := &NumverifyValidateResponse{
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

	gock.New("https://api.apilayer.com").
		Get("/number_verification/validate").
		MatchHeader("Apikey", apikey).
		MatchParam("number", number).
		Reply(200).
		JSON(expectedResult)

	s := NewNumverifySupplier()

	got, err := s.Request().SetApiKey(apikey).ValidateNumber(number)
	assert.Nil(t, err)

	assert.Equal(t, expectedResult, got)
}

func TestNumverifySupplierError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "11115551212"
	apikey := "5ad5554ac240e4d3d31107941b35a5eb"

	expectedResult := &NumverifyErrorResponse{
		Message: "You have exceeded your daily\\/monthly API rate limit. Please review and upgrade your subscription plan at https:\\/\\/apilayer.com\\/subscriptions to continue.",
	}

	gock.New("https://api.apilayer.com").
		Get("/number_verification/validate").
		MatchHeader("Apikey", apikey).
		MatchParam("number", number).
		Reply(429).
		JSON(expectedResult)

	s := NewNumverifySupplier()

	got, err := s.Request().SetApiKey(apikey).ValidateNumber(number)
	assert.Nil(t, got)
	assert.Equal(t, errors.New("You have exceeded your daily\\/monthly API rate limit. Please review and upgrade your subscription plan at https:\\/\\/apilayer.com\\/subscriptions to continue."), err)
}

func TestNumverifySupplierFallsBackToLegacyEndpointOn401(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "17027515054"
	apikey := "legacy-style-free-key"

	gock.New("https://api.apilayer.com").
		Get("/number_verification/validate").
		MatchParam("number", number).
		Reply(401).
		JSON(map[string]string{"message": "Invalid authentication credentials"})

	expectedResult := &NumverifyValidateResponse{
		Valid:               true,
		Number:              "17027515054",
		LocalFormat:         "7027515054",
		InternationalFormat: "+17027515054",
		CountryPrefix:       "+1",
		CountryCode:         "US",
		CountryName:         "United States of America",
		Location:            "Indian Spg",
		Carrier:             "",
		LineType:            "landline",
	}

	gock.New("http://apilayer.net").
		Get("/api/validate").
		MatchParam("access_key", apikey).
		MatchParam("number", number).
		Reply(200).
		JSON(expectedResult)

	s := NewNumverifySupplier()

	got, err := s.Request().SetApiKey(apikey).ValidateNumber(number)
	assert.Nil(t, err)

	assert.Equal(t, expectedResult, got)
}

func TestNumverifySupplierHTTPError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "11115551212"

	_ = os.Setenv("NUMVERIFY_API_KEY", "5ad5554ac240e4d3d31107941b35a5eb")
	defer os.Clearenv()

	dummyError := errors.New("test")

	gock.New("https://api.apilayer.com").
		Get("/number_verification/validate").
		ReplyError(dummyError)

	s := NewNumverifySupplier()

	got, err := s.Request().ValidateNumber(number)
	assert.Nil(t, got)
	assert.Equal(t, &url.Error{
		Op:  "Get",
		URL: "https://api.apilayer.com/number_verification/validate?number=11115551212",
		Err: dummyError,
	}, err)
}
