package main

import (
	"log"
	"sort"

	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tcell"
	"github.com/diamondburned/tview/v2"
	"gitlab.com/diamondburned/6cord/keyring"
)

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	// Store the token in a keyring
	go func() {
		log.Println("Storing token inside keyring...")
		keyring.Set(d.Token)
	}()

	rstore.Relationships = r.Relationships

	// loadChannel()

	guildView.SetSingleClick(true)

	guildNode := tview.NewTreeNode("Guilds")
	guildNode.SetColor(tcell.Color(cfg.Prop.ForegroundColor))
	guildNode.SetSelectedColor(tcell.ColorBlack)

	pNode := tview.NewTreeNode("Direct Messages")
	pNode.SetReference("Direct Messages")
	pNode.Collapse()
	pNode.SetColor(tcell.Color(cfg.Prop.ForegroundColor))
	pNode.SetSelectedColor(tcell.ColorBlack)

	// https://github.com/Bios-Marcel/cordless
	sort.Slice(r.PrivateChannels, func(a, b int) bool {
		channelA := r.PrivateChannels[a]
		channelB := r.PrivateChannels[b]

		return channelA.LastMessageID > channelB.LastMessageID
	})

	for _, ch := range r.PrivateChannels {
		var display = generateName(ch) + "[-::-]"

		if isUnread(ch) {
			display = unreadChannelColorPrefix + display
		} else {
			display = readChannelColorPrefix + display
		}

		chNode := tview.NewTreeNode(display)
		chNode.SetReference(ch)
		chNode.SetColor(tcell.Color(cfg.Prop.ForegroundColor))
		chNode.SetSelectedColor(tcell.ColorBlack)
		chNode.SetIndent(cfg.Prop.SidebarIndent - 1)

		pNode.AddChild(chNode)
	}

	guildNode.AddChild(pNode)
	guildView.SetCurrentNode(pNode)

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
		this := tview.NewTreeNode(readChannelColorPrefix + g.Name + "[-::-]")
		this.SetReference(g)
		this.Collapse()
		this.SetColor(tcell.Color(cfg.Prop.ForegroundColor))
		this.SetSelectedColor(tcell.ColorBlack)

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

			var name = generateName(ch)

			switch ch.Type {
			case discordgo.ChannelTypeGuildCategory:
				chNode := tview.NewTreeNode(unreadChannelColorPrefix + name + "[::-]")
				chNode.SetSelectable(false)
				chNode.SetColor(tcell.Color(cfg.Prop.ForegroundColor))
				chNode.SetSelectedColor(tcell.ColorBlack)
				chNode.SetIndent(cfg.Prop.SidebarIndent - 1)

				this.AddChild(chNode)

			case discordgo.ChannelTypeGuildVoice:
				chNode := tview.NewTreeNode("[-::-]" + name + "[-::-]")
				chNode.SetReference(ch)
				chNode.SetColor(tcell.Color(cfg.Prop.ForegroundColor))
				chNode.SetSelectedColor(tcell.ColorBlack)

				if ch.ParentID == 0 {
					chNode.SetIndent(cfg.Prop.SidebarIndent - 1)
				} else {
					chNode.SetIndent(cfg.Prop.SidebarIndent*2 - 1)
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
				chNode := tview.NewTreeNode(readChannelColorPrefix + name + "[-::-]")
				chNode.SetReference(ch)
				chNode.SetColor(tcell.Color(cfg.Prop.ForegroundColor))
				chNode.SetSelectedColor(tcell.ColorBlack)

				if ch.ParentID == 0 {
					chNode.SetIndent(cfg.Prop.SidebarIndent - 1)
				} else {
					chNode.SetIndent(cfg.Prop.SidebarIndent*2 - 1)
				}

				this.AddChild(chNode)
			}
		}

		checkGuildNode(g, this)
		guildNode.AddChild(this)
	}

	guildView.SetRoot(guildNode)
	guildView.SetMouseFunc(func(ev *tcell.EventMouse) bool {
		return false
	})
	guildView.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			CollapseAll(guildNode)
			node.SetExpanded(!node.IsExpanded())
			return
		}

		switch r := reference.(type) {
		case nil:

		case *discordgo.Channel:
			node.SetChildren(nil)
			loadChannel(r.ID)

		case *discordgo.Guild:
			node.SetText(readChannelColorPrefix + r.Name + "[-::-]")

			if !node.IsExpanded() {
				CollapseAll(guildNode)
				node.SetExpanded(true)
			} else {
				node.SetExpanded(false)
			}

			checkGuildNode(r, node)

		default: // Private Channels
			children := pNode.GetChildren()
			n := make([]*tview.TreeNode, 0, len(children))
			for i, c := range children {
				if c == node {
					n = append(n, c)
					n = append(n, children[:i]...)
					n = append(n, children[i+1:]...)

					pNode.SetChildren(n)
					break
				}
			}

			if !node.IsExpanded() {
				CollapseAll(guildNode)
				node.SetExpanded(true)
			} else {
				node.SetExpanded(false)
			}
		}
	})

	guildView.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyRight, tcell.KeyCtrlL:
			app.SetFocus(input)
			return nil
		case tcell.KeyLeft:
			return nil
		case tcell.KeyTab:
			return nil
		case tcell.KeyCtrlD, tcell.KeyCtrlU:
			children := guildNode.GetChildren()
			_, _, _, height := messagesView.GetInnerRect()

			hop(children, func(i int) {
				// Change fn to up/down functions accordingly
				switch ev.Key() {
				case tcell.KeyCtrlU:
					// Make sure we can scroll up by the desired amount
					if i-height/2 >= 0 {
						guildView.SetCurrentNode(children[i-height/2])
					} else {
						guildView.SetCurrentNode(children[0])
					}
				case tcell.KeyCtrlD:
					// Make sure we can scroll down by the desired amount
					if i+height/2 < len(children) {
						guildView.SetCurrentNode(children[i+height/2])
					} else {
						guildView.SetCurrentNode(children[len(children)-1])
					}
				}
			})

			return nil
		case tcell.KeyCtrlJ, tcell.KeyCtrlK,
			tcell.KeyPgDn, tcell.KeyPgUp:
			children := guildNode.GetChildren()

			hop(children, func(i int) {
				// Change fn to up/down functions accordingly
				switch ev.Key() {
				case tcell.KeyCtrlK, tcell.KeyPgUp:
					// If we're not the first guild
					if i > 0 {
						// Collapse all nodes
						CollapseAll(guildNode)
						// Set the next node as expanded
						children[i-1].SetExpanded(true)
						// Set the current node focus
						guildView.SetCurrentNode(children[i-1])
					}
				case tcell.KeyCtrlJ, tcell.KeyPgDn:
					// If we're not the last guild
					if i != len(children)-1 {
						// Collapse all nodes
						CollapseAll(guildNode)
						// Set the previous node as expanded
						children[i+1].SetExpanded(true)
						// Set the current node focus
						guildView.SetCurrentNode(children[i+1])
					}
				}
			})

			return nil
		}

		if ev.Rune() == '/' {
			app.SetFocus(input)
			input.SetText("/")

			return nil
		}

		return ev
	})

	app.Draw()
}

