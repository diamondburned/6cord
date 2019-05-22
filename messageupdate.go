package main

import (
	"github.com/diamondburned/discordgo"
)

func messageUpdate(s *discordgo.Session, u *discordgo.MessageUpdate) {
	if d == nil || Channel == nil {
		return
	}

	if Channel.ID != u.ChannelID {
		return
	}

	m, err := d.State.Message(Channel.ID, u.ID)
	if err != nil {
		Warn(err.Error())
		return
	}

	if rstore.Check(m.Author, RelationshipBlocked) && cfg.Prop.HideBlocked {
		return
	}

	messageRender <- u
}
