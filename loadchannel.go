package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/RumbleFrog/discordgo"
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

			sentTime, err := m.Timestamp.Parse()
			if err != nil {
				sentTime = time.Now()
			}

			var msg string
			if getLastAuthor() != m.Author.ID {
				username, color := us.DiscordThis(m)

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

			setLastAuthor(m.Author.ID)
		}(m, i)
	}

	wg.Wait()

	messagesView.Clear()
	messagesView.Write([]byte(
		strings.Join(messages, ""),
	))

	messagesView.ScrollToEnd()

	app.SetFocus(input)
}

func scrollChat() {
	if !messagesView.HasFocus() {
		messagesView.ScrollToEnd()
	}
}
