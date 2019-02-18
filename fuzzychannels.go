package main

import (
	"fmt"
	"strings"

	"github.com/rumblefrog/discordgo"
	"github.com/sahilm/fuzzy"
)

type allChannels []*discordgo.Channel

// String returns the fuzzy search part of the struct
func (ac allChannels) String(i int) string {
	g, err := d.State.Guild(ac[i].GuildID)
	if err != nil {
		if len(ac[i].Recipients) > 0 {
			recips := make([]string, len(ac[i].Recipients))
			for i, r := range ac[i].Recipients {
				recips[i] = r.Username
			}

			return HumanizeStrings(recips)
		}

		return ac[i].Name
	}

	return ac[i].Name + " (" + g.Name + ")"
}

// Len returns the length
func (ac allChannels) Len() int {
	return len(ac)
}

func fuzzyChannels(last string) {
	var (
		channels = make(allChannels, len(d.State.PrivateChannels))
		fuzzied  []fuzzy.Match
	)

	if len(last) > 0 {
		copy(channels, d.State.PrivateChannels)
		for _, g := range d.State.Guilds {
			channels = append(channels, g.Channels...)
		}

		fuzzied = fuzzy.FindFrom(
			strings.TrimPrefix(last, "#"), channels,
		)
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, fz := range fuzzied {
			if i == 10 {
				break
			}

			autocomp.InsertItem(
				i,
				fz.Str, "",
				rune(0x31+i),
				nil,
			)
		}

		rightflex.ResizeItem(autocomp, min(len(fuzzied), 10), 1)

		autofillfunc = func(i int) {
			words := strings.Fields(input.GetText())

			withoutlast := words[:len(words)-1]
			withoutlast = append(withoutlast, fmt.Sprintf(
				"<#%d> ", channels[fuzzied[i].Index].ID,
			))

			input.SetText(strings.Join(withoutlast, " "))

			clearList()

			app.SetFocus(input)
		}

	} else {
		rightflex.ResizeItem(autocomp, 1, 1)
	}

	app.Draw()
}
