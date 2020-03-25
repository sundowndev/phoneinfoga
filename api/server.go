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
func Serve(router *gin.Engine, port int, disableClient bool) (*gin.Engine, *http.Server) {
	httpPort := ":" + strconv.Itoa(port)

	router.Group("/api").
		GET("/", healthHandler).
		GET("/numbers", getAllNumbers).
		GET("/numbers/:number/validate", ValidateScanURL, validate).
		GET("/numbers/:number/scan/local", ValidateScanURL, localScan).
		GET("/numbers/:number/scan/numverify", ValidateScanURL, numverifyScan).
		GET("/numbers/:number/scan/googlesearch", ValidateScanURL, googleSearchScan).
		GET("/numbers/:number/scan/ovh", ValidateScanURL, ovhScan)

	if !disableClient {
		box := packr.New("static", "../client/dist")
		registerClientRoute(router, box)
	}

	router.Use(func(c *gin.Context) {
		c.JSON(404, JSONResponse{
			Success: false,
			Error:   "Resource not found",
		})
	})

	srv := &http.Server{
		Addr:    httpPort,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	return router, srv
}
