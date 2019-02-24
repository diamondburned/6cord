package md

import (
	"fmt"
	"io"
	"math"

	"github.com/alecthomas/chroma"
)

var c = chroma.MustParseColour

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

	theme := styleToEscapeSequence(c.table, style)
	for token := it(); token != chroma.EOF; token = it() {
		clr, ok := theme[token.Type]
		if !ok {
			clr, ok = theme[token.Type.SubCategory()]
			if !ok {
				clr = theme[token.Type.Category()]
			}
		}
		if clr != "" {
			fmt.Fprint(w, clr)
		}
		fmt.Fprint(w, token.Value)
	}

	return nil
}
