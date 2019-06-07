package md

import (
	"log"
	"regexp"
	"strings"

	md "github.com/diamondburned/markdown"
	"github.com/diamondburned/markdown/ast"
	ps "github.com/diamondburned/markdown/parser"
	"github.com/diamondburned/tview"
)

const extensions = 0 |
	ps.Spoilers |
	ps.FencedCode |
	ps.NoIntraEmphasis |
	ps.Strikethrough |
	ps.NoIndentCodeBlock |
	ps.HardLineBreak

// HighlightStyle determines the syntax highlighting colorstyle:
// https://xyproto.github.io/splash/docs/all.html
var HighlightStyle = "vs"

var trashyCodeBlockMatching = regexp.MustCompile("(.)```")

// Parse parses md into tview strings
func Parse(s string) string {
	return parse(s, true)
}

// ParseNoInsert parses md without inserting no-width spaces
func ParseNoInsert(s string) string {
	return parse(s, false)
}

func parse(s string, obf bool) (results string) {
	results = s
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	s = trashyCodeBlockMatching.ReplaceAllString(s, "$1\n```\n")
	s = fixQuotes(s)

	var builder strings.Builder

	node := md.Parse([]byte(s), ps.NewWithExtensions(extensions))
	ast.WalkFunc(node, func(node ast.Node, entering bool) ast.WalkStatus {
		switch node := node.(type) {
		case *ast.Text:
			builder.WriteString(tview.Escape(string(node.Literal)))
		case *ast.Softbreak:
			builder.WriteByte('\n')
		case *ast.Hardbreak:
			builder.WriteByte('\n')
		case *ast.Und:
			if entering {
				builder.Write([]byte("[::u]"))
			} else {
				builder.Write([]byte("[::-]"))
			}
		case *ast.Emph:
			if entering {
				builder.Write([]byte("[::i]"))
			} else {
				builder.Write([]byte("[::-]"))
			}
		case *ast.Strong:
			if entering {
				builder.Write([]byte("[::b]"))
			} else {
				builder.Write([]byte("[::-]"))
			}
		case *ast.Del:
			if entering {
				builder.Write([]byte("[::s]"))
			} else {
				builder.Write([]byte("[::-]"))
			}
		case *ast.Spoiler:
			if entering {
				builder.Write([]byte("[#777777]"))
			} else {
				builder.Write([]byte("[-]"))
			}
		case *ast.Code:
			if entering {
				builder.Write([]byte("[:#4f4f4f:]"))
				builder.WriteString(tview.Escape(string(node.Literal)))
				builder.Write([]byte("[:-:]"))
			}
		case *ast.CodeBlock:
			builder.WriteString(RenderCodeBlock(
				node.Info, node.Literal,
			))

			builder.WriteByte('\n')
		case *ast.BlockQuote:
			if entering {
				if _, ok := node.Parent.(*ast.BlockQuote); !ok {
					builder.Write([]byte("[#789922]"))
				}

				builder.Write([]byte(">"))
			} else {
				if _, ok := node.Parent.(*ast.BlockQuote); !ok {
					builder.Write([]byte("[-]\n"))
				}
			}
		case *ast.Paragraph:
			if !entering {
				if _, ok := node.Parent.(*ast.BlockQuote); !ok {
					builder.Write([]byte("\n"))
				}
			}
		default:
			if l := node.AsLeaf(); l != nil {
				builder.Write(l.Literal)
			} else if c := node.AsContainer(); c != nil {
				builder.Write(c.Literal)
			}
		}

		return ast.GoToNext
	})

	return strings.Trim(builder.String(), "\n")
}
