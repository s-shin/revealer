package main

import (
	"path/filepath"
	"strings"
)

type Theme string

func NewTheme(theme string) Theme {
	return Theme(theme)
}

func (t Theme) IsBuiltin() bool {
	return filepath.Ext(t.String()) == ""
}

func (t Theme) IsExternal() bool {
	return strings.HasPrefix(t.String(), "http")
}

func (t Theme) IsCustom() bool {
	return !t.IsBuiltin() && !t.IsExternal()
}

func (t Theme) Href() string {
	switch {
	case t.IsBuiltin():
		return MakeBuiltinThemePath(t.String())
	case t.IsExternal():
		return t.String()
	case t.IsCustom():
		return "/" + CustomThemePath
	}
	panic("error")
}

func (t Theme) String() string {
	return string(t)
}

type Separator struct {
	Horizontal string
	Vertical   string
	Note       string
}

type Config struct {
	SlideFilePath   string
	UseBuiltinTheme bool
	Theme           Theme
	Docroot         string
	Separator       *Separator
}
