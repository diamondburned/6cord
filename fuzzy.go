package main

import (
	"github.com/sahilm/fuzzy"
)

func stateResetter() {
	channelFuzzyCache = allChannels([]fuzzyReadState{})
	allMessages = []string{}
	autocomp.SetChangedFunc(nil)
	messagesView.Highlight()
}

func formatNeedle(m fuzzy.Match) (f string) {
	isHL := false

	for i := 0; i < len(m.Str); i++ {
		if fuzzyHasNeedle(i, m.MatchedIndexes) {
			f += "[::u]" + string(m.Str[i])
			isHL = true
		} else {
			if isHL {
				f += "[::-]"
			}

			f += string(m.Str[i])
		}
	}

	return
}

func fuzzyHasNeedle(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}

func min(i, j int) int {
	if i < j {
		return i
	}

	return j
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
var onhoverfn func(i int)
