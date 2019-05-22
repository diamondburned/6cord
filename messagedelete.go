package main

import (
	"github.com/diamondburned/discordgo"
)

func messageDelete(s *discordgo.Session, rm *discordgo.MessageDelete) {
	if d == nil || Channel == nil {
		return
	}

	if rm.ChannelID != Channel.ID {
		return
	}

	messageRender <- rm
}

func messageDeleteBulk(s *discordgo.Session, rmb *discordgo.MessageDeleteBulk) {
	if d == nil || Channel == nil {
		return
	}

	if rmb.ChannelID != Channel.ID {
		return
	}

	for _, m := range rmb.Messages {
		messageRender <- &discordgo.MessageDelete{
			Message: &discordgo.Message{
				ChannelID: rmb.ChannelID,
				ID:        m,
			},
		}
	}
}
