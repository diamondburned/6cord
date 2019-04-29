package main

import (
	"fmt"
	"time"
	"sync"

	"github.com/diamondburned/discordgo"
)

var (
	messageRender = make(chan *discordgo.Message, 3)
	msgRenderLock sync.Mutex
)

// Function takes in a messageCreate buffer
// do NOT use for update+delete (yet)
func messageRenderer() {
	var lastmsg *discordgo.Message
	msgRenderLock.Lock()
	defer msgRenderLock.Unlock()

	for m := range messageRender {
		if getLastAuthor() != m.Author.ID || (lastmsg != nil && messageisOld(m, lastmsg)) {
			sentTime, err := m.Timestamp.Parse()
			if err != nil {
				sentTime = time.Now()
			}

			setLastAuthor(m.Author.ID)

			username, color := us.DiscordThis(m)

			msg := fmt.Sprintf(
				authorFormat,
				color, username,
				sentTime.Local().Format(time.Stamp),
			)
			
			messagesView.Write([]byte(msg))
			messageStore = append(messageStore, msg)
		}

		msg := fmt.Sprintf(
			messageFormat+"[::-]",
			m.ID, fmtMessage(m),
		)

		app.QueueUpdateDraw(func() {
			messagesView.Write([]byte(msg))
		})

		messageStore = append(messageStore, msg)

		lastmsg = m
	}
}
