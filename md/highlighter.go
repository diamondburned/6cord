package md

import (
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// HighlightStyle determines the syntax highlighting colorstyle:
// https://xyproto.github.io/splash/docs/all.html
var HighlightStyle = "vs"

var (
	style  *chroma.Style
	fmtter chroma.Formatter
)

// RenderCodeBlock renders the node to a syntax
// highlighted code
func RenderCodeBlock(lang, content string) string {
	// Get the highlight style
	if style == nil {
		// Try and find it
		if s := styles.Get(HighlightStyle); s != nil {
			style = s
		} else {
			// Can't find, use fallback
			style = styles.Fallback
		}
	}

	if fmtter == nil {
		fmtter = formatters.Get("tview-256bit")
	}

	var lexer = lexers.Fallback
	if lang := string(lang); lang != "" {
		if l := lexers.Get(lang); l != nil {
			lexer = l
		} else {
			content = lang + "\n" + content
		}
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return wrapBlock(content)
	}

	var code strings.Builder

	if err := fmtter.Format(&code, style, iterator); err != nil {
		return wrapBlock(content)
	}

	return wrapBlock(code.String())
}

func wrapBlock(content string) string {
	var s strings.Builder

	// wrapped := tview.WordWrap(code.String(), 80)
	wrapped := strings.Split(content, "\n")

	for i := 0; i < len(wrapped); i++ {
		if wrapped[i] != "[-]" {
			s.WriteString("\n[grey]┃[-] " + wrapped[i])
		}
	}

	return s.String()[1:] + "\n"
}
