package web

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"strings"
)

//go:embed client/dist/*
var clientFS embed.FS

const (
	staticPath = "/"
)

func detectContentType(path string, data []byte) string {
	arr := strings.Split(path, ".")
	ext := arr[len(arr)-1]

	if mimeType := mime.TypeByExtension(fmt.Sprintf(".%s", ext)); mimeType != "" {
		return mimeType
	}

	return http.DetectContentType(data)
}

func walkEmbededClient(dir string, router *gin.Engine) error {
	files, err := clientFS.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			err := walkEmbededClient(path.Join(dir, file.Name()), router)
			if err != nil {
				return err
			}
			continue
		}

		assetPath := strings.ReplaceAll(path.Join(dir, file.Name()), "client/dist/", staticPath)
		f, err := clientFS.Open(path.Join(dir, file.Name()))
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		if assetPath == staticPath+"index.html" {
			assetPath = staticPath
		}

		router.GET(assetPath, func(c *gin.Context) {
			c.Header("Content-Type", detectContentType(assetPath, data))
			c.Writer.WriteHeader(http.StatusOK)
			_, _ = c.Writer.Write(data)
			c.Abort()
		})
	}
	return nil
}

func registerClientRoutes(router *gin.Engine) error {
	return walkEmbededClient("client/dist", router)
}
