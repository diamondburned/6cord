package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/RumbleFrog/discordgo"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	ch, err := s.Channel(ChannelID) // todo: state first
	if err != nil {
		log.Panicln(err)
	}

	messagesFrame.AddText("\t#"+ch.Name, true, tview.AlignLeft, tcell.ColorWhite)

	msgs, err := s.ChannelMessages(ChannelID, 75, 0, 0, 0)
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
			i := messagesView.GetRowCount()
			setMessage(msg, i)
		}

		messagesView.ScrollToEnd()

	})
}

func setMessage(msg *discordgo.Message, i int) (n int) {
	if msg.Author != nil {
		messagesView.SetCellSimple(i, 0, msg.Author.Username)
	}

	lines := strings.Split(msg.Content, "\n")

	c := tview.NewTableCell(lines[0]).SetExpansion(1)

	messagesView.SetCell(i, 1, c)

	for e := 1; e < len(lines); e++ {
		c := tview.NewTableCell(lines[e]).SetExpansion(1)

		messagesView.InsertRow(i + e)
		messagesView.SetCell(i+e, 1, c)
		messagesView.SetCellSimple(i+e, 2, strconv.FormatInt(msg.ID, 10))
	}

	n = len(lines)

	messagesView.SetCellSimple(i, 2, strconv.FormatInt(msg.ID, 10))

	LastAuthor = msg.Author.ID

	return
}
