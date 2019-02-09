package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/RumbleFrog/discordgo"
	"github.com/rivo/tview"
)

const (
	authorFormat  = "\n\n[#%X::b]%s[-::-] ([::d]%s[::-])"
	messageFormat = "\n" + `["%d"]%s[""]`
)

var (
	highlightInterval = time.Duration(time.Second * 7)
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != ChannelID {
		return
	}

	if rstore.Check(m.Author, RelationshipBlocked) {
		return
	}

	username, color := getUserData(m.Message)

	sentTime, err := m.Timestamp.Parse()
	if err != nil {
		sentTime = time.Now()
	}

	app.QueueUpdateDraw(func() {
		if LastAuthor != m.Author.ID {
			messagesView.Write([]byte(
				fmt.Sprintf(
					authorFormat,
					color, username,
					sentTime.Format(time.Stamp),
				),
			))
		}

		messagesView.Write([]byte(
			fmt.Sprintf(
				messageFormat,
				m.ID, fmtMessage(m.Message),
			),
		))

		messagesView.ScrollToEnd()

		messagesView.Write([]byte("[-:-:-]"))

		LastAuthor = m.Author.ID
	})
}

func messageUpdate(s *discordgo.Session, u *discordgo.MessageUpdate) {
	if ChannelID != u.ChannelID {
		return
	}

	m, err := d.State.Message(ChannelID, u.ID)
	if err != nil {
		log.Println(err)
		return
	}

	username, _ := getUserData(m)

	app.QueueUpdateDraw(func() {
		messagesView.Write([]byte(
			fmt.Sprintf(
				"\n\n"+`[::d]%s edited message ID %d:`+"\n",
				username, u.ID,
			),
		))

		messagesView.Highlight(fmt.Sprintf("%d", u.ID))
	})

	st := fmtMessage(m) + "[::-][\"\"]"
	app.QueueUpdateDraw(func() {
		messagesView.Write([]byte(st))
	})

	time.Sleep(highlightInterval)
	app.QueueUpdateDraw(func() {
		messagesView.Highlight()
	})
}

func fmtMessage(m *discordgo.Message) string {
	var (
		ct     = m.ContentWithMentionsReplaced()
		edited string
		c      []string
		l      = strings.Split(ct, "\n")
	)

	for i := 0; i < len(l); i++ {
		c = append(c, "\t"+tview.Escape(l[i]))
	}

	if len(m.Attachments) > 0 {
		if m.Content != "" {
			c = append(c, "\n")
		}

		for _, a := range m.Attachments {
			c = append(c, "\t"+tview.Escape(a.URL))
		}
	}

	if m.EditedTimestamp != "" {
		edited = " [::d](edited)[::-]"
	}

	return strings.Join(c, "\n") + edited
}
