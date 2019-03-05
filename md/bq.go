package md

import (
	"strings"

	"github.com/diamondburned/mark"
)

// RenderBlockQuote recursively renders a block quote
func RenderBlockQuote(n mark.Node) (s string) {
	q, _ := n.(*mark.BlockQuoteNode)

	for _, c := range q.Nodes {
		switch c := c.(type) {
		case *mark.ParagraphNode:
			for _, t := range c.Nodes {
				if t, ok := t.(*mark.TextNode); ok {
					for _, l := range strings.Split(t.Text, "\n") {
						s += "[green]>" + l + "[-]\n"
					}
				} else {
					s += c.Render()
				}
			}

		case *mark.BlockQuoteNode:
			// recursion recursion recursion recursion recursion
			s += "[green]>" + RenderBlockQuote(c) + "[-]"

		default:
			s += c.Render()
		}
	}

	return strings.TrimSuffix(s, "\n")
}
