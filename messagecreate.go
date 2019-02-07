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

	IDstring := strconv.FormatInt(m.ID, 10)

	for i := 0; i < messagesView.GetRowCount(); i++ {
		if c := messagesView.GetCell(i, 2); c != nil {
			if c.Text == IDstring {
				messagesView.InsertRow(i)

				// At this point, i is the empty row.
				// We will edit the i row
				n := setMessage(m.Message, i)

				for e := n - 1 + i; i < messagesView.GetRowCount(); e++ {
					if c := messagesView.GetCell(i, 2); c != nil {
						if c.Text == IDstring {
							messagesView.RemoveRow(e)
						}
					}
				}

				app.Draw()

				return
			}
		}
	}

}
