package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/discordgo"
	"github.com/sahilm/fuzzy"
)

type allChannels []fuzzyReadState

type fuzzyReadState struct {
	*discordgo.Channel
	Unread bool
	Format string
}

var channelFuzzyCache = allChannels([]fuzzyReadState{})

// String returns the fuzzy search part of the struct
func (ac allChannels) String(i int) string {
	if ac[i].Unread {
		return "[::b]#" + ac[i].Format + "[::-]"
	}

	return "[::d]#" + ac[i].Format + "[::-]"
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
				var name = c.Name

				if c.Name == "" {
					recips := make([]string, len(c.Recipients))
					for i, r := range c.Recipients {
						recips[i] = r.Username
					}

					name = HumanizeStrings(recips)
				}

				channelFuzzyCache = append(
					channelFuzzyCache,
					fuzzyReadState{c, isUnread(c), name},
				)
			}

			for _, g := range d.State.Guilds {
				for _, c := range g.Channels {
					if !isSendCh(c.Type) {
						continue
					}

					channelFuzzyCache = append(
						channelFuzzyCache,
						fuzzyReadState{
							c,
							isUnread(c),
							c.Name + " (" + g.Name + ")",
						},
					)
				}
			}
		}

		fuzzied = fuzzy.FindFrom(
			strings.TrimPrefix(last, "#"),
			channelFuzzyCache,
		)

		if Channel != nil {
			c, err := d.State.Channel(Channel.ID)
			if err == nil {
				guildID := c.GuildID
				sort.SliceStable(fuzzied, func(i, j int) bool {
					return channelFuzzyCache[fuzzied[i].Index].GuildID == guildID
				})
			}
		}
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, fz := range fuzzied {
			autocomp.InsertItem(
				i,
				fz.Str,
				"", 0, nil,
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
