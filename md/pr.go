package md

import (
	"fmt"

	"github.com/diamondburned/mark"
	"github.com/diamondburned/tview"
)

// RenderParagraph recursively renders a paragraph
func RenderParagraph(n mark.Node) (s string) {
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

		case *mark.ParagraphNode:
			s += RenderParagraph(n)

		case *mark.BlockQuoteNode:
			s += RenderBlockQuote(n)

		case *mark.CodeNode:
			s += n.Lang + n.Text

		default:
			s += tview.Escape(n.Render())
		}
	}

	s += "\n"

	return
}
