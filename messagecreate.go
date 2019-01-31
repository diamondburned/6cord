package main

import (
	"strconv"

	"github.com/RumbleFrog/discordgo"
)

const (
	messageFormat = "  %s"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != ChannelID {
		return
	}

	app.QueueUpdateDraw(func() {
		i := messagesView.GetRowCount()
		setMessage(m.Message, i)

		messagesView.ScrollToEnd()
	})
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.ChannelID != ChannelID {
		return
	}

	edited := false

	IDstring := strconv.FormatInt(m.ID, 10)

	for i := 0; i < messagesView.GetRowCount(); i++ {
		if c := messagesView.GetCell(i, 2); c != nil {
			if edited && (c.Text == IDstring) {
				messagesView.RemoveRow(i)
				continue
			}

			if c.Text == IDstring {
				messagesView.InsertRow(i)

				// At this point, i is the empty row.
				// We will edit the i row
				n := setMessage(m.Message, i)

				// Now we have leftover rows that belong to the old message.
				// These rows should have our message ID.
				messagesView.RemoveRow(n + i + 1)

				edited = true
				continue
			}
		}
	}

	app.Draw()
}
