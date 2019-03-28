package md

import (
	"fmt"
	"io"
	"strings"

	"github.com/diamondburned/tview"
	"github.com/russross/blackfriday"
	bf "github.com/russross/blackfriday"
)

type tviewMarkdown struct{}

// RenderNode is a big tree for rendering markdown
func (*tviewMarkdown) RenderNode(w io.Writer, node *bf.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case bf.Softbreak:
		w.Write([]byte("\n"))

	case bf.Hardbreak:
		w.Write([]byte("\n\n"))

	case bf.Emph:
		if entering {
			w.Write([]byte("[::i]"))
		} else {
			w.Write([]byte("[::-]"))
		}

	case bf.Strong:
		if entering {
			w.Write([]byte("[::b]"))
		} else {
			w.Write([]byte("[::-]"))
		}

	case bf.Del:
		if entering {
			w.Write([]byte("[::s]"))
		} else {
			w.Write([]byte("[::-]"))
		}

	case bf.Code:
		fmt.Fprintf(w, "[:#4f4f4f:]%s[:-:]", string(node.Literal))

	case bf.Link:
		fmt.Fprintf(w, "[::u]%s[::-]", string(node.Literal))

	case bf.CodeBlock:
		w.Write([]byte(RenderCodeBlock(node)))

	case bf.BlockQuote:
		if entering {
			node.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
				switch node.Type {
				case bf.BlockQuote:
				default:
					if entering {
						literal := string(node.Literal)

						if literal == "" {
							return blackfriday.GoToNext
						}

						for _, l := range strings.Split(literal, "\n") {
							fmt.Fprint(w, "[green]>"+l+"[-]\n")
						}
					}
				}

				return blackfriday.GoToNext
			})

			return blackfriday.SkipChildren
		}

	case bf.Paragraph:
		if checkIsBlockquote(node) {
			break
		}

	case bf.Document:
		break

	default:
		if checkIsBlockquote(node) {
			break
		}

		fmt.Fprint(w, tview.Escape(string(node.Literal)))
	}

	return blackfriday.GoToNext
}

// RenderHeader ..
func (*tviewMarkdown) RenderHeader(w io.Writer, ast *blackfriday.Node) {}

// RenderFooter ..
func (*tviewMarkdown) RenderFooter(w io.Writer, ast *blackfriday.Node) {}

func checkIsBlockquote(node *bf.Node) bool {
	if node.Parent != nil {
		switch node.Parent.Type {
		case bf.BlockQuote:
			return true
		case bf.Paragraph:
			return checkIsBlockquote(node.Parent)
		}
	}

	return false
}
