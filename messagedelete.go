package main

import (
	"strings"

	"github.com/rumblefrog/discordgo"
)

func messageDelete(s *discordgo.Session, rm *discordgo.MessageDelete) {
	if rm.ChannelID != ChannelID {
		return
	}

	for i := len(messageStore) - 1; i >= 0; i-- {
		if ID := getIDfromindex(i); ID != 0 {
			if rm.ID != ID {
				continue
			}

			prev := 0

			if i > 0 {
				if strings.HasPrefix(messageStore[i-1], "\n\n[#") {
					prev = 1
					setLastAuthor(0)
				}
			}

			messageStore = append(
				messageStore[:i-prev],
				messageStore[i+1:]...,
			)

			messagesView.SetText(strings.Join(messageStore, ""))
			app.Draw()

			scrollChat()

			return
		}
	}

	Message("Can't delete message with content: " + rm.Content)
}
