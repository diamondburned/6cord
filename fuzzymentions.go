package main

import (
	"fmt"
	"strings"

	"github.com/diamondburned/tview"
	"github.com/sahilm/fuzzy"
)

// String returns the fuzzy search part of the struct
func (gm UserStoreArray) String(i int) string {
	var s = gm[i].Name

	if gm[i].Nick != "" {
		s += " " + gm[i].Nick
	}

	return s
}

// Len returns the length
func (gm UserStoreArray) Len() int {
	return len(gm)
}

// FuzzyMembers fuzzy searches and returns the slice of results
func FuzzyMembers(pattern string, s *UserStore, guildID int64) (fzr UserStoreArray) {
	s.RLock()
	defer s.RUnlock()

	this, ok := s.Guilds[guildID]
	if !ok {
		return
	}

	results := fuzzy.FindFrom(pattern, this)
	for i := 0; i < len(results) && i < 10; i++ {
		fzr = append(fzr, this[results[i].Index])
	}

	return
}

func fuzzyMentions(last string) {
	var fuzzied UserStoreArray

	if len(last) > 0 && Channel != nil {
		fuzzied = FuzzyMembers(
			strings.TrimPrefix(last, "@"), us, Channel.GuildID,
		)
	}

	clearList()

	if len(fuzzied) > 0 {
		g, _ := d.State.Guild(Channel.GuildID)

		for i, u := range fuzzied {
			var username = u.Name + "[::d]#" + u.Discrim + "[::-]"
			if u.Nick != "" {
				username += " [::d](" + u.Nick + ")[::-]"
			}

			if g != nil {
				for _, p := range g.Presences {
					if p.User.ID == fuzzied[i].ID {
						username = fmt.Sprintf(
							"[%s]%s[-]",
							ReflectStatusColor(p.Status),
							username,
						)
					}
				}
			}

			autocomp.InsertItem(i, &tview.ListItem{"@" + username, "", 0, nil})
		}

		rightflex.ResizeItem(autocomp, min(len(fuzzied), 10), 1)

		autofillfunc = func(i int) {
			words := strings.Fields(input.GetText())

			withoutlast := words[:len(words)-1]
			withoutlast = append(withoutlast, fmt.Sprintf(
				"<@%d> ", fuzzied[i].ID,
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
