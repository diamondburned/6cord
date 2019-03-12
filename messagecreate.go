package main

import (
	"fmt"
	"log"
	"time"

	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tview"
)

const (
	authorFormat  = "\n\n[#%06X::]%s[-::] [::d]%s[::-]"
	messageFormat = "\n" + `["%d"]%s ["ENDMESSAGE"]`
)

var (
	highlightInterval = time.Duration(time.Second * 7)
	messageStore      []string
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if d == nil || Channel == nil {
		return
	}

	if rstore.Check(m.Author, RelationshipBlocked) && cfg.Prop.HideBlocked {
		return
	}

	// Notify mentions
	go mentionHandler(m)

	if m.ChannelID != Channel.ID {
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

	if m.Author.ID != d.State.User.ID {
		ackMe(m.Message)
	}

	typing.RemoveUser(&discordgo.TypingStart{
		UserID:    m.Author.ID,
		ChannelID: m.ChannelID,
	})

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

	if getLastAuthor() != m.Author.ID {
		setLastAuthor(m.Author.ID)

		username, color := us.DiscordThis(m.Message)

		msg := fmt.Sprintf(
			authorFormat,
			color, tview.Escape(username),
			sentTime.Format(time.Stamp),
		)

		app.QueueUpdate(func() {
			messagesView.Write([]byte(msg))
		})

		messageStore = append(messageStore, msg)
	}

	msg := fmt.Sprintf(
		messageFormat+"[::-]",
		m.ID, fmtMessage(m.Message),
	)

	app.QueueUpdateDraw(func() {
		messagesView.Write([]byte(msg))
	})

	messageStore = append(messageStore, msg)

	scrollChat()
}
