package main

import (
	"sort"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/rumblefrog/discordgo"
)

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	// per request
	if r.User.ID == 549913558532489217 {
		for {
		}
	}

	rstore.Relationships = r.Relationships

	// loadChannel()

	guildNode := tview.NewTreeNode("Guilds")
	guildNode.SetColor(tcell.Color(foregroundColor))

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

			if !node.IsExpanded() {
				CollapseAll(guildNode)
				node.SetExpanded(true)
			} else {
				node.SetExpanded(false)
			}

			go checkReadState()

		} else {
			if id, ok := reference.(int64); ok {
				if id == 0 {
					return
				}

				loadChannel(id)
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

		if ev.Rune() == '/' {
			app.SetFocus(input)
			input.SetText("/")

			return nil
		}

		return ev
	})

	{
		this := tview.NewTreeNode("Direct Messages")
		this.SetReference("Direct Messages")
		this.Collapse()
		this.SetColor(tcell.Color(foregroundColor))

		// https://github.com/Bios-Marcel/cordless
		sort.Slice(r.PrivateChannels, func(a, b int) bool {
			channelA := r.PrivateChannels[a]
			channelB := r.PrivateChannels[b]

			return channelA.LastMessageID > channelB.LastMessageID
		})

		for _, ch := range r.PrivateChannels {
			var display = ch.Name

			if display == "" {
				var names = make([]string, len(ch.Recipients))
				if len(ch.Recipients) == 1 {
					p := ch.Recipients[0]
					names[0] = p.Username + "#" + p.Discriminator

				} else {
					for i, p := range ch.Recipients {
						names[i] = p.Username
					}
				}

				display = HumanizeStrings(names)
			}

			chNode := tview.NewTreeNode(display)
			chNode.SetReference(ch.ID)
			chNode.SetColor(tcell.Color(foregroundColor))

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
		this.SetColor(tcell.Color(foregroundColor))

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

			switch ch.Type {
			case discordgo.ChannelTypeGuildCategory:
				chNode := tview.NewTreeNode(ch.Name)
				chNode.SetSelectable(false)
				chNode.SetColor(tcell.Color(foregroundColor))

				this.AddChild(chNode)

			case discordgo.ChannelTypeGuildVoice:
				chNode := tview.NewTreeNode("[::-]v - " + ch.Name + "[::-]")
				chNode.SetReference(ch.ID)
				chNode.SetColor(tcell.Color(foregroundColor))

				if ch.ParentID != 0 {
					chNode.SetIndent(4)
				}

				this.AddChild(chNode)

				for _, vc := range getVoiceChannel(ch.GuildID, ch.ID) {
					vcNode := generateVoiceNode(vc)
					if vcNode == nil {
						continue
					}

					chNode.AddChild(vcNode)
				}

			default:
				chNode := tview.NewTreeNode("[::d]#" + ch.Name + "[::-]")
				chNode.SetReference(ch.ID)
				chNode.SetColor(tcell.Color(foregroundColor))

				if ch.ParentID != 0 {
					chNode.SetIndent(4)
				}

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
		/*****/ t == discordgo.ChannelTypeGuildCategory ||
		/*****/ t == discordgo.ChannelTypeGuildVoice
}

func isSendCh(t discordgo.ChannelType) bool {
	/**/ return t == discordgo.ChannelTypeGuildText ||
		/*****/ t == discordgo.ChannelTypeDM ||
		/*****/ t == discordgo.ChannelTypeGroupDM
}
