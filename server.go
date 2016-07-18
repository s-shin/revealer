package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

type RoutePath interface {
	RoutePath() string
}

// RunServer runs server with options.
func RunServer(config *Config) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		data, err := config.Mode.TemplateAsset().Load()
		if err != nil {
			panic(err)
		}
		t, err := template.New(config.Mode.String()).Parse(string(data[:]))
		if err != nil {
			return err
		}
		switch config.Mode {
		case ModeNormal:
			return t.Execute(c.Response().Writer(), config)
		case ModeAggressive:
			markdown := mustLoadSlideMarkdown(config.SlideFilePath)
			sections := NewAggressiveMode(config).ConvertMarkdown(markdown)
			return t.Execute(c.Response().Writer(), &struct {
				Config        *Config
				SlideSections [][]*SlideSection
			}{
				Config:        config,
				SlideSections: sections,
			})
		}
		panic("error")
	})

	if config.Mode == ModeNormal {
		e.GET("/__slide__.md", func(c echo.Context) error {
			markdown := mustLoadSlideMarkdown(config.SlideFilePath)
			return c.String(http.StatusOK, markdown)
		})
	}

	if config.Theme.Type == ThemeTypeCustom {
		e.GET(config.Theme.RoutePath(), func(c echo.Context) error {
			return respondStatic(config.Theme.String(), c)
		})
	}

	e.GET("/"+BaseRevealAssetPath+"/*", func(c echo.Context) error {
		path := c.Request().URL().Path()
		return respondRevealAsset(path[1:], c)
	})

	e.Static("/", config.Docroot)

	fmt.Printf("Server started. Let's open http://localhost:%d/ in browser.\n", config.Port)
	e.Run(standard.New(fmt.Sprintf(":%d", config.Port)))
}

func mustLoadSlideMarkdown(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func respondStatic(path string, c echo.Context) error {
	fd, err := os.Open(path)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	fi, err := fd.Stat()
	if err != nil {
		return err
	}
	return c.ServeContent(fd, fi.Name(), fi.ModTime())
}

func respondRevealAsset(path string, c echo.Context) error {
	data, err := Asset(path)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	res := c.Response()
	res.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
	res.Write(data)
	return nil
}
