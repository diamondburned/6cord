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

			if i > 0 && len(messageStore)-1 > i {
				if (strings.HasPrefix(messageStore[i-1], "\n\n[#")) &&
					!strings.HasPrefix(messageStore[i+1], "\n\n[#") {

					prev = 1
					setLastAuthor(0)
				}
			}

			if i+1 == len(messageStore) {
				messageStore = messageStore[:i-prev]
			} else {
				messageStore = append(
					messageStore[:i-prev],
					messageStore[i+1:]...,
				)
			}

			app.QueueUpdate(func() {
				messagesView.SetText(strings.Join(messageStore, ""))
			})

			scrollChat()

			return
		}
	}
	/*
		if rm.Content != "" {
			Message("Can't delete message with content: " + rm.Content)
		} else {
			Message(fmt.Sprintf("Can't delete message with ID: %d", rm.ID))
		}
	*/
}

func messageDeleteBulk(s *discordgo.Session, rmb *discordgo.MessageDeleteBulk) {
	for _, m := range rmb.Messages {
		messageDelete(s, &discordgo.MessageDelete{
			Message: &discordgo.Message{
				ChannelID: rmb.ChannelID,
				ID:        m,
			},
		})
	}
}
