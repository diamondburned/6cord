package main

import (
	"fmt"

	"github.com/RumbleFrog/discordgo"
	"github.com/rivo/tview"
)

const (
	messageFormat = "\n[\"%d\"][::b]%s  [::-]%s[\"\"]"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != ChannelID {
		return
	}

	app.QueueUpdateDraw(func() {
		if LastAuthor != m.Author.ID {
			messagesView.Write([]byte("\n"))
		}

		messagesView.Write([]byte(
			fmt.Sprintf(messageFormat, m.ID, m.Author.Username, tview.Escape(m.Content)),
		))

		messagesView.ScrollToEnd()

		LastAuthor = m.Author.ID
	})
}
