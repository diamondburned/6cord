package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rumblefrog/discordgo"
	"github.com/sahilm/fuzzy"
)

type allChannels []*discordgo.Channel

var channelFuzzyStore = make(map[int64]string)
var channelFuzzyReadstateStore = make(map[int64]bool)

func isUnreadFuzzy(ch *discordgo.Channel) bool {
	bl, ok := channelFuzzyReadstateStore[ch.ID]
	if !ok {
		bl = isUnread(ch)
		channelFuzzyReadstateStore[ch.ID] = bl
	}

	return bl
}

func isUnreadFuzzyReset() {
	channelFuzzyReadstateStore = make(map[int64]bool)
}

// String returns the fuzzy search part of the struct
func (ac allChannels) String(i int) (name string) {
	var color = "d"
	if isUnreadFuzzy(ac[i]) {
		color = "b"
	}

	name, ok := channelFuzzyStore[ac[i].ID]
	if !ok {
		g, err := d.State.Guild(ac[i].GuildID)
		if err != nil {
			if len(ac[i].Recipients) > 0 && ac[i].Name == "" {
				recips := make([]string, len(ac[i].Recipients))
				for i, r := range ac[i].Recipients {
					recips[i] = r.Username
				}

				r := HumanizeStrings(recips)
				channelFuzzyStore[ac[i].ID] = r

			} else {
				name = ac[i].Name
				channelFuzzyStore[ac[i].ID] = name
			}

		} else {
			name = ac[i].Name + " (" + g.Name + ")"
			channelFuzzyStore[ac[i].ID] = name
		}
	}

	return "[::" + color + "]" + name + "[::-]"
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

		var guildID int64

		c, err := d.State.Channel(ChannelID)
		if err == nil {
			guildID = c.GuildID
		}

		sort.SliceStable(fuzzied, func(i, j int) bool {
			return channels[fuzzied[i].Index].GuildID == guildID
		})
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

			switch {
			case strings.HasPrefix(input.GetText(), "/goto "):
				input.SetText("")
				gotoChannel(withoutlast)
				return

			default:
				input.SetText(strings.Join(withoutlast, " "))
			}

			clearList()

			app.SetFocus(input)
		}

	} else {
		rightflex.ResizeItem(autocomp, 1, 1)
	}

	app.Draw()
}
