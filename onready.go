package main

import (
	"sort"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/rumblefrog/discordgo"
)

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	rstore.Relationships = r.Relationships

	// loadChannel()

	guildNode := tview.NewTreeNode("Guilds")

	guildView.SetRoot(guildNode)
	guildView.SetCurrentNode(guildNode)
	guildView.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			CollapseAll(guildNode)
			node.SetExpanded(!node.IsExpanded())
			return
		}

		if name, ok := reference.(string); ok {
			node.SetText("[::d]" + name + "[::-]")

			CollapseAll(guildNode)
			node.SetExpanded(!node.IsExpanded())

		} else {
			if id, ok := reference.(int64); ok {
				if id == 0 {
					return
				}

				ChannelID = id
				loadChannel()
			}
		}
	})

	guildView.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyRight:
			app.SetFocus(input)
			return nil
		case tcell.KeyLeft:
			return nil
		case tcell.KeyTab:
			return nil
		}

		return ev
	})

	{
		this := tview.NewTreeNode("Direct Messages")
		this.SetReference("Direct Messages")
		this.Collapse()

		// https://github.com/Bios-Marcel/cordless
		sort.Slice(r.PrivateChannels, func(a, b int) bool {
			channelA := r.PrivateChannels[a]
			channelB := r.PrivateChannels[b]

			return channelA.LastMessageID > channelB.LastMessageID
		})

		for _, ch := range r.PrivateChannels {
			var names = make([]string, len(ch.Recipients))
			for i, p := range ch.Recipients {
				if p.Username == "" {
					continue
				}

				names[i] = p.Username + "#" + p.Discriminator
			}

			chNode := tview.NewTreeNode(HumanizeStrings(names))
			chNode.SetReference(ch.ID)

			this.AddChild(chNode)
		}

		guildNode.AddChild(this)
	}

	// https://github.com/Bios-Marcel/cordless
	sort.Slice(r.Guilds, func(a, b int) bool {
		aFound := false
		for _, guild := range r.Settings.GuildPositions {
			if aFound {
				if guild == r.Guilds[b].ID {
					return true
				}
			} else {
				if guild == r.Guilds[a].ID {
					aFound = true
				}
			}
		}

		return false
	})

	for _, g := range r.Guilds {
		this := tview.NewTreeNode("[::d]" + g.Name + "[::-]")
		this.SetReference(g.Name)
		this.Collapse()

		sorted := SortChannels(g.Channels)

		for _, ch := range sorted {
			if !isValidCh(ch.Type) {
				continue
			}

			perm, err := d.State.UserChannelPermissions(
				d.State.User.ID,
				ch.ID,
			)

			if err != nil {
				continue
			}

			if perm&discordgo.PermissionReadMessages == 0 {
				continue
			}

			if ch.Type == discordgo.ChannelTypeGuildCategory {
				chNode := tview.NewTreeNode("> " + ch.Name)
				chNode.SetReference(0)

				this.AddChild(chNode)
			} else {
				chNode := tview.NewTreeNode("[::d]#" + ch.Name + "[::-]")
				chNode.SetReference(ch.ID)

				this.AddChild(chNode)
			}
		}

		guildNode.AddChild(this)
	}

	app.Draw()

	checkReadState()
}

func isValidCh(t discordgo.ChannelType) bool {
	/**/ return t == discordgo.ChannelTypeGuildText ||
		/*****/ t == discordgo.ChannelTypeDM ||
		/*****/ t == discordgo.ChannelTypeGroupDM ||
		/*****/ t == discordgo.ChannelTypeGuildCategory
}
