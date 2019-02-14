package main

import (
	"strings"

	"github.com/sahilm/fuzzy"
	"gitlab.com/diamondburned/6cord/demojis"
)

func fuzzyEmojis(last string) {
	var fuzzied []fuzzy.Match

	if len(last) > 0 {
		fuzzied = demojis.FuzzyEmojis(
			strings.TrimPrefix(last, "@"),
		)
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, m := range fuzzied {
			autocomp.InsertItem(
				i,
				m.Str, "",
				rune(0x31+i),
				nil,
			)
		}

		rightflex.ResizeItem(autocomp, 10, 1)

		autofillfunc = func(i int) {
			var (
				words = strings.Fields(input.GetText())
				emoji = demojis.MatchEmoji(fuzzied[i])
			)

			withoutlast := words[:len(words)-1]
			withoutlast = append(
				withoutlast,
				emoji+" ",
			)

			input.SetText(strings.Join(withoutlast, " "))

			clearList()

			app.SetFocus(input)
		}

	} else {
		rightflex.ResizeItem(autocomp, 1, 1)
	}

	app.Draw()
}

// fuck you you fucking tview dev
// you could just fucking globalize the goddamn selection
// function, but no. you didn't. why the fuck didn't you?
// are you fucking retarded in the head? stop trying to
// handle everything by you yourself and your shitty shoddy
// little functions, you fucking stupid asshat
// bloody fucking jesus i fucking hate doing this, but
// this is literally my only fucking choice
var autofillfunc func(i int)
