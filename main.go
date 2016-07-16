package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "revealgo"
	app.Usage = "simple server to show markdown as slide by reveal.js"
	app.Version = "v0.1.0"

	cli.AppHelpTemplate = `NAME:
     {{.Name}} - {{.Usage}}

USAGE:
    {{.HelpName}} [options] {{.ArgsUsage}}

OPTIONS:
    {{range .VisibleFlags}}{{.}}
    {{end}}
VERSION:
    {{.Version}}

`

	app.ArgsUsage = "<slide.md>"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 3000,
			Usage: "port number for http server",
		},
		cli.StringFlag{
			Name:  "theme, t",
			Value: "black",
			Usage: "",
		},
		cli.StringFlag{
			Name:  "docroot, d",
			Value: "",
			Usage: "default is base directory of <slide.md>",
		},
		cli.StringFlag{
			Name:  "separator-horizontal",
			Value: "^---",
			Usage: "",
		},
		cli.StringFlag{
			Name:  "separator-vertical",
			Value: "^___",
			Usage: "",
		},
		cli.StringFlag{
			Name:  "separator-note",
			Value: "^Note:",
			Usage: "",
		},
	}

	app.Action = action

	app.Run(os.Args)
}

func action(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		fmt.Printf("one argument required.\n\n")
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	slideFilePath := args[0]
	if _, err := os.Stat(slideFilePath); err != nil {
		fmt.Printf("file not found: %s\n", slideFilePath)
		os.Exit(1)
	}

	theme := NewTheme(c.String("theme"))
	if theme.IsBuiltin() {
		if !IsAvailableBuiltintTheme(theme.String()) {
			fmt.Printf("theme not included: %s\n", theme)
			os.Exit(1)
		}
	} else if theme.IsCustom() {
		if _, err := os.Stat(theme.String()); err != nil {
			fmt.Printf("file not found: %s\n", theme)
			os.Exit(1)
		}
	}

	docroot := c.String("docroot")
	if docroot == "" {
		docroot = filepath.Dir(slideFilePath)
	}

	RunServer(&Config{
		SlideFilePath: slideFilePath,
		Theme:         theme,
		Docroot:       docroot,
		Separator: &Separator{
			Horizontal: c.String("separator"),
			Vertical:   c.String("separator-vertical"),
			Note:       c.String("separator-note"),
		},
	})
}
