package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rumblefrog/discordgo"
	"github.com/sahilm/fuzzy"
)

type allChannels []fuzzyReadState

var channelFuzzyStore = make(map[int64]string)

type fuzzyReadState struct {
	*discordgo.Channel
	Unread bool
}

var channelFuzzyCache = allChannels([]fuzzyReadState{})

// String returns the fuzzy search part of the struct
func (ac allChannels) String(i int) string {
	name, ok := channelFuzzyStore[ac[i].ID]
	if !ok {
		g, err := d.State.Guild(ac[i].GuildID)
		if err != nil {
			if len(ac[i].Recipients) > 0 && ac[i].Name == "" {
				recips := make([]string, len(ac[i].Recipients))
				for i, r := range ac[i].Recipients {
					recips[i] = r.Username
				}

				name = HumanizeStrings(recips)
				channelFuzzyStore[ac[i].ID] = name

			} else {
				name = ac[i].Name
				channelFuzzyStore[ac[i].ID] = name
			}

		} else {
			name = ac[i].Name + " (" + g.Name + ")"
			channelFuzzyStore[ac[i].ID] = name
		}
	}

	if ac[i].Unread {
		return "[::b]" + name + "[::-]"
	}

	return "[::d]" + name + "[::-]"
}

// Len returns the length
func (ac allChannels) Len() int {
	return len(ac)
}

func fuzzyChannels(last string) {
	var fuzzied []fuzzy.Match

	if len(last) > 0 {
		if len(channelFuzzyCache) == 0 {
			for _, c := range d.State.PrivateChannels {
				channelFuzzyCache = append(
					channelFuzzyCache,
					fuzzyReadState{c, isUnread(c)},
				)
			}

			for _, g := range d.State.Guilds {
				for _, c := range g.Channels {
					if isSendCh(c.Type) {
						channelFuzzyCache = append(
							channelFuzzyCache,
							fuzzyReadState{c, isUnread(c)},
						)
					}
				}
			}
		}

		fuzzied = fuzzy.FindFrom(
			strings.TrimPrefix(last, "#"),
			channelFuzzyCache,
		)

		var guildID int64

		c, err := d.State.Channel(ChannelID)
		if err == nil {
			guildID = c.GuildID
		}

		sort.SliceStable(fuzzied, func(i, j int) bool {
			return channelFuzzyCache[fuzzied[i].Index].GuildID == guildID
		})
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, fz := range fuzzied {
			autocomp.InsertItem(
				i,
				fz.Str, "",
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

			words := strings.Fields(input.GetText())

			withoutlast := words[:len(words)-1]
			withoutlast = append(withoutlast, fmt.Sprintf(
				"<#%d> ", channelFuzzyCache[fuzzied[i].Index].ID,
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
