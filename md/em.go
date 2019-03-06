package md

import "github.com/diamondburned/mark"

func tagReflect(t string) string {
	switch t {
	case "strong":
		return "[::b]"
	case "em":
		return "[::b]"
	case "del":
		return "[::d]"
	case "code":
		return "[:#4f4f4f:]"
	}

	return ""
}

// RenderEmphasis recursively renders emphasis
func RenderEmphasis(n mark.Node) (s string) {
	em, _ := n.(*mark.EmphasisNode)

	s += tagReflect(em.Tag())

	for _, n := range em.Nodes {
		if _, ok := n.(*mark.EmphasisNode); ok {
			s += RenderEmphasis(n)
		} else {
			s += n.Render()
		}
	}

	s += "[:-:-]"

	return
}
