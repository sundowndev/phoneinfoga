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

	t.Run("Serve", func(t *testing.T) {
		go Serve(router, 8000, true)
		// defer srv.Shutdown()

		t.Run("healthHandler", func(t *testing.T) {
			res, err := http.Get("http://127.0.0.1:8000/api")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(err, nil, "should be equal")
			assert.Equal(res.StatusCode, 200, "should be equal")
			assert.Equal(string(body), "{\"success\":true,\"version\":\""+config.Version+"\"}\n", "should be equal")
		})

		t.Run("404 error", func(t *testing.T) {
			res, err := http.Get("http://127.0.0.1:8000/api/notfound")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(err, nil, "should be equal")
			assert.Equal(res.StatusCode, 404, "should be equal")
			assert.Equal(string(body), "{\"success\":false,\"error\":\"Resource not found\"}\n", "should be equal")
		})
	})
}
