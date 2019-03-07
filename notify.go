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

	if !messagePingable(m.Message, m.GuildID) {
		return
	}

	var submessage = "mentioned you"
	var name = m.Author.Username

	if c, err := d.State.Channel(m.ChannelID); err == nil {
		if len(c.Recipients) > 0 {
			submessage = "messaged you"
		}

		if c.Name != "" {
			submessage += " in #" + c.Name

			m, err := d.State.Member(c.GuildID, m.Author.ID)
			if err == nil {
				if m.Nick != "" {
					name = m.Nick
				}
			}

		} else {
			if len(c.Recipients) > 1 {
				var names = make([]string, len(c.Recipients))

				for i, p := range c.Recipients {
					names[i] = p.Username
				}

				submessage += " in " + HumanizeStrings(names)
			}
		}
	}

	// Skip if user is busy
	if d.State.Settings.Status != discordgo.StatusDoNotDisturb {
		if err := beeep.Notify(
			name+" "+submessage,
			html.EscapeString(m.ContentWithMentionsReplaced()),
			"",
		); err != nil {
			Warn(err.Error())
		}
	}

	// Walk the tree for the sake of a (1)

	if Channel != nil && m.ChannelID == Channel.ID {
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
