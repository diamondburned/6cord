package main

import (
	"log"
	"sort"

	"github.com/RumbleFrog/discordgo"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
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
			node.SetExpanded(!node.IsExpanded())
			return
		}

		if id, ok := reference.(int64); ok {
			ChannelID = id
			loadChannel()
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
		this.Collapse()

		guildNode.AddChild(this)

		for _, ch := range r.PrivateChannels {
			var names = make([]string, len(ch.Recipients))
			for i, p := range ch.Recipients {
				names[i] = p.Username
			}

			chNode := tview.NewTreeNode(HumanizeStrings(names))
			chNode.SetReference(ch.ID)

			this.AddChild(chNode)
		}
	}

	for _, gID := range r.Settings.GuildPositions {
		g, e := d.State.Guild(gID)
		if e != nil {
			g, e = d.Guild(gID)
			if e != nil {
				log.Println("Can't populate guild list on Guild ID:", gID, e)
				continue
			}
		}

		this := tview.NewTreeNode(g.Name)
		this.Collapse()

		sort.Slice(g.Channels, func(i, j int) bool {
			return g.Channels[i].Position < g.Channels[j].Position
		})

		for _, ch := range g.Channels {
			if !isValidCh(ch.Type) {
				continue
			}

			chNode := tview.NewTreeNode("#" + ch.Name)
			chNode.SetReference(ch.ID)

			this.AddChild(chNode)
		}

		guildNode.AddChild(this)
	}

	app.Draw()
}

func isValidCh(t discordgo.ChannelType) bool {
	/**/ return t == discordgo.ChannelTypeGuildText ||
		/*****/ t == discordgo.ChannelTypeDM ||
		/*****/ t == discordgo.ChannelTypeGroupDM
}
