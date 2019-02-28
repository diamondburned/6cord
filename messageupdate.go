package main

import (
	"fmt"
	"strings"

	"github.com/rumblefrog/discordgo"
)

func messageUpdate(s *discordgo.Session, u *discordgo.MessageUpdate) {
	if ChannelID != u.ChannelID {
		return
	}

	if d == nil {
		return
	}

	m, err := d.State.Message(ChannelID, u.ID)
	if err != nil {
		Warn(err.Error())
		return
	}

	if rstore.Check(m.Author, RelationshipBlocked) && HideBlocked {
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

	messagesView.SetText(strings.Join(messageStore, ""))
	app.Draw()

	scrollChat()
}
