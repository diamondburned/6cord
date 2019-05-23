package main

import (
	"log"
	"strings"

	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tview"
)

const readChannelColorPrefix = "[#808080::]"

func messageAck(s *discordgo.Session, a *discordgo.MessageAck) {
	// Sets ReadState to the message you read
	for _, c := range d.State.ReadState {
		if c.ID == a.ChannelID && c.LastMessageID != 0 {
			c.LastMessageID = a.MessageID
		}
	}

	// update
	checkReadState(a.ChannelID)
}

// "[::b]actual string[::-]"
func stripFormat(a string) string {
	if len(a) <= 10 {
		return a
	}

	if strings.HasPrefix(a, readChannelColorPrefix) {
		a = a[len(readChannelColorPrefix):]
	}

	return strings.TrimSuffix(a, "[-::-]")
}

func checkReadState(chID ...int64) {
	var guildSettings *discordgo.UserGuildSettings

	if d.State == nil {
		return
	}

	if d.State.Settings == nil {
		return
	}

	if guildView == nil {
		return
	}

	changed := false

	root := guildView.GetRoot()
	if root == nil {
		return
	}

	root.Walk(func(node, parent *tview.TreeNode) bool {
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

		// This is true when the current node is a voice state
		// AsÜ the voice state is the channel's children, the channel (parent)
		// will have an int64 reference
		_, ok = parent.GetReference().(int64)
		if ok {
			return true
		}

		if len(chID) > 0 {
			for _, chid := range chID {
				if chid == id {
					node.ClearChildren()
					g, ok := parent.GetReference().(string)
					if ok {
						if !strings.HasPrefix(g, readChannelColorPrefix) {
							parent.SetText(readChannelColorPrefix + g + "[-::-]")
						}
					}
				}
			}

			return true
		}

		c, err := d.State.Channel(id)
		if err != nil {
			return true
		}

		if guildSettings == nil || guildSettings.GuildID != c.GuildID {
			guildSettings = getGuildFromSettings(c.GuildID)
		}

		var (
			chSettings   = getChannelFromGuildSettings(c.ID, guildSettings)
			originalName = stripFormat(node.GetText())
		)

		name := readChannelColorPrefix + originalName + "[-::-]"

		var (
			chMuted = settingChannelIsMuted(chSettings, guildSettings)
			guMuted = settingGuildIsMuted(guildSettings)
		)

		if isUnread(c) && !chMuted {
			changed = true

			name = "[::b]" + originalName + "[-::-]"

			if !guMuted {
				g, ok := parent.GetReference().(string)
				if ok {
					if !strings.HasSuffix(parent.GetText(), " [#DC143C](!)[-::-]") {
						parent.SetText("[::b]" + g + "[-::-]")
					}
				}
			}
		}

		node.SetText(name)

		return true
	})

	if changed == true {
		app.Draw()
	}
}

// true if channelID has unread msgs
func isUnread(ch *discordgo.Channel) bool {
	if ch.LastMessageID == 0 {
		return false
	}

	for _, c := range d.State.ReadState {
		if c.ID != ch.ID {
			continue
		}

		if c.LastMessageID != ch.LastMessageID {
			return true
		}
	}

	return false
}

var lastAck string

func ackMe(m *discordgo.Message) {
	c, err := d.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	if !isUnread(c) {
		return
	}

	// triggers messageAck
	a, err := d.ChannelMessageAck(m.ChannelID, m.ID, lastAck)

	if err != nil {
		log.Println(err)
		return
	}

	lastAck = a.Token
}
