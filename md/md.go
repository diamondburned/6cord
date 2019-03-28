package md

import (
	"log"
	"regexp"

	bf "github.com/russross/blackfriday"
)

// HighlightStyle determines the syntax highlighting colorstyle:
// https://xyproto.github.io/splash/docs/all.html
var HighlightStyle = "vs"

var trashyCodeBlockMatching = regexp.MustCompile("(.)```")

const mdExtensions = 0 |
	bf.NoIntraEmphasis |
	bf.FencedCode |
	bf.Autolink |
	bf.Strikethrough |
	bf.NoEmptyLineBeforeBlock |
	bf.HardLineBreak

// Parse parses md into tview strings
func Parse(s string) (results string) {
	results = s
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	s = trashyCodeBlockMatching.ReplaceAllString(s, "$1\n```")
	s = fixQuotes(s)

	r := &tviewMarkdown{}
	return string(bf.Run([]byte(s),
		bf.WithNoExtensions(),
		bf.WithRenderer(r),
		bf.WithExtensions(mdExtensions),
	))
}
