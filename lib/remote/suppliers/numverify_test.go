package suppliers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/url"
	"os"
	"testing"
)

func TestNumverifySupplierSuccess(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "11115551212"

	_ = os.Setenv("NUMVERIFY_API_KEY", "5ad5554ac240e4d3d31107941b35a5eb")
	_ = os.Setenv("NUMVERIFY_ENABLE_SSL", "1")
	defer os.Setenv("NUMVERIFY_API_KEY", "")
	defer os.Setenv("NUMVERIFY_ENABLE_SSL", "")

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

	gock.New("https://apilayer.net").
		Get("/api/validate").
		MatchParam("access_key", "5ad5554ac240e4d3d31107941b35a5eb").
		MatchParam("number", number).
		Reply(200).
		JSON(expectedResult)

	s := NewNumverifySupplier()

	assert.True(t, s.IsAvailable())

	got, err := s.Validate(number)
	assert.Nil(t, err)

	assert.Equal(t, expectedResult, got)
}

func TestNumverifySupplierWithoutSSL(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "11115551212"

	_ = os.Setenv("NUMVERIFY_API_KEY", "5ad5554ac240e4d3d31107941b35a5eb")
	defer os.Setenv("NUMVERIFY_API_KEY", "")

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

	gock.New("http://apilayer.net").
		Get("/api/validate").
		MatchParam("access_key", "5ad5554ac240e4d3d31107941b35a5eb").
		MatchParam("number", number).
		Reply(200).
		JSON(expectedResult)

	s := NewNumverifySupplier()

	assert.True(t, s.IsAvailable())

	got, err := s.Validate(number)
	assert.Nil(t, err)

	assert.Equal(t, expectedResult, got)
}

func TestNumverifySupplierError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "11115551212"

	_ = os.Setenv("NUMVERIFY_API_KEY", "5ad5554ac240e4d3d31107941b35a5eb")
	defer os.Setenv("NUMVERIFY_API_KEY", "")

	expectedResult := &NumverifyValidateResponse{
		Valid: false,
		Error: numverifyError{
			Code: 100,
			Info: "Access Restricted - Your current Subscription Plan does not support HTTPS Encryption.",
		},
	}

	gock.New("http://apilayer.net").
		Get("/api/validate").
		MatchParam("access_key", "5ad5554ac240e4d3d31107941b35a5eb").
		MatchParam("number", number).
		Reply(400).
		JSON(expectedResult)

	s := NewNumverifySupplier()

	assert.True(t, s.IsAvailable())

	got, err := s.Validate(number)
	assert.Nil(t, got)
	assert.Equal(t, errors.New("Access Restricted - Your current Subscription Plan does not support HTTPS Encryption."), err)
}

func TestNumverifySupplierHTTPError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	number := "11115551212"

	_ = os.Setenv("NUMVERIFY_API_KEY", "5ad5554ac240e4d3d31107941b35a5eb")
	_ = os.Setenv("NUMVERIFY_ENABLE_SSL", "1")
	defer os.Setenv("NUMVERIFY_API_KEY", "")
	defer os.Setenv("NUMVERIFY_ENABLE_SSL", "")

	dummyError := errors.New("test")

	gock.New("https://apilayer.net").
		Get("/api/validate").
		ReplyError(dummyError)

	s := NewNumverifySupplier()

	assert.True(t, s.IsAvailable())

	got, err := s.Validate(number)
	assert.Nil(t, got)
	assert.Equal(t, &url.Error{
		Op:  "Get",
		URL: "https://apilayer.net/api/validate?access_key=5ad5554ac240e4d3d31107941b35a5eb&number=11115551212",
		Err: dummyError,
	}, err)
}

func TestNumverifySupplierWithoutAPIKey(t *testing.T) {
	s := NewNumverifySupplier()
	assert.False(t, s.IsAvailable())
}
