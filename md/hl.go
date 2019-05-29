package md

import (
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/rivo/tview"
)

// RenderCodeBlock renders the node to a syntax
// highlighted code
func RenderCodeBlock(lang, literal []byte) string {
	var s strings.Builder

	content := strings.TrimFunc(
		string(literal),
		func(r rune) bool {
			return r == '\n'
		},
	)

	var lexer = lexers.Fallback
	if lang := string(lang); lang != "" {
		if l := lexers.Get(lang); l != nil {
			lexer = l
		} else {
			content = lang + "\n" + content
		}
	}

	var fmtter = formatters.Get("tview-256bit")
	if fmtter == nil {
		fmtter = formatters.Fallback
	}

	var style = styles.Get(HighlightStyle)
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return string(literal)
	}

	var code strings.Builder

	if err := fmtter.Format(&code, style, iterator); err != nil {
		return string(literal)
	}

	wrapped := tview.WordWrap(code.String(), 80)

	for _, l := range wrapped {
		if l != "[-]" {
			s.WriteString("\n[grey]┃[-] " + l)
		}
	}

	if !strings.HasSuffix(s.String(), "\n") {
		return s.String() + "\n"
	}

	return s.String()
}
