package main

import (
	"fmt"
	"strings"

	"github.com/diamondburned/discordgo"
)

func messageUpdate(s *discordgo.Session, u *discordgo.MessageUpdate) {
	if d == nil || Channel == nil {
		return
	}

	if Channel.ID != u.ChannelID {
		return
	}

	m, err := d.State.Message(Channel.ID, u.ID)
	if err != nil {
		Warn(err.Error())
		return
	}

	if rstore.Check(m.Author, RelationshipBlocked) && cfg.HideBlocked {
		return
	}

	// username, _ := us.DiscordThis(m)

	for i, msg := range messageStore {
		if strings.HasPrefix(msg, fmt.Sprintf("\n"+`["%d"]`, u.ID)) {
			msg := fmt.Sprintf(
				messageFormat+"[::-]",
				m.ID, fmtMessage(m),
			)

			messageStore[i] = msg

			break
		}
	}

	app.QueueUpdateDraw(func() {
		messagesView.SetText(strings.Join(messageStore, ""))
	})

	scrollChat()
}
