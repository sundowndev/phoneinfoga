// Package web includes code for the web server of PhoneInfoga
//
//go:generate swag init -g ./server.go --parseDependency
package web

import (
	"github.com/gin-gonic/gin"
	v2 "github.com/sundowndev/phoneinfoga/v2/web/v2/api/server"
	"net/http"
)

// @title PhoneInfoga REST API
// @description Advanced information gathering & OSINT framework for phone numbers.
// @version v2
// @host localhost:5000
// @BasePath /api
// @schemes http https
// @license.name GNU General Public License v3.0
// @license.url https://github.com/sundowndev/phoneinfoga/blob/master/LICENSE

type Server struct {
	router *gin.Engine
}

func NewServer(disableClient bool) (*Server, error) {
	s := &Server{
		router: gin.Default(),
	}
	if err := s.registerRoutes(disableClient); err != nil {
		return s, err
	}
	return s, nil
}

func (s *Server) registerRoutes(disableClient bool) error {
	group := s.router.Group("/api")

	group.
		GET("/", healthHandler).
		GET("/numbers", getAllNumbers).
		GET("/numbers/:number/validate", ValidateScanURL, validate).
		GET("/numbers/:number/scan/local", ValidateScanURL, localScan).
		GET("/numbers/:number/scan/numverify", ValidateScanURL, numverifyScan).
		GET("/numbers/:number/scan/googlesearch", ValidateScanURL, googleSearchScan).
		GET("/numbers/:number/scan/ovh", ValidateScanURL, ovhScan)

	v2routes := v2.NewServer().Routes()
	for _, r := range v2routes {
		group.Handle(r.Method, r.Path, r.HandlerFunc)
	}

	if !disableClient {
		err := registerClientRoutes(s.router)
		if err != nil {
			return err
		}
	}

	s.router.Use(func(c *gin.Context) {
		c.JSON(404, JSONResponse{
			Success: false,
			Error:   "resource not found",
		})
	})

	return nil
}

func (s *Server) ListenAndServe(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
