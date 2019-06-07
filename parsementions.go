package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tview"
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

func parseMessageContent(m *discordgo.Message) (content string, emojiMap map[string][]string) {
	channel, err := d.State.Channel(m.ChannelID)
	if err != nil {
		content = ParseMentionsFallback(m)
		return
	}

	_c, emojiMap := parseEmojis(m.Content)
	content = md.Parse(_c)

	for _, user := range m.Mentions {
		var username = tview.Escape(user.Username)

		member, err := d.State.Member(channel.GuildID, user.ID)
		if err == nil && member.Nick != "" {
			username = tview.Escape(member.Nick)
		}

		var color = "[-:" + cfg.Prop.MentionColor + ":-]"
		if user.ID == d.State.User.ID {
			color = "[-:" + cfg.Prop.MentionSelfColor + ":-]"
		}

		content = strings.NewReplacer(
			// <@ID>
			fmt.Sprintf("<@%d>", user.ID),
			color+"@"+username+"[:-:]",
			// <@!ID>
			fmt.Sprintf("<@!%d>", user.ID),
			color+"@"+username+"[:-:]",
		).Replace(content)
	}

	var color = "[-:" + cfg.Prop.MentionColor + ":-]"

	for _, roleID := range m.MentionRoles {
		role, err := d.State.Role(channel.GuildID, roleID)
		if err != nil {
			continue
		}

		content = strings.Replace(
			content,
			fmt.Sprintf("<@&%d>", role.ID),
			color+"@"+role.Name+"[:-:]",
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

		return color + "#" + channel.Name + "[:-:]"
	})

	return
}
