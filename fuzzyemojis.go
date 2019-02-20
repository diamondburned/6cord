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

func fuzzyEmojis(last string) {
	var (
		fuzzied []fuzzy.Match
		emojis  DiscordEmojis
	)

	if len(last) > 0 {
		c, err := d.State.Channel(ChannelID)
		if err != nil {
			return
		}

		g, err := d.State.Guild(c.GuildID)
		if err != nil {
			return
		}

		emojis = g.Emojis
		emojis = append(emojis, demojis.DiscordEmojis...)

		fuzzied = fuzzy.FindFrom(
			strings.TrimPrefix(last, ":"),
			emojis,
		)
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, m := range fuzzied {
			if i == 26 {
				break
			}

			autocomp.InsertItem(
				i,
				":"+m.Str+":", "",
				rune(0x31+i),
				nil,
			)
		}

		rightflex.ResizeItem(autocomp, min(len(fuzzied), 10), 1)

		autofillfunc = func(i int) {
			var (
				words  = strings.Fields(input.GetText())
				emoji  = emojis[fuzzied[i].Index]
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

// fuck you you fucking tview dev
// you could just fucking globalize the goddamn selection
// function, but no. you didn't. why the fuck didn't you?
// are you fucking retarded in the head? stop trying to
// handle everything by you yourself and your shitty shoddy
// little functions, you fucking stupid asshat
// bloody fucking jesus i fucking hate doing this, but
// this is literally my only fucking choice
var autofillfunc func(i int)
