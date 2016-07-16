package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

// RunServer runs server with options.
func RunServer(config *Config) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		data, err := Asset("assets/index.html")
		if err != nil {
			return err
		}
		t, err := template.New("index").Parse(string(data[:]))
		if err != nil {
			return err
		}
		return t.Execute(c.Response().Writer(), config)
	})

	emojiReplacer := mustCreateReplacerForEmojiInMarkdown()

	e.GET("/__slide__.md", func(c echo.Context) error {
		data, err := ioutil.ReadFile(config.SlideFilePath)
		if err != nil {
			return err
		}
		_, err = emojiReplacer.WriteString(c.Response().Writer(), string(data[:]))
		return err
	})

	if config.Theme.IsCustom() {
		e.GET(config.Theme.Href(), func(c echo.Context) error {
			return respondStatic(config.Theme.String(), c)
		})
	}

	e.GET("/"+BaseRevealAssetPath+"/*", func(c echo.Context) error {
		path := c.Request().URL().Path()
		fmt.Println(path)
		return respondRevealAsset(path[1:], c)
	})

	e.Static("/", config.Docroot)

	e.Run(standard.New(":3000"))
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

func mustCreateReplacerForEmojiInMarkdown() *strings.Replacer {
	replacer, err := createReplacerForEmojiInMarkdown()
	if err != nil {
		panic(err)
	}
	return replacer
}

func createReplacerForEmojiInMarkdown() (*strings.Replacer, error) {
	data, err := Asset(BaseGemojiAssetPath + "/db/emoji.json")
	if err != nil {
		return nil, err
	}
	type Emoji struct {
		Emoji       string
		Description string
		Aliases     []string
		Tags        []string
	}
	var emojis []Emoji
	if err := json.Unmarshal(data, &emojis); err != nil {
		panic(err)
	}
	oldnew := make([]string, 0, len(emojis)*2)
	for _, emoji := range emojis {
		fileName := ""
		if len(emoji.Emoji) > 0 {
			r, _ := utf8.DecodeRuneInString(emoji.Emoji)
			fileName = fmt.Sprintf("%x", r)
		}
		for _, alias := range emoji.Aliases {
			if fileName == "" {
				fileName = alias
			}
			href, err := GetEmojiImageHref(fileName)
			if err != nil {
				continue
			}
			oldnew = append(
				oldnew,
				":"+alias+":",
				fmt.Sprintf("<img src='%s' />", href),
			)
		}
	}
	return strings.NewReplacer(oldnew...), nil
}
