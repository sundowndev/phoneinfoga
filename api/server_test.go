package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
				assert.Equal(res.Result().StatusCode, 500, "should be equal")
				assert.Equal(string(body), `{"success":false,"error":"invalid country code"}`, "should be equal")
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
