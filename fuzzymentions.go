package main

import (
	"fmt"
	"strings"

	"github.com/sahilm/fuzzy"
)

var (
	// tell me a better way
	currentFuzzy UserStoreArray
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

// Len returns the length of the Emojis slice
func (gm UserStoreArray) Len() int {
	return len(gm)
}

// FuzzyMembers fuzzy searches the list of emojis and returns the slice of results
func FuzzyMembers(pattern string, s *UserStore) (fzr UserStoreArray) {
	results := fuzzy.FindFrom(pattern, s.Data)
	for i := 0; i < len(results) && i < 8; i++ {
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
		currentFuzzy = fuzzied

		for i, u := range fuzzied {
			var username = u.Name + "[::d]#" + u.Discrim + "[::-]"
			if u.Nick != "" {
				username += " (" + u.Nick + ")"
			}

			autocomp.InsertItem(
				i,
				username, "",
				rune(0x31+i),
				nil,
			)
		}

		rightflex.ResizeItem(autocomp, 10, 1)
	}

	app.Draw()
}

func applyMention(i int) {
	words := strings.Fields(input.GetText())

	withoutlast := words[:len(words)-1]
	withoutlast = append(withoutlast, fmt.Sprintf(
		"<@%d> ", currentFuzzy[i].ID,
	))

	input.SetText(strings.Join(withoutlast, " "))

	clearList()
}
