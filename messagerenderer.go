package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/diamondburned/discordgo"
)

var (
	messageRender = make(chan interface{}, 12)
)

// Function takes in a messageCreate buffer
func messageRenderer() {
	var lastmsg *discordgo.Message

	for i := range messageRender {
		switch m := i.(type) {
		case *discordgo.MessageCreate:
			if !isRegularMessage(m.Message) {
				break
			}

			rendererCreate(m.Message, lastmsg)

			lastmsg = m.Message
			scrollChat()

		case *discordgo.Message:
			rendererCreate(m, lastmsg)

			lastmsg = m
			messagesView.ScrollToEnd()

		case *discordgo.MessageDelete:
			for i := len(messageStore) - 1; i >= 0; i-- {
				if ID := getIDfromindex(i); ID != 0 {
					if m.ID != ID {
						continue
					}

					prev := 0

					if (i > 1 && i == len(messageStore)-1 && strings.HasPrefix(messageStore[i-1], authorFormat[:4])) ||
						(i > 0 &&
							strings.HasPrefix(messageStore[i-1], authorFormat[:4]) &&
							!strings.HasPrefix(messageStore[i+1], messageFormat[:3]) &&
							i != len(messageStore)-1) {

						prev = 1
						setLastAuthor(0)
					}

					messageStore = append(
						messageStore[:i-prev],
						messageStore[i+1:]...,
					)

					messagesView.SetText(strings.Join(messageStore, ""))

					break
				}
			}

			lastmsg = nil

		case *discordgo.MessageUpdate:
			for i, msg := range messageStore {
				if strings.HasPrefix(msg, fmt.Sprintf("\n"+`["%d"]`, m.ID)) {
					msg := fmt.Sprintf(
						messageFormat+"[::-]",
						m.ID, fmtMessage(m.Message),
					)

					messageStore[i] = msg

					break
				}
			}

		case nil:
			messagesView.Clear()
			messageStore = make([]string, 0, prefetchMessageCount)

		default:
			Warn(fmt.Sprintf("Message renderer received event type:\n%T", i))
			log.Println(fmt.Sprintf("%#v", i))

			continue
		}

		app.Draw()
	}
}

func rendererCreate(m, lastmsg *discordgo.Message) {
	msgFmt := fmt.Sprintf(
		messageFormat+"[::-]",
		m.ID, fmtMessage(m),
	)

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

		messagesView.Write([]byte(msg + msgFmt))
		messageStore = append(messageStore, msg, msgFmt)

	} else {
		messagesView.Write([]byte(msgFmt))
		messageStore = append(messageStore, msgFmt)
	}
}
