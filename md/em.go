package md

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/diamondburned/mark"
	"github.com/diamondburned/tview"
)

func tagReflect(t string) string {
	switch t {
	case "strong":
		return "[::b]"
	case "em":
		return "[::i]"
	case "del":
		return "[::s]"
	case "code":
		return "[:#4f4f4f:]"
	}

	return ""
}

// RenderEmphasis recursively renders emphasis
func RenderEmphasis(n mark.Node) (s string) {
	em, _ := n.(*mark.EmphasisNode)

	log.Println(spew.Sdump(em))

	s += tagReflect(em.Tag())

	for _, n := range em.Nodes {
		switch n := n.(type) {
		case *mark.EmphasisNode:
			s += RenderEmphasis(n)
		case *mark.TextNode:
			s += n.Text
		default:
			s += tview.Escape(n.Render())
		}
	}

	s += "[:-:-]"

	return
}
