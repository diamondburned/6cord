package md

import (
	"fmt"
	"io"
	"math"

	"github.com/alecthomas/chroma"
	"github.com/rivo/tview"
)

var (
	c = chroma.MustParseColour
)

func entryToEscapeSequence(table *ttyTable, entry chroma.StyleEntry) string {
	var out string
	if entry.Colour.IsSet() {
		out += table.foreground[findClosest(table, entry.Colour)]
	}
	return out
}

func findClosest(table *ttyTable, seeking chroma.Colour) chroma.Colour {
	closestColour := chroma.Colour(0)
	closest := float64(math.MaxFloat64)
	for colour := range table.foreground {
		distance := colour.Distance(seeking)
		if distance < closest {
			closest = distance
			closestColour = colour
		}
	}
	return closestColour
}

func styleToEscapeSequence(table *ttyTable, style *chroma.Style) map[chroma.TokenType]string {
	out := map[chroma.TokenType]string{}
	for _, ttype := range style.Types() {
		entry := style.Get(ttype)
		out[ttype] = entryToEscapeSequence(table, entry)
	}
	return out
}

type indexedTTYFormatter struct {
	table *ttyTable
}

func (c *indexedTTYFormatter) Format(w io.Writer, style *chroma.Style, it chroma.Iterator) (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			err = perr.(error)
		}
	}()

	lastColor := false

	theme := styleToEscapeSequence(c.table, style)
	for token := it(); token != chroma.EOF; token = it() {
		color, ok := theme[token.Type]
		if !ok {
			color, ok = theme[token.Type.SubCategory()]
			if !ok {
				color = theme[token.Type.Category()]
			}
		}

		if color != "" {
			fmt.Fprint(w, color)
			lastColor = true
		}

		fmt.Fprint(w, tview.Escape(token.Value))

		if lastColor {
			fmt.Fprint(w, "[-]")
			lastColor = false
		}
	}

	return nil
}
