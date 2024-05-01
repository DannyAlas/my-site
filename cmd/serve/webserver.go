package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	// print cwd
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("cwd: %s", cwd)
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	GenerateRoutes(e, "output")
	e.Static("/", "output/static")

	// Custom routes page
	e.GET("/routes", func(c echo.Context) error {
		routes_str := ""
		for _, route := range e.Routes() {
			routes_str += route.Path + "\n"
		}
		return c.String(http.StatusOK, routes_str)
	})

	// Start server
	e.Logger.Fatal(e.Start(":42069"))
}
