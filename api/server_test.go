package api

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/pkg/config"
)

func TestApi(t *testing.T) {
	assert := assert.New(t)
	router := gin.Default()

	apiURL := "http://127.0.0.1:8000"

	t.Run("Serve", func(t *testing.T) {
		go Serve(router, 8000, true)
		// defer srv.Shutdown()

		t.Run("getAllNumbers - /api/numbers", func(t *testing.T) {
			res, err := http.Get(apiURL + "/api/numbers")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(err, nil, "should be equal")
			assert.Equal(res.StatusCode, 200, "should be equal")
			assert.Equal(string(body), "{\"success\":true,\"error\":\"\",\"numbers\":[]}\n", "should be equal")
		})

		t.Run("validate - /api/numbers/:number/validate", func(t *testing.T) {
			t.Run("valid number", func(t *testing.T) {
				res, err := http.Get(apiURL + "/api/numbers/3312345253/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.StatusCode, 200, "should be equal")
				assert.Equal(string(body), "{\"success\":true,\"error\":\"The number is valid\"}\n", "should be equal")
			})

			t.Run("invalid number", func(t *testing.T) {
				res, err := http.Get(apiURL + "/api/numbers/azerty/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.StatusCode, 400, "should be equal")
				assert.Equal(string(body), "{\"success\":false,\"error\":\"Parameter 'number' must be a valid integer.\"}\n", "should be equal")
			})

			t.Run("invalid country code", func(t *testing.T) {
				res, err := http.Get(apiURL + "/api/numbers/09880/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.StatusCode, 400, "should be equal")
				assert.Equal(string(body), "{\"success\":false,\"error\":\"invalid country code\"}\n", "should be equal")
			})
		})

		t.Run("healthHandler - /api/", func(t *testing.T) {
			res, err := http.Get(apiURL + "/api")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(err, nil, "should be equal")
			assert.Equal(res.StatusCode, 200, "should be equal")
			assert.Equal(string(body), "{\"success\":true,\"version\":\""+config.Version+"\"}\n", "should be equal")
		})

		t.Run("404 error - /api/notfound", func(t *testing.T) {
			res, err := http.Get(apiURL + "/api/notfound")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(err, nil, "should be equal")
			assert.Equal(res.StatusCode, 404, "should be equal")
			assert.Equal(string(body), "{\"success\":false,\"error\":\"Resource not found\"}\n", "should be equal")
		})
	})
}
