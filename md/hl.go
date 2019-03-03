package md

import (
	"strings"

	"github.com/diamondburned/mark"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// RenderCodeBlock renders the node to a syntax
// highlighted code
func RenderCodeBlock(n mark.Node) (s string) {
	c, _ := n.(*mark.CodeNode)

	var lexer = lexers.Fallback
	if c.Lang != "" {
		if l := lexers.Get(c.Lang); l != nil {
			lexer = l
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

	content := strings.TrimFunc(
		c.Text,
		func(r rune) bool {
			return r == '\n'
		},
	)

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return c.Text
	}

	code := strings.Builder{}

	if err := fmtter.Format(&code, style, iterator); err != nil {
		return c.Text
	}

	for _, l := range strings.Split(code.String(), "\n") {
		s += "\n[grey]┃[-] " + l
	}

	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}

	return
}
