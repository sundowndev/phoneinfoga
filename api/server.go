// Package api is about the REST API of PhoneInfoga
//go:generate cp -r ../client/dist ./static
//go:generate go run github.com/swaggo/swag/cmd/swag init -g ./server.go --parseDependency
package api

import (
	"embed"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	staticPath = "/"
)

//go:embed static
var clientFiles embed.FS

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

func walkDir(dir string, efs embed.FS) ([]string, error) {
	var arr []string
	files, err := efs.ReadDir(dir)
	if err != nil {
		return arr, err
	}

	for _, file := range files {
		if file.IsDir() {
			filesInDir, err := walkDir(filepath.Join(dir, file.Name()), efs)
			if err != nil {
				return arr, err
			}
			arr = append(arr, filesInDir...)
			continue
		}
		arr = append(arr, filepath.Join(dir, file.Name()))
	}

	return arr, nil
}

func registerClientRoute(router *gin.Engine) error {
	files, err := walkDir("static", clientFiles)
	if err != nil {
		return err
	}

	for _, name := range files {
		path := strings.ReplaceAll(name, "static/", staticPath)
		data, err := clientFiles.ReadFile(name)
		if err != nil {
			return err
		}

		if path == staticPath+"index.html" {
			path = staticPath
		}

		router.GET(path, func(c *gin.Context) {
			c.Header("Content-Type", detectContentType(path, data))
			c.Writer.WriteHeader(http.StatusOK)
			_, _ = c.Writer.Write(data)
			c.Abort()
		})
	}
	return nil
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
		err := registerClientRoute(router)
		if err != nil {
			panic(err)
		}
	}

	router.Use(func(c *gin.Context) {
		c.JSON(404, JSONResponse{
			Success: false,
			Error:   "Resource not found",
		})
	})

	return router
}
