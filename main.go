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
	app.Usage = "instant presentation server for markdown by reveal.js"
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
		// TODO: config file support
		// cli.StringFlag{
		// 	Name:  "config, c",
		// 	Value: "",
		// 	Usage: "",
		// },
		cli.IntFlag{
			Name:  "port, p",
			Value: 3000,
			Usage: "port number for http server",
		},
		cli.StringFlag{
			Name:  "docroot, d",
			Value: "",
			Usage: "default is base directory of <slide.md>",
		},
		cli.StringFlag{
			Name:  "mode, m",
			Value: "normal",
			Usage: "normal or aggressive",
		},
		cli.StringFlag{
			Name:  "theme, t",
			Value: "black",
			Usage: "builtin theme name, external theme url, or custom theme path",
		},
		cli.StringFlag{
			Name:  "theme-type",
			Value: "auto",
			Usage: "builtin, external, custom, or auto",
		},
		cli.StringFlag{
			Name:  "separator-horizontal",
			Value: "^---",
			Usage: "the value of data-separator attribute of reveal.js",
		},
		cli.StringFlag{
			Name:  "separator-vertical",
			Value: "^___",
			Usage: "the value of data-separator-vertical attribute of reveal.js",
		},
		cli.StringFlag{
			Name:  "separator-note",
			Value: "^Note:",
			Usage: "the value of data-separator-notes attribute of reveal.js",
		},
	}

	app.Action = action

	app.Run(os.Args)
}

func action(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		fmt.Printf("ERROR: one argument required.\n\n")
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	slideFilePath := args[0]
	theme := NewTheme(c.String("theme"), NewThemeType(c.String("theme-type")))
	docroot := c.String("docroot")
	if docroot == "" {
		docroot = filepath.Dir(slideFilePath)
	}

	config := &Config{
		SlideFilePath: slideFilePath,
		Port:          c.Int("port"),
		Docroot:       docroot,
		Mode:          NewMode(c.String("mode")),
		Theme:         theme,
		Separator: Separator{
			Horizontal: c.String("separator-horizontal"),
			Vertical:   c.String("separator-vertical"),
			Note:       c.String("separator-note"),
		},
	}

	if err := Validator(config).CheckValidity(); err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(1)
	}

	RunServer(config)
}
