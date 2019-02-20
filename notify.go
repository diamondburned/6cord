package main

import (
	"html"

	"github.com/gen2brain/beeep"
	"github.com/rivo/tview"
	"github.com/rumblefrog/discordgo"
)

func mentionHandler(m *discordgo.MessageCreate) {
	// Crash-prevention
	if d.State.Settings == nil {
		return
	}

	if m.GuildID != 0 {
		for _, mention := range m.Mentions {
			if mention.ID == d.State.User.ID {
				goto PingConfirmed
			}
		}
	} else {
		if m.Author.ID != d.State.User.ID {
			goto PingConfirmed
		}
	}

	return

PingConfirmed:
	var channel string
	var name = m.Author.Username

	if c, err := d.State.Channel(m.ChannelID); err == nil {
		channel = " in #" + c.Name

		m, err := d.State.Member(c.GuildID, m.Author.ID)
		if err == nil {
			if m.Nick != "" {
				name = m.Nick
			}
		}
	}

	// Skip if user is busy
	if d.State.Settings.Status != discordgo.StatusDoNotDisturb {
		if err := beeep.Notify(
			name+" mentioned you"+channel,
			html.EscapeString(m.ContentWithMentionsReplaced()),
			"",
		); err != nil {
			Warn(err.Error())
		}
	}

	// Walk the tree for the sake of a (1)

	if m.ChannelID == ChannelID {
		return
	}

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

		if id != m.ChannelID {
			return true
		}

		pingNode := tview.NewTreeNode(
			"[red]" + name + "[-] mentioned you",
		)

		pingNode.SetSelectable(false)

		node.AddChild(pingNode)
		node.Expand()

		return false
	})

	app.Draw()

}
