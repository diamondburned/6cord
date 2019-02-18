package main

import (
	"html"

	"github.com/rumblefrog/discordgo"
	"github.com/gen2brain/beeep"
)

func mentionHandler(m *discordgo.MessageCreate) {
	// Crash-prevention
	if d.State.Settings == nil {
		return
	}

	// Skip if user is busy
	if d.State.Settings.Status == discordgo.StatusDoNotDisturb {
		return
	}

	for _, mention := range m.Mentions {
		if mention.ID == d.State.User.ID {
			var channel string
			if c, err := d.State.Channel(m.ChannelID); err == nil {
				channel = " in #" + c.Name
			}

			if err := beeep.Notify(
				m.Author.Username+" mentioned you"+channel,
				html.EscapeString(m.ContentWithMentionsReplaced()),
				"",
			); err != nil {
				Warn(err.Error())
			}

			return
		}
	}
}
