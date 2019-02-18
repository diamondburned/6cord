package main

import (
	"github.com/sahilm/fuzzy"
)

func formatNeedle(m fuzzy.Match, a string) (f string) {
	for i := 0; i < len(a); i++ {
		if fuzzyHasNeedle(i, m.MatchedIndexes) {
			f += "[::u]" + string(a[i]) + "[::-]"
		} else {
			f += string(a[i])
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
