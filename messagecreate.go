package main

import "github.com/RumbleFrog/discordgo"

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != ChannelID {
		return
	}

	app.QueueUpdateDraw(func() {
		messagesList.AddItem(
			m.Author.Username,
			m.Content,
			0, nil,
		)

		messagesList.SetCurrentItem(
			messagesList.GetItemCount(),
		)
	})
}
