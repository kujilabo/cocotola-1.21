package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed web_dist
var web embed.FS

func main() {
	// ...

	router := gin.Default()
	// WEB ---------
	viteStaticFS, err := fs.Sub(web, "web_dist")
	if err != nil {
		panic(err)
	}

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/assets") {
			c.FileFromFS(c.Request.URL.Path, http.FS(viteStaticFS))
			return
		}
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.FileFromFS("", http.FS(viteStaticFS))
			return
		}
	})
	router.Run(":8080")
}
