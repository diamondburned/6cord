// Package antitele attempts to insert invisible runes into messages, so that
// collected telemetric data aren't usable for selling to companies, protecting
// privacy
package antitele

import (
	"math/rand"
	"strings"
	"time"
	"unicode"
)

// Probability modifies the probability of an invisible character appearing.
// The lower this is, the less visible characters you can send.
// Rules: [0, n) 0 > n >= +Inf
var Probability = 5

// ZeroWidthRunes is the array containing all invisible runes.
// U+200B is used for obfuscating.
var ZeroWidthRunes = []rune{
	'\u200b', '\u200c', '\u200d', '\ufeff',
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Insert works its magic
func Insert(s string) string {
	var words = strings.Fields(s)

	var amount = max(0, 2048-len(s))
	var needle int

	for i, w := range words {
		needle++
		if amount < needle {
			break
		}

		// Skip over links for them to be clickable
		if strings.HasPrefix(w, "http") {
			continue
		}

		// Skip if it's not a word
		if strings.IndexFunc(w, func(c rune) bool {
			return !unicode.IsLetter(c)
		}) != -1 {
			continue
		}

		// Words too short probably doens't need
		// to be obfuscated
		if len(w) < 3 {
			continue
		}

		words[i] = obf(w)
	}

	return strings.Join(words, " ")
}

func obf(s string) string {
	var r = []rune(s)
	var i = rand.Intn(len(r))

	return string(r[:i]) + "\u200b" + string(r[i:])
}

func max(i, j int) int {
	if i > j {
		return i
	}

	return j
}

func containsRunes(s string, trs ...rune) bool {
	var i int
	for _, r := range []rune(s) {
		for _, tr := range trs {
			if tr == r {
				i++
			}
		}
	}

	return i == len(trs)
}
