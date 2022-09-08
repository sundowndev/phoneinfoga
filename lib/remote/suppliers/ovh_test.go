package suppliers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"gopkg.in/h2non/gock.v1"
	"net/url"
	"testing"
)

func TestOVHSupplierSuccess(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	num, _ := number.NewNumber("33365172812")

	gock.New("https://api.ovh.com").
		Get("/1.0/telephony/number/detailedZones").
		MatchParam("country", "fr").
		Reply(200).
		JSON([]OVHAPIResponseNumber{
			{
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

	s := NewOVHSupplier()

	got, err := s.Search(*num)
	assert.Nil(t, err)

	expectedResult := &OVHScannerResponse{
		Found:       true,
		NumberRange: "036517xxxx",
		City:        "Abbeville",
	}

	assert.Equal(t, expectedResult, got)
}

func TestOVHSupplierError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	num, _ := number.NewNumber("33365172812")

	dummyError := errors.New("test")

	gock.New("https://api.ovh.com").
		Get("/1.0/telephony/number/detailedZones").
		MatchParam("country", "fr").
		ReplyError(dummyError)

	s := NewOVHSupplier()

	got, err := s.Search(*num)
	assert.Nil(t, got)
	assert.Equal(t, &url.Error{
		Op:  "Get",
		URL: "https://api.ovh.com/1.0/telephony/number/detailedZones?country=fr",
		Err: dummyError,
	}, err)
}

func TestOVHSupplierCountryCodeError(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://api.ovh.com").
		Get("/1.0/telephony/number/detailedZones").
		MatchParam("country", "co").
		Reply(400).
		JSON(OVHAPIErrorResponse{Message: "[country] Given data (co) does not belong to the NumberCountryEnum enumeration"})

	num, err := number.NewNumber("+575556661212")
	if err != nil {
		t.Fatal(err)
	}

	s := NewOVHSupplier()

	got, err := s.Search(*num)
	assert.Nil(t, got)
	assert.EqualError(t, err, "[country] Given data (co) does not belong to the NumberCountryEnum enumeration")
}
