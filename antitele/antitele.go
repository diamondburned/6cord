// Package antitele attempts to insert invisible runes into messages, so that
// collected telemetric data aren't usable for selling to companies, protecting
// privacy
package antitele

import (
	"math/rand"
	"strings"
	"time"
)

// Probability modifies the probability of an invisible character appearing.
// The lower this is, the less visible characters you can send.
// Rules: [0, n) 0 > n >= +Inf
var Probability = 5

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Insert works its magic
func Insert(s string) string {
	var r = []rune(s)
	var b strings.Builder

	for _, r := range r {
		b.WriteRune(r)

		if rand.Intn(Probability) == 0 {
			b.WriteRune('\u200d')
		}
	}

	return b.String()
}
