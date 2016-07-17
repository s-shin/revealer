package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode/utf8"
)

const templateEmojiGo = `package main

import "regexp"

var emojiMapping = map[string]string{ {{range $k, $v := .EmojiMapping}}
	"{{$k}}": "{{$v}}",{{end}}
}

var reEmoji = regexp.MustCompile(":({{.AliasesStr}}):")

func ReplaceEmoji(str string, replace func(string) string) string {
	return reEmoji.ReplaceAllStringFunc(str, func(m string) string {
		parts := reEmoji.FindStringSubmatch(m)
		return replace(emojiMapping[parts[1]])
	})
}
`

type emoji struct {
	Emoji       string
	Description string
	Aliases     []string
	Tags        []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <gemoji-dir>\n", os.Args[0])
		os.Exit(1)
	}
	gemojiDir := os.Args[1]
	fd, err := os.Open(filepath.Join(gemojiDir, "db/emoji.json"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	decoder := json.NewDecoder(fd)
	var emojis []emoji
	if err := decoder.Decode(&emojis); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fd.Close()
	emojiMapping := make(map[string]string)
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
			imagePath := "images/emoji/" + fileName + ".png"
			if _, err := os.Stat(filepath.Join(gemojiDir, imagePath)); err != nil {
				imagePath = "images/emoji/unicode/" + fileName + ".png"
				if _, err := os.Stat(filepath.Join(gemojiDir, imagePath)); err != nil {
					continue
				}
			}
			emojiMapping[alias] = imagePath
		}
		t := template.Must(template.New("emoji.go").Parse(templateEmojiGo))
		fd, err := os.OpenFile("emoji.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		aliases := make([]string, 0, len(emojiMapping))
		for alias := range emojiMapping {
			aliases = append(aliases, strings.Replace(regexp.QuoteMeta(alias), "\\", "\\\\", -1))
		}
		t.Execute(fd, &struct {
			EmojiMapping map[string]string
			AliasesStr   string
		}{
			EmojiMapping: emojiMapping,
			AliasesStr:   strings.Join(aliases, "|"),
		})
		fd.Close()
	}
}
