package main

import (
	"log"

	"github.com/RumbleFrog/discordgo"
	"github.com/rivo/tview"
)

func messageAck(s *discordgo.Session, a *discordgo.MessageAck) {
	for _, c := range d.State.ReadState {
		if c.ID == a.ChannelID {
			c.LastMessageID = a.MessageID
		}
	}

	checkReadState()
}

func checkReadState() {
	guildView.GetRoot().Walk(func(node, parent *tview.TreeNode) bool {
		if parent == nil {
			return true
		}

		reference := node.GetReference()
		if reference == nil {
			return true
		}

		id, ok := reference.(int64)
		if !ok {
			return true
		}

		c, err := d.State.Channel(id)
		if err != nil {
			return true
		}

		var name = "[::d]" + c.Name + "[::-]"
		if isUnread(c) {
			name = "[::b]" + c.Name + "[::-]"

			g, err := d.State.Guild(c.GuildID)
			if err == nil {
				parent.SetText("[::b]" + g.Name + "[::-]")
			}
		}

		node.SetText(name)

		return true
	})

	app.Draw()
}

// true if channelID has unread msgs
func isUnread(ch *discordgo.Channel) bool {
	for _, c := range d.State.ReadState {
		if c.ID == ch.ID && c.LastMessageID != ch.LastMessageID {
			return true
		}
	}

	return false
}

var lastAck string

func ackMe(c *discordgo.Channel) {
	ack, err := d.ChannelMessageAck(
		c.ID,
		c.LastMessageID,
		lastAck,
	)

	if err != nil {
		log.Println(err)
		return
	}

	lastAck = ack.Token
}
