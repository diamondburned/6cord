package main

import (
	"fmt"
	"time"

	"github.com/RumbleFrog/discordgo"
)

func messageUpdate(s *discordgo.Session, u *discordgo.MessageUpdate) {
	if ChannelID != u.ChannelID {
		return
	}

	m, err := d.State.Message(ChannelID, u.ID)
	if err != nil {
		Warn(err.Error())
		return
	}

	// self-bots? not today.
	if m.Author.Bot && len(m.Embeds) > 0 {
		return
	}

	if rstore.Check(m.Author, RelationshipBlocked) {
		return
	}

	username, _ := us.DiscordThis(m)

	app.QueueUpdateDraw(func() {
		messagesView.Write([]byte(
			fmt.Sprintf(
				"\n\n"+`[::d]%s edited message ID %d:`+"\n",
				username, u.ID,
			),
		))

		messagesView.Highlight(fmt.Sprintf("%d", u.ID))
	})

	st := fmtMessage(m) + "[::-][\"\"]\n"
	app.QueueUpdateDraw(func() {
		messagesView.Write([]byte(st))
	})

	time.Sleep(highlightInterval)
	app.QueueUpdateDraw(func() {
		messagesView.Highlight()
		scrollChat()
	})

	setLastAuthor(0)
}
