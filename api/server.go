package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

func registerClientRoute(router *gin.Engine, box *packr.Box) {
	router.Group("/static").
		StaticFS("/", box)

	router.GET("/", func(c *gin.Context) {
		html, err := box.Find("index.html")

		if err != nil {
			log.Fatal()
		}

		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(html))
		c.Abort()
	})
}

// Serve launches the web client
// Using Gin & Vue.js
func Serve(router *gin.Engine, port int, disableClient bool) *gin.Engine {
	httpPort := ":" + strconv.Itoa(port)

	router.Group("/api").
		GET("/", healthHandler).
		GET("/numbers", getAllNumbers).
		GET("/numbers/:number/validate", ValidateScanURL, validate).
		GET("/numbers/:number/scan/local", ValidateScanURL, localScan).
		GET("/numbers/:number/scan/numverify", ValidateScanURL, numverifyScan).
		GET("/numbers/:number/scan/googlesearch", ValidateScanURL, googleSearchScan)

	if !disableClient {
		box := packr.NewBox("../client/dist")
		registerClientRoute(router, box)
	}

	router.Use(func(c *gin.Context) {
		c.JSON(404, JSONResponse{
			Success: false,
			Error:   "Resource not found",
		})
	})

	err := router.Run(httpPort)

	if err != nil {
		log.Fatal(err)
	}

	return router
}
