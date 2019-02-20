package main

import (
	"github.com/sahilm/fuzzy"
)

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
