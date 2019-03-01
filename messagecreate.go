package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rivo/tview"
	"github.com/rumblefrog/discordgo"
)

const (
	authorFormat  = "\n\n[#%06X::b]%s[-::-] [::d]%s[::-]"
	messageFormat = "\n" + `["%d"]%s ["ENDMESSAGE"]`
)

var (
	highlightInterval = time.Duration(time.Second * 7)
	messageStore      []string
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if d == nil {
		return
	}

	if rstore.Check(m.Author, RelationshipBlocked) && HideBlocked {
		return
	}

	// Notify mentions
	go mentionHandler(m)

	if m.ChannelID != ChannelID {
		c, err := d.State.Channel(m.ChannelID)
		if err == nil {
			c.LastMessageID = m.ID

		} else {
			log.Println(err)
		}

		checkReadState()

		return
	}

	if !isRegularMessage(m.Message) {
		return
	}

	if len(m.Embeds) == 1 {
		m := m.Embeds[0]
		// edgiest case ever
		if m.Description == "" && m.Title == "" && len(m.Fields) == 0 {
			return
		}
	}

	sentTime, err := m.Timestamp.Parse()
	if err != nil {
		sentTime = time.Now()
	}

	app.QueueUpdateDraw(func() {
		if getLastAuthor() != m.Author.ID {
			username, color := us.DiscordThis(m.Message)

			msg := fmt.Sprintf(
				authorFormat,
				color, tview.Escape(username),
				sentTime.Format(time.Stamp),
			)

			messagesView.Write([]byte(msg))
			messageStore = append(messageStore, msg)
		}

		msg := fmt.Sprintf(
			messageFormat+"[::-]",
			m.ID, fmtMessage(m.Message),
		)

		messagesView.Write([]byte(msg))
		messageStore = append(messageStore, msg)

		scrollChat()

		setLastAuthor(m.Author.ID)
	})

}
