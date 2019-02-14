package discordgo

import "github.com/sahilm/fuzzy"

// Emojis contains a list of emojis
var Emojis = makeArray()

func makeArray() (a []string) {
	for e := range emojiCodeMap {
		a = append(a, e)
	}

	return
}

// FuzzyEmojis fuzzy searches the emojis
// Argument: p == pattern
func FuzzyEmojis(p string) []fuzzy.Match {
	return fuzzy.Find(p, Emojis)
}

// MatchEmoji matches a fuzzy search with the emoji
func MatchEmoji(m fuzzy.Match) string {
	vl, ok := emojiCodeMap[m.Str]
	if !ok {
		// should never happen
		return ""
	}

	return vl
}

// GetEmojiFromKey returns "", false if emoji isn't found
func GetEmojiFromKey(k string) (string, bool) {
	v, ok := emojiCodeMap[k]
	return v, ok
}
