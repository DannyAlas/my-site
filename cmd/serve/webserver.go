package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GenerateRoutes(e *echo.Echo, dir string) {
	// walk the directory and add routes for each file
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			relPath, err := filepath.Rel(dir, path)
			// remmove the /pages/ prefix
			if err != nil {
				return err
			}
			// if we're in the root of the pages directory, remove the prefix
			if strings.HasPrefix(relPath, "pages/") {
				relPath = strings.TrimPrefix(relPath, "pages/")
			}

			// if it's an index file, set the route to the directory
			if strings.Contains(relPath, "index.html") {
				relPath = strings.TrimSuffix(relPath, "index.html")
				relPath = strings.TrimSuffix(relPath, "posts/")

			} else {
				relPath = strings.TrimSuffix(relPath, ".html")
			}

			// \ -> /
			relPath = strings.ReplaceAll(relPath, "\\", "/")
			// remove trailing slash
			relPath = strings.TrimSuffix(relPath, "/")
			e.GET("/"+relPath, func(c echo.Context) error {
				return c.File(path)
			})
		}

		return nil
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func main() {

	// Parse flags
	var distDir string
	flag.StringVar(&distDir, "dist", "dist", "The directory to serve files from")
	flag.Parse()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	GenerateRoutes(e, distDir)
	e.Static("/", distDir+"/static")
	// server favicom in dist/static/imgs
	e.File("/favicon.ico", distDir+"/static/imgs/favicon.ico")
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
