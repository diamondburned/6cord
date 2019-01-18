package main

import (
	"fmt"
	"log"

	"github.com/RumbleFrog/discordgo"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	ch, err := s.Channel(ChannelID) // todo: state first
	if err != nil {
		log.Panicln(err)
	}

	messagesFrame.AddText(ch.Name, true, tview.AlignLeft, tcell.ColorWhite)

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
			if LastAuthor != msg.Author.ID {
				messagesView.Write([]byte("\n"))
			}

			messagesView.Write([]byte(
				fmt.Sprintf(messageFormat, msg.Author.Username, tview.Escape(msg.Content)),
			))

			LastAuthor = msg.Author.ID
		}

		messagesView.ScrollToEnd()
	})
}
