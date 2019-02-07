package main

import (
	"fmt"
	"log"

	"github.com/RumbleFrog/discordgo"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func loadChannel() {
	ch, err := d.Channel(ChannelID) // todo: state first
	if err != nil {
		log.Panicln(err)
	}

	messagesFrame.AddText("#"+ch.Name, true, tview.AlignLeft, tcell.ColorWhite)

	msgs, err := d.ChannelMessages(ChannelID, 75, 0, 0, 0)
	if err != nil {
		log.Panicln(err)
	}

	// reverse
	for i := len(msgs)/2 - 1; i >= 0; i-- {
		opp := len(msgs) - 1 - i
		msgs[i], msgs[opp] = msgs[opp], msgs[i]
	}

	app.QueueUpdateDraw(func() {
		for _, msg := range msgs {
			if LastAuthor != msg.Author.ID {
				messagesView.Write([]byte("\n"))
			}

			messagesView.Write([]byte(
				fmt.Sprintf(messageFormat, msg.ID, msg.Author.Username, tview.Escape(msg.Content)),
			))

			LastAuthor = msg.Author.ID
		}

		messagesView.ScrollToEnd()
	})
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
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
		}

		return ev
	})

	for _, gID := range r.Settings.GuildPositions {
		g, e := d.State.Guild(gID)
		if e != nil {
			log.Fatalln(e)
		}

		this := tview.NewTreeNode(g.Name)
		this.Collapse()

		for _, ch := range g.Channels {
			chNode := tview.NewTreeNode("#" + ch.Name)
			chNode.SetReference(ch.ID)

			this.AddChild(chNode)
		}

		guildNode.AddChild(this)
	}

	app.Draw()
}
