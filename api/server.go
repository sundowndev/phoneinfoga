package api

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerClientRoute(router *gin.Engine) {
	for name, file := range Assets.Files {
		if !file.IsDir() {
			println(111111, name, file.Path)

			h, _ := ioutil.ReadAll(file)

			router.GET(strings.ReplaceAll(name, "/client/dist", "/"), func(c *gin.Context) {
				// c.Header("content-type", getMimeType(file.Path))
				c.Writer.WriteHeader(http.StatusOK)
				c.Writer.Write([]byte(h))
				c.Abort()
			})
		}
	}
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
