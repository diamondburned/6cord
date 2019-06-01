package main

import (
	"log"
	"time"

	"github.com/diamondburned/discordgo"
)

var (
	highlightInterval = time.Duration(time.Second * 7)
	messageStore      []string
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if d == nil || Channel == nil {
		return
	}

	if rstore.Check(m.Author, RelationshipBlocked) && cfg.Prop.HideBlocked {
		return
	}

	// Notify mentions
	go mentionHandler(m)

	if m.ChannelID != Channel.ID {
		c, err := d.State.Channel(m.ChannelID)
		if err == nil {
			c.LastMessageID = m.ID

		} else {
			log.Println(err)
		}

		markUnread(m.Message)

		return
	}

	// ackMe(m.ChannelID, m.ID)

	typing.RemoveUser(&discordgo.TypingStart{
		UserID:    m.Author.ID,
		ChannelID: m.ChannelID,
	})

	// messagerenderer.go
	messageRender <- m
	scrollChat()
}
