package md

import (
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/russross/blackfriday"
)

// RenderCodeBlock renders the node to a syntax
// highlighted code
func RenderCodeBlock(node *blackfriday.Node) (s string) {
	var (
		lang = string(node.CodeBlockData.Info)
	)

	content := strings.TrimFunc(
		string(node.Literal),
		func(r rune) bool {
			return r == '\n'
		},
	)

	if content == "" {
		content = lang
	}

	var lexer = lexers.Fallback
	if lang != "" {
		if l := lexers.Get(lang); l != nil {
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

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return content
	}

	code := strings.Builder{}

	if err := fmtter.Format(&code, style, iterator); err != nil {
		return content
	}

	for _, l := range strings.Split(code.String(), "\n") {
		if l != "[-]" {
			s += "\n[grey]┃[-] " + l
		}
	}

	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}

	return
}
