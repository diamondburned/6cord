package main

import (
	"fmt"
	"strings"

	"github.com/RumbleFrog/discordgo"
	"github.com/rivo/tview"
)

const (
	authorFormat = "\n\n[::b]%s"

	messageFormat = "\n[\"%d\"][::-]%s[\"\"]"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != ChannelID {
		return
	}

	app.QueueUpdateDraw(func() {
		if LastAuthor != m.Author.ID {
			messagesView.Write([]byte(
				fmt.Sprintf(authorFormat, m.Author.Username),
			))
		}

		messagesView.Write([]byte(
			fmt.Sprintf(messageFormat, m.ID, func() string {
				var c string
				for _, l := range strings.Split(m.Content, "\n") {
					c += "\t\t" + tview.Escape(l) + "\n"
				}

				return c
			}()),
		))

		messagesView.ScrollToEnd()

		LastAuthor = m.Author.ID
	})
}
