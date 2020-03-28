package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/pkg/scanners"
	gock "gopkg.in/h2non/gock.v1"
)

func performRequest(r http.Handler, method, path string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, err
}

func TestApi(t *testing.T) {
	assert := assert.New(t)
	r := gin.Default()
	r = Serve(r, true)

	t.Run("Serve", func(t *testing.T) {
		t.Run("getAllNumbers - /api/numbers", func(t *testing.T) {
			res, err := performRequest(r, "GET", "/api/numbers")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(err, nil, "should be equal")
			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(string(body), "{\"success\":true,\"error\":\"\",\"numbers\":[]}", "should be equal")
		})

		t.Run("validate - /api/numbers/:number/validate", func(t *testing.T) {
			t.Run("valid number", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/3312345253/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 200, "should be equal")
				assert.Equal(string(body), "{\"success\":true,\"error\":\"The number is valid\"}", "should be equal")
			})

			t.Run("invalid number", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/azerty/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 400, "should be equal")
				assert.Equal(string(body), "{\"success\":false,\"error\":\"Parameter 'number' must be a valid integer.\"}", "should be equal")
			})

			t.Run("invalid country code", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/09880/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 400, "should be equal")
				assert.Equal(string(body), "{\"success\":false,\"error\":\"invalid country code\"}", "should be equal")
			})
		})

		t.Run("localScan - /api/numbers/:number/scan/local", func(t *testing.T) {
			t.Run("valid number", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/3312345253/scan/local")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 200, "should be equal")
				assert.Equal(string(body), `{"success":true,"error":"","result":{"rawLocal":"12345253","local":"12345253","E164":"+3312345253","international":"3312345253","countryCode":33,"country":"FR","carrier":""}}`, "should be equal")
			})

			t.Run("invalid number", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/9999999999/scan/local")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 400, "should be equal")
				assert.Equal(string(body), `{"success":false,"error":"invalid country code"}`, "should be equal")
			})
		})

		t.Run("numverifyScan - /api/numbers/:number/scan/numverify", func(t *testing.T) {
			t.Run("should succeed", func(t *testing.T) {
				defer gock.Off() // Flush pending mocks after test execution

				expectedResult := scanners.NumverifyScannerResponse{
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

				gock.New("http://numverify.com").
					Get("/").
					Reply(200).BodyString(`<html><body><input type="hidden" name="scl_request_secret" value="secret"/></body></html>`)

				gock.New("https://numverify.com").
					Get("/php_helper_scripts/phone_api.php").
					MatchParam("secret_key", "5ad5554ac240e4d3d31107941b35a5eb").
					MatchParam("number", "79516566591").
					Reply(200).
					JSON(expectedResult)

				res, err := performRequest(r, "GET", "/api/numbers/79516566591/scan/numverify")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 200, "should be equal")
				assert.Equal(string(body), `{"success":true,"error":"","result":{"valid":true,"number":"79516566591","local_format":"9516566591","international_format":"+79516566591","country_prefix":"+7","country_code":"RU","country_name":"Russian Federation","location":"Saint Petersburg and Leningrad Oblast","carrier":"OJSC St. Petersburg Telecom (OJSC Tele2-Saint-Petersburg)","line_type":"mobile"}}`, "should be equal")

				assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
			})
		})

		t.Run("ovhScan - /api/numbers/:number/scan/ovh", func(t *testing.T) {
			t.Run("should find number on OVH", func(t *testing.T) {
				defer gock.Off() // Flush pending mocks after test execution

				gock.New("https://api.ovh.com").
					Get("/1.0/telephony/number/detailedZones").
					MatchParam("country", "fr").
					Reply(200).
					JSON([]scanners.OVHAPIResponseNumber{
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

				res, err := performRequest(r, "GET", "/api/numbers/330365179268/scan/ovh")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 200, "should be equal")
				assert.Equal(string(body), `{"success":true,"error":"","result":{"found":true,"numberRange":"036517xxxx","city":"Abbeville","zipCode":""}}`, "should be equal")

				assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
			})
		})

		// t.Run("healthHandler - /api/", func(t *testing.T) {
		// 	res, err := performRequest(r, "GET", "/api")

		// 	body, _ := ioutil.ReadAll(res.Body)

		// 	assert.Equal(err, nil, "should be equal")
		// 	assert.Equal(res.Result().StatusCode, 200, "should be equal")
		// 	assert.Equal(string(body), "{\"success\":true,\"version\":\""+config.Version+"\"}", "should be equal")
		// })

		t.Run("404 error - /api/notfound", func(t *testing.T) {
			res, err := performRequest(r, "GET", "/api/notfound")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(err, nil, "should be equal")
			assert.Equal(res.Result().StatusCode, 404, "should be equal")
			assert.Equal(string(body), "{\"success\":false,\"error\":\"Resource not found\"}", "should be equal")
		})
	})
}
