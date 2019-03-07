package main

import (
	"fmt"
	"strings"

	"github.com/rumblefrog/discordgo"
	"github.com/sahilm/fuzzy"
	"gitlab.com/diamondburned/6cord/demojis"
)

// DiscordEmojis ..
type DiscordEmojis []*discordgo.Emoji

func (de DiscordEmojis) String(i int) string {
	return de[i].Name
}

func (de DiscordEmojis) Len() int {
	return len(de)
}

var discordEmojis = DiscordEmojis([]*discordgo.Emoji{})

func fuzzyEmojis(last string) {
	var fuzzied []fuzzy.Match

	if len(last) > 0 && Channel != nil {
		if len(discordEmojis) < 1 {
			c, err := d.State.Channel(Channel.ID)
			if err != nil {
				return
			}

			emojis := DiscordEmojis{}

			g, err := d.State.Guild(c.GuildID)
			if err == nil {
				emojis = g.Emojis
				emojis = append(emojis, demojis.DiscordEmojis...)
			} else {
				emojis = demojis.DiscordEmojis
			}

			discordEmojis = emojis
		}

		fuzzied = fuzzy.FindFrom(
			strings.TrimPrefix(last, ":"),
			discordEmojis,
		)
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, m := range fuzzied {
			autocomp.InsertItem(
				i,
				":"+m.Str+":", "",
				rune(0x31+i),
				nil,
			)

			if i == 25 {
				break
			}
		}

		rightflex.ResizeItem(autocomp, min(len(fuzzied), 10), 1)

		autofillfunc = func(i int) {
			defer stateResetter()

			var (
				words  = strings.Fields(input.GetText())
				emoji  = discordEmojis[fuzzied[i].Index]
				insert string
			)

			if emoji.ID == -2 {
				e, ok := demojis.GetEmojiFromKey(emoji.Name)
				if ok {
					insert = e
				}
			} else {
				var a string
				if emoji.Animated {
					a = "a"
				}

				insert = fmt.Sprintf(
					"<%s:%s:%d>",
					a, emoji.Name, emoji.ID,
				)
			}

			withoutlast := words[:len(words)-1]
			withoutlast = append(
				withoutlast,
				insert+" ",
			)

			input.SetText(strings.Join(withoutlast, " "))

			clearList()

			app.SetFocus(input)
		}

	} else {
		rightflex.ResizeItem(autocomp, 1, 1)
	}

	app.Draw()
}
