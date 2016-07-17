package main

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/russross/blackfriday"
)

type SlideSection struct {
	HTML string
}

func (s *SlideSection) SafeHTML() template.HTML {
	// TODO: sanitize
	return template.HTML(s.HTML)
}

func NewSlideSection(html string) *SlideSection {
	return &SlideSection{html}
}

type AggressiveMode struct {
	Config *Config
}

func NewAggressiveMode(config *Config) *AggressiveMode {
	return &AggressiveMode{config}
}

var reCodeSeparator = regexp.MustCompile("(?m)^```")

func (m *AggressiveMode) ConvertMarkdown(input string) [][]*SlideSection {
	// TODO: move these to config.js and validation
	reHorizontalSeparator := regexp.MustCompile("(?m)" + m.Config.Separator.Horizontal)
	reVerticalSeparator := regexp.MustCompile("(?m)" + m.Config.Separator.Vertical)
	reNoteSeparator := regexp.MustCompile("(?m)" + m.Config.Separator.Note)

	var result [][]*SlideSection

	horizontalChunks := splitMarkdownToChunksWithCodeBlockSafe(input, reHorizontalSeparator, -1, func(str string) string {
		return ReplaceEmoji(str, func(path string) string {
			return fmt.Sprintf("<span class='revealer-emoji' style='background-image: url(%s)'></span>", "/"+BaseGemojiAssetPath+"/"+path)
		})
	})
	for i, horizontalChunk := range horizontalChunks {
		result = append(result, make([]*SlideSection, 0))
		verticalChunks := splitMarkdownToChunksWithCodeBlockSafe(horizontalChunk, reVerticalSeparator, -1, nil)
		for _, verticalChunk := range verticalChunks {
			var sectionInput string
			chunks := splitMarkdownToChunksWithCodeBlockSafe(verticalChunk, reNoteSeparator, 2, nil)
			if len(chunks) >= 2 {
				sectionInput = chunks[0] + "<aside class='notes'>\n" + strings.Join(chunks[1:], "") + "\n</aside>"
			} else {
				sectionInput = verticalChunk
			}
			output := blackfriday.MarkdownCommon([]byte(sectionInput))
			result[i] = append(result[i], NewSlideSection(string(output)))
		}
	}

	return result
}

func splitMarkdownToChunksWithCodeBlockSafe(input string, reSeparator *regexp.Regexp, numSeparated int, filter func(string) string) []string {
	var chunks []string

	for i, chunk := range reCodeSeparator.Split(input, -1) {
		lastIdx := len(chunks) - 1
		switch i % 2 {
		case 0:
			if filter != nil {
				chunk = filter(chunk)
			}
			tmp := reSeparator.Split(chunk, numSeparated)
			if lastIdx >= 0 {
				chunks[lastIdx] += tmp[0]
			} else {
				chunks = append(chunks, tmp[0])
			}
			chunks = append(chunks, tmp[1:]...)
		case 1:
			// code content
			chunks[lastIdx] += "```" + chunk + "```"
		}
	}

	return chunks
}
