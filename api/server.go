// Package api is about the REST API of PhoneInfoga
//go:generate go run github.com/jessevdk/go-assets-builder ../client/dist -o ./assets.go -p api
//go:generate go run github.com/swaggo/swag/cmd/swag init -g ./server.go --parseDependency
package api

import (
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	clientDistPath = "/client/dist/"
	staticPath     = "/"
)

// @title PhoneInfoga REST API
// @description Advanced information gathering & OSINT framework for phone numbers.
// @version v2
// @host demo.phoneinfoga.crvx.fr
// @BasePath /api
// @schemes http https
// @license.name GNU General Public License v3.0
// @license.url https://github.com/sundowndev/phoneinfoga/blob/master/LICENSE

func detectContentType(path string, data []byte) string {
	arr := strings.Split(path, ".")
	ext := arr[len(arr)-1]

	if mimeType := mime.TypeByExtension(fmt.Sprintf(".%s", ext)); mimeType != "" {
		return mimeType
	}

	return http.DetectContentType(data)
}

func registerClientRoute(router *gin.Engine) {
	for name, file := range Assets.Files {
		if file.IsDir() {
			continue
		}

		path := strings.ReplaceAll(name, clientDistPath, staticPath)
		data := file.Data

		if path == staticPath+"index.html" {
			path = staticPath
		}

		router.GET(path, func(c *gin.Context) {
			c.Header("Content-Type", detectContentType(path, data))
			c.Writer.WriteHeader(http.StatusOK)
			c.Writer.Write(data)
			c.Abort()
		})
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
