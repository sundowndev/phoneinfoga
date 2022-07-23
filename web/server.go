// Package web includes code for the web server of PhoneInfoga
//go:generate swag init -g ./server.go --parseDependency
package web

import (
	"github.com/gin-gonic/gin"
)

// @title PhoneInfoga REST API
// @description Advanced information gathering & OSINT framework for phone numbers.
// @version v2
// @host demo.phoneinfoga.crvx.fr
// @BasePath /api
// @schemes http https
// @license.name GNU General Public License v3.0
// @license.url https://github.com/sundowndev/phoneinfoga/blob/master/LICENSE

// Serve launches the web client
// Using Gin & Vue.js
func Serve(router *gin.Engine, disableClient bool) (*gin.Engine, error) {
	router.Group("/api").
		GET("/", healthHandler).
		GET("/numbers", getAllNumbers).
		GET("/numbers/:number/validate", ValidateScanURL, validate).
		GET("/numbers/:number/scan/local", ValidateScanURL, localScan).
		GET("/numbers/:number/scan/numverify", ValidateScanURL, numverifyScan).
		GET("/numbers/:number/scan/googlesearch", ValidateScanURL, googleSearchScan).
		GET("/numbers/:number/scan/ovh", ValidateScanURL, ovhScan)

	if !disableClient {
		err := registerClientRoute(router)
		if err != nil {
			return router, err
		}
	}

	router.Use(func(c *gin.Context) {
		c.JSON(404, JSONResponse{
			Success: false,
			Error:   "Resource not found",
		})
	})

	return router, nil
}
