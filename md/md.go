package md

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/diamondburned/mark"
	"github.com/rivo/tview"
)

// HighlightStyle determines the syntax highlighting colorstyle:
// https://xyproto.github.io/splash/docs/all.html
var HighlightStyle = "vs"

var trashyCodeBlockMatching = regexp.MustCompile("(.)```")

// Parse parses md into tview string
func Parse(s string) string {
	return ParseNoEscape(tview.Escape(s))
}

// ParseNoEscape parses md into tview string without escaping it
func ParseNoEscape(s string) string {
	s = trashyCodeBlockMatching.ReplaceAllString(s, "$1\n```")
	s = fixQuotes(s)

	m := mark.New(s, &mark.Options{})
	if m == nil {
		return s
	}

	m.AddRenderFn(mark.NodeText, func(n mark.Node) (s string) {
		t, _ := n.(*mark.TextNode)
		return t.Text
	})

	m.AddRenderFn(mark.NodeEmphasis, func(n mark.Node) (s string) {
		e, _ := n.(*mark.EmphasisNode)
		for _, n := range e.Nodes {
			s += n.Render()
		}

		return
	})

	m.AddRenderFn(mark.NodeBlockQuote, RenderBlockQuote)
	m.AddRenderFn(mark.NodeCode, RenderCodeBlock)

	m.AddRenderFn(mark.NodeParagraph, func(n mark.Node) (s string) {
		p, _ := n.(*mark.ParagraphNode)
		for _, n := range p.Nodes {
			switch n := n.(type) {
			case *mark.EmphasisNode:
				s += tagReflect(n.Tag())

				for _, n := range n.Nodes {
					s += n.Render()
				}

				s += "[:-:-]"

			case *mark.LinkNode:
				if n.Title == "" {
					s += "[::u]" + n.Href + "[::-]"
				} else {
					s += fmt.Sprintf(
						"[%s](%s)",
						n.Title, n.Href,
					)
				}

			default:
				s += n.Render()
			}
		}

		s += "\n"

		return
	})

	return strings.TrimSpace(m.Render())
}
