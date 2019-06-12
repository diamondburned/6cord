package main

import (
	"fmt"
	"log"
	"strconv"
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
			imageRendererPipeline.cache.markUnfetch(m.Message)
			rendererCreate(m.Message, lastmsg)

			lastmsg = m.Message
			if m.Author.ID == d.State.User.ID {
				messagesView.ScrollToEnd()
			} else if !scrollChat() {
				messagesView.SetTitle(
					generateTitle(Channel, "[red]Unread messages[-]"),
				)
			}

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

					if (i > 1 && i == len(messageStore)-1 && strings.HasPrefix(messageStore[i-1], authorPrefix)) ||
						(i > 0 &&
							strings.HasPrefix(messageStore[i-1], authorPrefix) &&
							!strings.HasPrefix(messageStore[i+1], messageRawFormat[:3]) &&
							i != len(messageStore)-1) {

						prev = 1
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
			message, err := d.State.Message(Channel.ID, m.ID)
			if err != nil {
				Warn(err.Error())
				break
			}

			id := strconv.FormatInt(m.ID, 10)
			for i, msg := range messageStore {
				if strings.HasPrefix(msg, messageRawFormat[:3]+id+"\"]") {
					msg := messageTmpl.ExecuteString(map[string]interface{}{
						"ID":      id,
						"content": fmtMessage(message),
					})

					messageStore[i] = msg

					messagesView.SetText(strings.Join(messageStore, ""))
					break
				}
			}

		case string:
			msg := authorTmpl.ExecuteString(map[string]interface{}{
				"color": fmtHex(defaultNameColor),
				"name":  "Not Clyde",
				"time":  time.Now().Format(time.Stamp),
			})

			var (
				l = strings.Split(m, "\n")
				c []string
			)

			for i := 0; i < len(l); i++ {
				c = append(c, chatPadding+l[i])
			}

			msg += messageTmpl.ExecuteString(map[string]interface{}{
				"ID":      "0",
				"content": strings.Join(c, "\n"),
			})

			app.QueueUpdateDraw(func() {
				messagesView.Write([]byte(msg))
			})

			messageStore = append(messageStore, msg)

			scrollChat()
			lastmsg = nil

		case nil:
			messagesView.Clear()
			messageStore = make([]string, 0, prefetchMessageCount*2)
			imageRendererPipeline.cache.gc()

		default:
			Warn(fmt.Sprintf("Message renderer received event type:\n%T", i))
			log.Println(fmt.Sprintf("%#v", i))

			continue
		}

		app.Draw()
	}
}

func rendererCreate(m, lastmsg *discordgo.Message) {
	if m.Type != discordgo.MessageTypeDefault {
		var messageText string

		// https://github.com/Bios-Marcel/cordless
		switch m.Type {
		case discordgo.MessageTypeGuildMemberJoin:
			messageText = "joined the server."
		case discordgo.MessageTypeCall:
			messageText = "is calling you."
		case discordgo.MessageTypeChannelIconChange:
			messageText = "changed the channel icon."
		case discordgo.MessageTypeChannelNameChange:
			messageText = "changed the channel name to " + m.Content + "."
		case discordgo.MessageTypeChannelPinnedMessage:
			messageText = fmt.Sprintf("pinned message %d.", m.ID)
		case discordgo.MessageTypeRecipientAdd:
			messageText = "added " + m.Mentions[0].Username + " to the group."
		case discordgo.MessageTypeRecipientRemove:
			messageText = "removed " + m.Mentions[0].Username + " from the group."
		}

		if messageText != "" {
			msg := fmt.Sprintf(
				"\n\n[::d][\"%d\"]%s %s[\"\"][::-]",
				m.ID, m.Author.Username, messageText,
			)

			messagesView.Write([]byte(msg))
			messageStore = append(messageStore, msg)
		}

		return
	}

	msgFmt := messageTmpl.ExecuteString(map[string]interface{}{
		"ID":      strconv.FormatInt(m.ID, 10),
		"content": fmtMessage(m),
	})

	go func() {
		if _, err := imageRendererPipeline.cache.upd(m); err != nil {
			Message(err.Error())
		}
	}()

	if cfg.Prop.CompactMode || (lastmsg == nil ||
		(lastmsg.Author.ID != m.Author.ID || messageisOld(m, lastmsg))) {

		sentTime, err := m.Timestamp.Parse()
		if err != nil {
			sentTime = time.Now()
		}

		username, color := us.DiscordThis(m)

		msg := authorTmpl.ExecuteString(map[string]interface{}{
			"color": fmtHex(color),
			"name":  username,
			"time":  sentTime.Format(time.Stamp),
		})

		messagesView.Write([]byte(msg + msgFmt))
		messageStore = append(messageStore, msg, msgFmt)
	} else {
		messagesView.Write([]byte(msgFmt))
		messageStore = append(messageStore, msgFmt)
	}
}
