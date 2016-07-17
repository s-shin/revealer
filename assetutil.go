package main

import "errors"

// Asset path constants
const (
	BaseAssetPath       = "assets"
	BaseRevealAssetPath = "assets/reveal"
	BaseGemojiAssetPath = "assets/gemoji"
	CustomThemePath     = BaseRevealAssetPath + "/css/theme/__custom__.css"
)

func MakeBuiltinThemePath(name string) string {
	return BaseRevealAssetPath + "/css/theme/" + name + ".css"
}

func IsAvailableBuiltintTheme(name string) bool {
	_, err := AssetInfo(MakeBuiltinThemePath(name))
	return err == nil
}

func GetEmojiImageHref(name string) (string, error) {
	nonUnicodeImagePath := BaseGemojiAssetPath + "/images/emoji/" + name + ".png"
	if _, err := AssetInfo(nonUnicodeImagePath); err == nil {
		return "/" + nonUnicodeImagePath, nil
	}
	unicodeImagePath := BaseGemojiAssetPath + "/images/emoji/unicode/" + name + ".png"
	if _, err := AssetInfo(unicodeImagePath); err == nil {
		return "/" + unicodeImagePath, nil
	}
	return "", errors.New("image of the name not found: " + name)
}

//---

type RoutableAsset struct {
	Path string
}

func NewRoutableAsset(path string) *RoutableAsset {
	return &RoutableAsset{path}
}

func (a *RoutableAsset) Load() ([]byte, error) {
	return Asset(a.Path)
}

// RoutePath is a function of RoutePath interface.
func (a *RoutableAsset) RoutePath() string {
	return "/" + a.Path
}
