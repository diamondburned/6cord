package main

import (
	"log"

	"github.com/RumbleFrog/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/rivo/tview"
)

func messageAck(s *discordgo.Session, a *discordgo.MessageAck) {
	log.Println(spew.Sdump(a))

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

			for _, ugs := range d.State.UserGuildSettings {
				if ugs.GuildID == c.GuildID && !ugs.Muted {
					g, ok := parent.GetReference().(string)
					if ok {
						parent.SetText("[::b]" + g + "[::-]")
					}
				}
			}
		}

		node.SetText(name)

		return true
	})

	app.Draw()
}

func getGuildFromSettings(guildID int64) *discordgo.UserGuildSetings {
	for _, ugs
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
