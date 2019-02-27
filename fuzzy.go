package main

import (
	"strings"

	"github.com/sahilm/fuzzy"
)

func fuzzyBuffer(buffer chan string) {
	for text := range buffer {
		words := strings.Fields(text)

		if len(words) < 1 {
			stateResetter()
			clearList()

			continue
		}

		switch last := words[len(words)-1]; {
		case strings.HasPrefix(last, "@"):
			fuzzyMentions(last)
		case strings.HasPrefix(last, "#"):
			fuzzyChannels(last)
		case strings.HasPrefix(last, ":"):
			fuzzyEmojis(last)
		case strings.HasPrefix(text, "/upload "):
			fuzzyUpload(text)
		case strings.HasPrefix(text, "/"):
			if len(words) == 1 {
				fuzzyCommands(text)
			}
		default:
			stateResetter()
			clearList()
		}
	}
}

func stateResetter() {
	// Function for future calls
	isUnreadFuzzyReset()
}

func clearList() {
	rightflex.ResizeItem(autocomp, 1, 1)
	autocomp.Clear()
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
