package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"github.com/rumblefrog/discordgo"
	"gitlab.com/diamondburned/6cord/md"
)

var (
	patternChannels = regexp.MustCompile("<#[^>]*>")
)

// ParseMentionsFallback parses mentions into strings without failing
func ParseMentionsFallback(m *discordgo.Message) (content string) {
	content = md.Parse(m.Content)

	for _, user := range m.Mentions {
		var username = tview.Escape(user.Username)

		content = strings.NewReplacer(
			// <@ID>
			fmt.Sprintf("<@%d>", user.ID),
			"[:blue:]@"+username+"[:-:]",
			// <@!ID>
			fmt.Sprintf("<@!%d>", user.ID),
			"[:blue:]@"+username+"[:-:]",
		).Replace(content)
	}

	return
}

// ParseAll parses everything into formatted strings
func ParseAll(m *discordgo.Message) (content string, emojiMap map[string][]string) {
	channel, err := d.State.Channel(m.ChannelID)
	if err != nil {
		content = ParseMentionsFallback(m)
		return
	}

	content, emojiMap = parseEmojis(md.Parse(m.Content))

	for _, user := range m.Mentions {
		var username = tview.Escape(user.Username)

		member, err := d.State.Member(channel.GuildID, user.ID)
		if err == nil && member.Nick != "" {
			username = tview.Escape(member.Nick)
		}

		content = strings.NewReplacer(
			// <@ID>
			fmt.Sprintf("<@%d>", user.ID),
			"[:blue:]@"+username+"[:-:]",
			// <@!ID>
			fmt.Sprintf("<@!%d>", user.ID),
			"[:blue:]@"+username+"[:-:]",
		).Replace(content)
	}

	for _, roleID := range m.MentionRoles {
		role, err := d.State.Role(channel.GuildID, roleID)
		if err != nil {
			continue
		}

		content = strings.Replace(
			content,
			fmt.Sprintf("<@&%d>", role.ID),
			"[:blue:]@"+role.Name+"[:-:]",
			1,
		)
	}

	content = patternChannels.ReplaceAllStringFunc(content, func(mention string) string {
		id, err := strconv.ParseInt(mention[2:len(mention)-1], 10, 64)
		if err != nil {
			return mention
		}

		channel, err := d.State.Channel(id)
		if err != nil || channel.Type == discordgo.ChannelTypeGuildVoice {
			return mention
		}

		return "[:blue:]#" + channel.Name + "[:-:]"
	})

	return
}
