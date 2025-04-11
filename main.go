package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	templates "myapp/templates"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RequestData struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
//go:embed static/*
var staticFS embed.FS

func main() {
	// expose Go functions to JS 
	// js.Global().Set("goSort", js.FuncOf(sort))//<-make(chan bool)
	// js.Global().Set("goSort", js.FuncOf(sort))<-make(chan bool) // keep alive
	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

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
	e.GET("/project", projectHandler)
	// e.GET("/status", sysInfoHandler)
	// e.GET("/wasm", serveWASM)
	// e.GET("/wasm_exec.js", serveWASMJS)
		

	e.GET("/accent-color", func(c echo.Context) error {
		colors := []string{"#ff4136", "#ff851b", "#ffdc00", "#2ecc40","#0074d9", "#b10dc9"}
		day := time.Now().Day() % len(colors)
		return c.String(http.StatusOK, colors[day])
	})
	e.GET("/server-status", func(c echo.Context) error {
		return c.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
	})

	e.GET("/weather", func(c echo.Context) error {
		lat := c.FormValue("latitude")
		long := c.FormValue("longitude")
		// resp, err := http.Get(fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&hourly=temperature_2m&current=relative_humidity_2m,temperature_2m&temperature_unit=fahrenheit", lat, long))
		// if err != nil {
		// 	return c.String(http.StatusBadRequest, err.Error())
		// }
		// defer resp.Body.Close()
		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	return c.String(http.StatusBadRequest, err.Error())
		// }
		// fmt.Println(string(body))
		fmt.Println(lat, long)
		return c.String(http.StatusOK, "75")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":"+port))
}

// func sort(this js.Value, args []js.Value) interface{} {
// 	arr := generateRandomArray(10) // 10 elements
// 	bubbleSort(arr, js.Global().Get("updateSortBars"))
// 	return nil
// }

// func bubbleSort(arr []int, callback js.Value) {
// 	n := len(arr)
// 	for i := 0; i < n-1; i++ {
// 		for j := 0; j < n-i-1; j++ {
// 			if arr[j] > arr[j+1] {
// 				arr[j], arr[j+1] = arr[j+1], arr[j]
// 			}
// 			// update visualization after each step

// 			callback.Invoke(js.ValueOf(arr))
// 			time.Sleep(100 * time.Millisecond) // for animation
// 		}
// 	}
// }

func generateRandomArray(size int) []int {
	arr := make([]int, size)
	for i := range arr {
		arr[i] = rand.Intn(100) + 1 // 1-100
	}
	return arr 
}

func serveWASM(c echo.Context) error {
	return c.File("static/wasm/main.wasm")
}

func serveWASMJS(c echo.Context) error {
	return c.File("static/wasm/wasm_exec.js")
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

func projectHandler(c echo.Context) error {
	// Get Projects somehow
	return Render(c, http.StatusOK, templates.Project("example"))
}

func sysInfoHandler(c echo.Context) error {
	ipAddress := c.RealIP()
	return Render(c, http.StatusOK, templates.SystemInfo(ipAddress))
}

func wasmHandler(c echo.Context) error {
	return c.File("static/wasm/main.wasm")
}

func wasmJsHandler(c echo.Context) error {
	return c.File("static/wasm/wasm_exec.js")
}

