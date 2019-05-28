package main

import (
	"html"
	"strings"

	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tview"
	"github.com/gen2brain/beeep"
)

const pingHighlightNode = " [#ED2939](!)[-]"

func mentionHandler(m *discordgo.MessageCreate) {
	// Crash-prevention
	if d.State.Settings == nil {
		return
	}

	var (
		submessage = "said in a heated channel"
		name       = m.Author.Username
		pinged     bool
	)

	if m.Author.ID != d.State.User.ID {
		if heatedChannelsExists(m.ChannelID) {
			goto Notify
		}
	}

	if !messagePingable(m.Message, m.GuildID) {
		return
	}

	pinged = true

Notify:
	if c, err := d.State.Channel(m.ChannelID); err == nil {
		switch {
		case !pinged:
			submessage = "said in a heated channel"
		case len(c.Recipients) > 0:
			submessage = "messaged"
		default:
			submessage = "mentioned you"
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
	if d.State.Settings.Status != discordgo.StatusDoNotDisturb || !pinged {
		// we ignore errors for users without dbus/notify-send
		beeep.Notify(
			name+" "+submessage,
			html.EscapeString(m.ContentWithMentionsReplaced()),
			"",
		)

		// if it's a heat signal
		if !pinged {
			return
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

		reference, ok := node.GetReference().(*discordgo.Channel)
		if !ok {
			return true
		}

		if reference.ID != m.ChannelID {
			return false
		}

		pingNode := tview.NewTreeNode(
			"[#ED2939]" + tview.Escape(name) + "[-] mentioned you",
		)

		pingNode.SetSelectable(false)
		pingNode.SetIndent(cfg.Prop.SidebarIndent - 1)

		node.AddChild(pingNode)
		node.Expand()

		if _, ok := parent.GetReference().(*discordgo.Guild); ok {
			if !strings.HasSuffix(parent.GetText(), pingHighlightNode) {
				parent.SetText(parent.GetText() + pingHighlightNode)
			}
		}

		return false
	})

	app.Draw()
}
