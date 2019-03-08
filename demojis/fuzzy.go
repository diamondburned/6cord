package demojis

import (
	"github.com/diamondburned/discordgo"
	"github.com/sahilm/fuzzy"
)

// Emojis contains a list of emojis
// DiscordEmojis generate Discordgo emojis
// with ID always being -2
var Emojis, DiscordEmojis = makeArray()

func makeArray() (a []string, d []*discordgo.Emoji) {
	a = make([]string, len(emojiCodeMap))
	d = make([]*discordgo.Emoji, len(emojiCodeMap))

	i := 0
	for e := range emojiCodeMap {
		a[i] = e
		d[i] = &discordgo.Emoji{
			ID:   -2,
			Name: e,
		}

		i++
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
