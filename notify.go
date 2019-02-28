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
		if c.Name != "" {
			channel = " in #" + c.Name

			m, err := d.State.Member(c.GuildID, m.Author.ID)
			if err == nil {
				if m.Nick != "" {
					name = m.Nick
				}
			}

		} else {
			var names = make([]string, len(c.Recipients))

			if len(c.Recipients) == 1 {
				p := c.Recipients[0]
				names[0] = p.Username + "#" + p.Discriminator

			} else {
				for i, p := range c.Recipients {
					names[i] = p.Username
				}
			}

			channel = HumanizeStrings(names)
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
			"[red]" + tview.Escape(name) + "[-] mentioned you",
		)

		pingNode.SetSelectable(false)

		node.AddChild(pingNode)
		node.Expand()

		if g, ok := parent.GetReference().(string); ok {
			parent.SetText("[::b]" + g + " [red](!)[-::-]")
		}

		return false
	})

	app.Draw()

}
