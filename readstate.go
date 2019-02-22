package main

import (
	"log"
	"strings"

	"github.com/rivo/tview"
	"github.com/rumblefrog/discordgo"
)

func messageAck(s *discordgo.Session, a *discordgo.MessageAck) {
	// Sets ReadState to the message you read
	for _, c := range d.State.ReadState {
		if c.ID == a.ChannelID {
			c.LastMessageID = a.MessageID
		}
	}

	ackData[a.ChannelID] = a.MessageID

	// update
	checkReadState(a.ChannelID)
}

// "[::b]actual string[::-]"
func stripFormat(a string) string {
	if len(a) <= 10 {
		return a
	}

	if strings.HasPrefix(a, "[::") {
		a = a[5:]
	}

	return strings.TrimSuffix(a, "[::-]")
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
		// As the voice state is the channel's children, the channel (parent)
		// will have an int64 reference
		_, ok = parent.GetReference().(int64)
		if ok {
			return true
		}

		for _, chid := range chID {
			if chid == id {
				node.ClearChildren()
			}
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

		name := "[::d]" + originalName + "[::-]"

		var (
			chMuted = settingChannelIsMuted(chSettings)
			guMuted = settingGuildIsMuted(guildSettings)
		)

		if isUnread(c) && !chMuted {
			changed = true

			name = "[::b]" + originalName + "[::-]"

			if !guMuted {
				g, ok := parent.GetReference().(string)
				if ok {
					if !strings.HasPrefix(g, "[::b]") {
						parent.SetText("[::b]" + g + "[::-]")
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
	for _, c := range d.State.ReadState {
		if c.ID == ch.ID && c.LastMessageID != ch.LastMessageID {
			return true
		}
	}

	return false
}

var (
	lastAck string

	// map[channelID]messageID
	ackData = make(map[int64]int64)
)

func ackMe(c *discordgo.Channel, m *discordgo.Message) {
	if mID, ok := ackData[m.ChannelID]; ok {
		if mID == m.ID {
			return
		}
	}

	// triggers messageAck
	ack, err := d.ChannelMessageAck(c.ID, m.ID, lastAck)

	if err != nil {
		log.Println(err)
		return
	}

	lastAck = ack.Token
}
