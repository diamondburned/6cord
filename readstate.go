package main

import (
	"github.com/RumbleFrog/discordgo"
	"github.com/rivo/tview"
)

func checkReadState() {
	guildView.GetRoot().Walk(func(node, parent *tview.TreeNode) bool {
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
