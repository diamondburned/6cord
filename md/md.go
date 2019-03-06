package md

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/diamondburned/mark"
	"github.com/rivo/tview"
)

// HighlightStyle determines the syntax highlighting colorstyle:
// https://xyproto.github.io/splash/docs/all.html
var HighlightStyle = "vs"

var trashyCodeBlockMatching = regexp.MustCompile("(.)```")

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

	m := mark.New(s, &mark.Options{})
	if m == nil {
		return s
	}

	m.AddRenderFn(mark.NodeText, func(n mark.Node) (s string) {
		t, _ := n.(*mark.TextNode)
		return tview.Escape(t.Text)
	})

	m.AddRenderFn(mark.NodeEmphasis, RenderEmphasis)
	m.AddRenderFn(mark.NodeBlockQuote, RenderBlockQuote)
	m.AddRenderFn(mark.NodeCode, RenderCodeBlock)

	m.AddRenderFn(mark.NodeParagraph, func(n mark.Node) (s string) {
		p, _ := n.(*mark.ParagraphNode)
		for _, n := range p.Nodes {
			switch n := n.(type) {
			case *mark.EmphasisNode:
				s += RenderEmphasis(n)

			case *mark.LinkNode:
				if n.Title == "" {
					s += "[::u]" + tview.Escape(n.Href) + "[::-]"
				} else {
					s += fmt.Sprintf(
						"[%s](%s)",
						tview.Escape(n.Title),
						tview.Escape(n.Href),
					)
				}

			default:
				s += tview.Escape(n.Render())
			}
		}

		s += "\n"

		return
	})

	return strings.TrimSpace(m.Render())
}
