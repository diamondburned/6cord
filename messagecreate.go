package main

import (
	"fmt"
	"time"

	"github.com/RumbleFrog/discordgo"
)

const (
	authorFormat  = "\n\n[#%06X::b]%s[-::-] [::d]%s[::-]"
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
		messageText = "pinned a message."
	case discordgo.MessageTypeRecipientAdd:
		messageText = "added " + m.Mentions[0].Username + " to the group."
	case discordgo.MessageTypeRecipientRemove:
		messageText = "removed " + m.Mentions[0].Username + " from the group."
	}

	if messageText != "" {
		setLastAuthor(0)

		messagesView.Write([]byte(
			fmt.Sprintf(
				"\n\n[::d]%s %s[::-]",
				m.Author.Username, messageText,
			),
		))

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

		scrollChat()

		messagesView.Write([]byte("[::-]"))

		setLastAuthor(m.Author.ID)
	})

}
