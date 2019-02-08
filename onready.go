package main

import (
	"fmt"
	"log"
	"strings"

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
		for _, m := range msgs {
			if rstore.Check(m.Author, RelationshipBlocked) {
				continue
			}

			if LastAuthor != m.Author.ID {
				messagesView.Write([]byte(
					fmt.Sprintf(authorFormat, m.Author.Username),
				))
			} else {
				messagesView.Write([]byte(
					"\n" + strings.Repeat(" ", len(m.Author.Username)+3),
				))
			}

			messagesView.Write([]byte(
				fmt.Sprintf(messageFormat, m.ID, func() string {
					var c []string
					for _, l := range strings.Split(m.Content, "\n") {
						c = append(c, tview.Escape(l))
					}

					return strings.Join(c, "\n")
				}()),
			))

			LastAuthor = m.Author.ID
		}

		messagesView.ScrollToEnd()
	})
}

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
