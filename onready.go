package main

import (
	"log"

	"github.com/RumbleFrog/discordgo"
)

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	msgs, err := s.ChannelMessages(ChannelID, 75, 0, 0, 0)
	if err != nil {
		log.Panicln(err)
	}

	// reverse
	for i := len(msgs)/2 - 1; i >= 0; i-- {
		opp := len(msgs) - 1 - i
		msgs[i], msgs[opp] = msgs[opp], msgs[i]
	}

	app.QueueUpdateDraw(func() {
		for _, msg := range msgs {
			messagesList.AddItem(
				msg.Author.Username,
				msg.Content,
				0, nil,
			)
		}

		messagesList.SetCurrentItem(
			messagesList.GetItemCount(),
		)
	})
}
