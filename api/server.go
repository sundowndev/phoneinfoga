package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

// Serve launches the web client
// Using Gin & Vue.js
func Serve(port int) *gin.Engine {
	httpPort := ":" + strconv.Itoa(port)

	router := gin.Default()

	router.Group("/api").
		GET("/", healthHandler).
		GET("/numbers", getAllNumbers).
		GET("/numbers/:number/validate", ValidateScanURL, validate).
		GET("/numbers/:number/scan/local", ValidateScanURL, localScan).
		GET("/numbers/:number/scan/numverify", ValidateScanURL, numverifyScan).
		GET("/numbers/:number/scan/googlesearch", ValidateScanURL, googleSearchScan)

	box := packr.NewBox("../client/dist")

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

		return
	})

	router.Use(func(c *gin.Context) {
		c.JSON(404, JSONResponse{
			Success: false,
			Error:   "Resource not found",
		})
	})

	router.Run(httpPort)

	return router
}
