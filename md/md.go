package md

import (
	"strings"

	md "github.com/gomarkdown/markdown"
	ast "github.com/gomarkdown/markdown/ast"
	"github.com/rivo/tview"
)

// Parse parses md into tview string
func Parse(s string) string {
	return ParseNoEscape(tview.Escape(s))
}

// ParseNoEscape parses md into tview string without escaping it
func ParseNoEscape(s string) string {
	b := strings.Builder{}

	doc := md.Parse([]byte(s), nil)
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		switch node := node.(type) {
		case *ast.Text:
			b.Write(node.Literal)
		case *ast.Softbreak, *ast.Hardbreak:
			b.WriteRune('\n')
		case *ast.Emph:
			b.WriteString(isFormatEnter(entering, "b"))
			b.Write(node.Content)
		case *ast.Strong:
			b.WriteString(isFormatEnter(entering, "b"))
			b.Write(node.Content)
		case *ast.Del:
			b.WriteString(isFormatEnter(entering, "d"))
			b.Write(node.Content)
		case *ast.BlockQuote:
			if entering {
				b.WriteString("[green]>")
			} else {
				b.WriteString("[-]")
			}

			b.Write(node.Content)
		case *ast.Link:
			b.WriteString(isFormatEnter(entering, "u"))
			b.Write(node.Title)
		case *ast.Code:
			if entering {
				b.WriteString("[:black:]")
			} else {
				b.WriteString("[:-:]")
			}

			b.Write(node.Literal)
		case *ast.CodeBlock:
			if entering {
				b.WriteString("[:black:]")
			} else {
				b.WriteString("[:-:]")
			}

			b.Write(node.Literal)
		case *ast.Aside:
			b.Write(node.Literal)
		case *ast.CrossReference:
			b.Write(node.Literal)
		case *ast.Citation:
			b.Write(node.Literal)
		case *ast.Image:
			b.Write(node.Literal)
		case *ast.Caption:
			b.Write(node.Literal)
		case *ast.CaptionFigure:
			b.Write(node.Literal)
		case *ast.Document:
			b.Write(node.Literal)
		case *ast.Paragraph:
			b.Write(node.Literal)
		case *ast.HTMLSpan:
			b.Write(node.Literal)
		case *ast.HTMLBlock:
			b.Write(node.Literal)
		case *ast.Heading:
			b.Write(node.Literal)
		case *ast.HorizontalRule:
			b.Write(node.Literal)
		case *ast.List:
			b.Write(node.Literal)
		case *ast.ListItem:
			b.Write(node.Literal)
		case *ast.Table:
			b.Write(node.Literal)
		case *ast.TableCell:
			b.Write(node.Literal)
		case *ast.TableHeader:
			b.Write(node.Literal)
		case *ast.TableBody:
			b.Write(node.Literal)
		case *ast.TableRow:
			b.Write(node.Literal)
		case *ast.TableFooter:
			b.Write(node.Literal)
		case *ast.Math:
			b.Write(node.Literal)
		case *ast.MathBlock:
			b.Write(node.Literal)
		case *ast.DocumentMatter:
			b.Write(node.Literal)
		case *ast.Callout:
			b.Write(node.Literal)
		case *ast.Index:
			b.Write(node.Literal)
		case *ast.Subscript:
			b.Write(node.Literal)
		case *ast.Superscript:
			b.Write(node.Literal)
		case *ast.Footnotes:
			b.Write(node.Literal)
		default:

		}

		return ast.GoToNext
	})

	return b.String()
}
