package main

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"github.com/sahilm/fuzzy"
)

// String returns the fuzzy search part of the struct
func (gm UserStoreArray) String(i int) string {
	var s = gm[i].Name

	if gm[i].Nick != "" {
		s += " "
		s += gm[i].Nick
	}

	return s
}

// Len returns the length
func (gm UserStoreArray) Len() int {
	return len(gm)
}

// FuzzyMembers fuzzy searches and returns the slice of results
func FuzzyMembers(pattern string, s *UserStore) (fzr UserStoreArray) {
	results := fuzzy.FindFrom(pattern, s.Data)
	for i := 0; i < len(results) && i < 10; i++ {
		fzr = append(fzr, s.Data[results[i].Index])
	}

	return
}

func fuzzyMentions(last string) {
	var fuzzied UserStoreArray

	if len(last) > 0 {
		fuzzied = FuzzyMembers(
			strings.TrimPrefix(last, "@"), us,
		)
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, u := range fuzzied {
			var username = u.Name + "[::d]#" + u.Discrim + "[::-]"
			if u.Nick != "" {
				username += " (" + tview.Escape(u.Nick) + ")"
			}

			autocomp.InsertItem(
				i,
				username, "",
				rune(0x31+i),
				nil,
			)
		}

		rightflex.ResizeItem(autocomp, min(len(fuzzied), 10), 1)

		autofillfunc = func(i int) {
			words := strings.Fields(input.GetText())

			withoutlast := words[:len(words)-1]
			withoutlast = append(withoutlast, fmt.Sprintf(
				"<@%d> ", fuzzied[i].ID,
			))

			input.SetText(strings.Join(withoutlast, " "))

			clearList()

			app.SetFocus(input)
		}

	} else {
		rightflex.ResizeItem(autocomp, 1, 1)
	}

	app.Draw()
}