func hop(children []*tview.TreeNode, fn func(i int)) {
	if n := guildView.GetCurrentNode(); n != nil {
		switch r := n.GetReference().(type) {
		// If the reference is a channel, we know the cursor is over a
		// guild's children
		case *discordgo.Channel:
			// Iterate over guild nodes
			for i, gNode := range children {
				// Get the dgo guild reference
				rg, ok := gNode.GetReference().(*discordgo.Guild)
				if !ok {
					// Probably not what we're looking for, next
					continue
				}

				// Not the guild we're in
				if rg.ID != r.GuildID {
					continue // next
				}

				fn(i)
			}

		// If the reference is a guild or the direct message thing
		case *discordgo.Guild, string:
			// Iterate over guild nodes
			for i, gNode := range children {
				// If the guild node is not the node we're on
				if gNode != n {
					continue // skip
				}

				fn(i)
			}
		}
	}
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

func generateName(ch *discordgo.Channel) string {
	switch ch.Type {
	case discordgo.ChannelTypeDM, discordgo.ChannelTypeGroupDM:
		return makeDMName(ch)
	case discordgo.ChannelTypeGuildVoice:
		return "v- " + ch.Name
	case discordgo.ChannelTypeGuildCategory:
		return ch.Name
	default:
		return "#" + ch.Name
	}
}

func makeDMName(ch *discordgo.Channel) string {
	if ch.Name != "" {
		return ch.Name
	}

	var names = make([]string, len(ch.Recipients))
	if len(ch.Recipients) == 1 {
		p := ch.Recipients[0]
		names[0] = p.Username + "#" + p.Discriminator
	} else {
		for i, p := range ch.Recipients {
			names[i] = p.Username
		}
	}

	return HumanizeStrings(names)
}
