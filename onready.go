package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/RumbleFrog/discordgo"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func loadChannel() {
	ch, err := d.Channel(ChannelID) // todo: state first
	if err != nil {
		log.Panicln(err)
	}

	wrapFrame.SetTitle("#" + ch.Name)
	typing.Reset()

	msgs, err := d.ChannelMessages(ChannelID, 75, 0, 0, 0)
	if err != nil {
		log.Panicln(err)
	}

	// reverse
	for i := len(msgs)/2 - 1; i >= 0; i-- {
		opp := len(msgs) - 1 - i
		msgs[i], msgs[opp] = msgs[opp], msgs[i]
	}

	var (
		messages = make([]string, len(msgs))
		wg       sync.WaitGroup
	)

	for i, m := range msgs {
		wg.Add(1)
		go func(m *discordgo.Message, i int) {
			defer wg.Done()

			if rstore.Check(m.Author, RelationshipBlocked) {
				return
			}

			username, color := getUserData(m)

			sentTime, err := m.Timestamp.Parse()
			if err != nil {
				sentTime = time.Now()
			}

			var msg string
			if LastAuthor != m.Author.ID {
				msg = fmt.Sprintf(
					authorFormat,
					color, username,
					sentTime.Format(time.Stamp),
				)
			}

			msg += fmt.Sprintf(
				messageFormat,
				m.ID, fmtMessage(m),
			)

			messages[i] = msg

			LastAuthor = m.Author.ID
		}(m, i)
	}

	wg.Wait()

	messagesView.Write([]byte(
		strings.Join(messages, ""),
	))

	messagesView.ScrollToEnd()

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
