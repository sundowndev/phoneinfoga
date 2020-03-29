//go:generate $GOPATH/bin/pkger

package api

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
)

func registerClientRoute(router *gin.Engine) {
	router.StaticFS("/static", pkger.Dir("/client/dist"))

	router.GET("/", func(c *gin.Context) {
		f, _ := pkger.Open("/client/dist/index.html")
		sl, _ := ioutil.ReadAll(f)

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(sl))
		c.Abort()
	})
}

// Serve launches the web client
// Using Gin & Vue.js
func Serve(router *gin.Engine, disableClient bool) *gin.Engine {
	router.Group("/api").
		GET("/", healthHandler).
		GET("/numbers", getAllNumbers).
		GET("/numbers/:number/validate", ValidateScanURL, validate).
		GET("/numbers/:number/scan/local", ValidateScanURL, localScan).
		GET("/numbers/:number/scan/numverify", ValidateScanURL, numverifyScan).
		GET("/numbers/:number/scan/googlesearch", ValidateScanURL, googleSearchScan).
		GET("/numbers/:number/scan/ovh", ValidateScanURL, ovhScan)

	if !disableClient {
		pkger.Include("/client/dist")
		registerClientRoute(router)
	}

	router.Use(func(c *gin.Context) {
		c.JSON(404, JSONResponse{
			Success: false,
			Error:   "Resource not found",
		})
	})

	return router
}
