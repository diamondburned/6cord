package main

import (
	"fmt"
	"strings"

	"github.com/RumbleFrog/discordgo"
)

func messageUpdate(s *discordgo.Session, u *discordgo.MessageUpdate) {
	if ChannelID != u.ChannelID {
		return
	}

	m, err := d.State.Message(ChannelID, u.ID)
	if err != nil {
		Warn(err.Error())
		return
	}

	if rstore.Check(m.Author, RelationshipBlocked) {
		return
	}

	// username, _ := us.DiscordThis(m)

	for i, msg := range messageStore {
		if strings.HasPrefix(msg, fmt.Sprintf("\n"+`["%d"]`, u.ID)) {
			msg := fmt.Sprintf(
				messageFormat+"[::-]",
				m.ID, fmtMessage(m),
			)

			messageStore[i] = msg

			break
		}
	}

	messagesView.Clear()
	messagesView.SetText(strings.Join(messageStore, ""))

	app.Draw()

	scrollChat()

	return

	//app.QueueUpdateDraw(func() {
	//messagesView.Write([]byte(
	//fmt.Sprintf(
	//"\n\n"+`[::d]%s edited message ID %d:`+"\n",
	//username, u.ID,
	//),
	//))

	//messagesView.Highlight(fmt.Sprintf("%d", u.ID))
	//})

	//st := fmtMessage(m) + "[::-][\"\"]\n"
	//app.QueueUpdateDraw(func() {
	//messagesView.Write([]byte(st))
	//})

	//time.Sleep(highlightInterval)
	//app.QueueUpdateDraw(func() {
	//messagesView.Highlight()
	//scrollChat()
	//})

	//setLastAuthor(0)
}
