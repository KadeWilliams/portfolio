package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	templates "myapp/templates"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	staticFiles, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/static/*", func(c echo.Context) error {
		filePath := c.Param("*")
		data, err := fs.ReadFile(staticFiles, filePath)
		if err != nil {
			return c.String(http.StatusNotFound, "File not found")
		}

		contentType := "text/plain"
		switch filepath.Ext(filePath) {
		case ".css": 
			contentType = "text/css"
		}
		return c.Blob(http.StatusOK, contentType, data)
	})

	e.GET("/", homeHandler)
	e.GET("/about", aboutHandler)
	e.GET("/projects", projectsHandler)

	e.GET("/accent-color", func(c echo.Context) error {
		colors := []string{"#ff4136", "#ff851b", "#ffdc00", "#2ecc40","#0074d9", "#b10dc9"}
		day := time.Now().Day() % len(colors)
		return c.String(http.StatusOK, colors[day])
	})
	e.GET("/server-status", func(c echo.Context) error {
		return c.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
	})
	e.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
	})

	i := 0 
	e.POST("/weather", func(c echo.Context) error {
		lat := c.FormValue("latitude")
		long := c.FormValue("longitude")
		fmt.Printf("\nlat: %s\nlong: %s\n", lat, long)
		i++
		time.Sleep(1 * time.Second)
		return c.String(http.StatusOK, fmt.Sprintf("%d", i))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":"+port))
}

func homeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Home())
}

func aboutHandler(c echo.Context) error {
	return Render(c, http.StatusOK, templates.About())
}

func projectsHandler(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Projects())
}