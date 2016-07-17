package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//---

type Validator interface {
	CheckValidity() error
}

//---

type ThemeType string

// ThemeType constants
const (
	ThemeTypeNil      = ThemeType("")
	ThemeTypeBuiltin  = ThemeType("builtin")
	ThemeTypeCustom   = ThemeType("custom")
	ThemeTypeExternal = ThemeType("external")
	ThemeTypeAuto     = ThemeType("auto")
)

func NewThemeType(t string) ThemeType {
	return ThemeType(t)
}

type Theme struct {
	Value string
	Type  ThemeType
}

func NewTheme(v string, t ThemeType) Theme {
	if t == ThemeTypeAuto {
		t = detectThemeType(v)
	}
	return Theme{
		Value: v,
		Type:  t,
	}
}

func detectThemeType(v string) ThemeType {
	if filepath.Ext(v) == "" {
		return ThemeTypeBuiltin
	}
	if strings.HasPrefix(v, "http") {
		return ThemeTypeExternal
	}
	return ThemeTypeCustom
}

// CheckValidity is a function of Validator interface.
func (t Theme) CheckValidity() error {
	switch t.Type {
	case ThemeTypeBuiltin:
		_, err := AssetInfo(MakeBuiltinThemePath(t.Value))
		return err
	case ThemeTypeCustom:
		_, err := os.Stat(t.Value)
		return err
	}
	return nil
}

// RoutePath is a function of RoutePath interface.
func (t Theme) RoutePath() string {
	switch t.Type {
	case ThemeTypeBuiltin:
		return "/" + MakeBuiltinThemePath(t.Value)
	case ThemeTypeExternal:
		return t.Value
	case ThemeTypeCustom:
		return "/" + CustomThemePath
	}
	panic("error")
}

func (t Theme) String() string {
	return t.Value
}

//---

type Mode string

// Mode constants
const (
	ModeNormal     = Mode("normal")
	ModeAggressive = Mode("aggressive")
)

var availableModes = map[Mode]bool{
	ModeNormal:     true,
	ModeAggressive: true,
}

func NewMode(mode string) Mode {
	return Mode(mode)
}

// CheckValidity checks whether this is valid or not.
func (mode Mode) CheckValidity() error {
	_, ok := availableModes[mode]
	if !ok {
		return errors.New("Unavailable mode: " + mode.String())
	}
	return nil
}

func (mode Mode) String() string {
	return string(mode)
}

func (mode Mode) TemplateAsset() *RoutableAsset {
	switch mode {
	case ModeNormal, ModeAggressive:
		return NewRoutableAsset(BaseAssetPath + "/" + mode.String() + ".html")
	}
	panic("TODO")
}

//---

type Separator struct {
	Horizontal string
	Vertical   string
	Note       string
}

//---

type ConfigError struct {
	Message string
}

func NewConfigError(msg string) *ConfigError {
	return &ConfigError{msg}
}

func (e *ConfigError) Error() string {
	return e.Message
}

//---

type Config struct {
	SlideFilePath string
	Port          int
	Docroot       string
	Mode          Mode
	Theme         Theme
	Separator     Separator
}

// CheckValidity is a function of Validator interface.
func (c *Config) CheckValidity() error {
	if _, err := os.Stat(c.SlideFilePath); err != nil {
		return NewConfigError(fmt.Sprintf("file not found: %s\n", c.SlideFilePath))
	}
	if _, err := os.Stat(c.Docroot); err != nil {
		return NewConfigError(fmt.Sprintf("docroot not found: %s\n", c.Docroot))
	}
	// TODO: check more
	return nil
}
