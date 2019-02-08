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
	authorFormat = "\n\n[::b]%s   "

	messageFormat = `[::-]["%d"]%s[""]`
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

	app.QueueUpdateDraw(func() {
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
			fmt.Sprintf(
				messageFormat,
				m.ID, spaceIndent(m.Author, m.Content, false),
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

	app.QueueUpdateDraw(func() {
		messagesView.Write([]byte(
			fmt.Sprintf(
				"[\"EDIT_MESSAGE\"]\n"+`[::d]User edited message ID %d:`+"\n",
				u.ID,
			),
		))

		// Disabled until the highlighting bug is fixed:
		// https://media.discordapp.net/attachments/361916911682060289/543338561089437706/unknown.png
		//messagesView.Highlight()

		// Ugh who am I kidding
		messagesView.Highlight(fmt.Sprintf("%d", u.ID))
	})

	m, err := d.State.Message(ChannelID, u.ID)
	if err != nil {
		log.Println(err)
		return
	}

	st := spaceIndent(m.Author, u.Content, true) + "[::-][\"\"]"
	app.QueueUpdateDraw(func() {
		messagesView.Write([]byte(st))
	})

	time.Sleep(highlightInterval)
	app.QueueUpdateDraw(func() {
		messagesView.Highlight()
	})
}

func spaceIndent(au *discordgo.User, ct string, editmode bool) string {
	var (
		c []string
		l = strings.Split(ct, "\n")
	)

	if !editmode {
		c = append(c, tview.Escape(l[0]))
		if len(l) <= 1 {
			goto Done
		}
	}

	{ // TODO: CLEAN UP THIS MESS PLEASE!
		var (
			a   = 1
			sfx = ""
		)

		if editmode {
			a = 0
			sfx += ">"
		}

		sp := strings.Repeat(" ", len(au.Username)+3) + sfx

		for i := a; i < len(l); i++ {
			c = append(c, sp+tview.Escape(l[i]))
		}
	}

Done:
	return strings.Join(c, "\n")
}
