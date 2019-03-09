package main

import (
	"fmt"
	"strings"

	"github.com/diamondburned/discordgo"
)

const (
	formatReactionConstant = `[:%s:] %d %s [:-:]  `

	unreactColor = "#383838"
	reactColor   = "#2196F3"
)

func removeAllReactions(rs []*discordgo.MessageReactions, i int) []*discordgo.MessageReactions {
	if i > 0 {
		return append(rs[:i], rs[i+1:]...)
	}

	return rs[i+1:]
}

func isSameEmoji(rs *discordgo.MessageReactions, r *discordgo.MessageReaction) bool {
	if rs.Emoji.ID != 0 || r.Emoji.ID != 0 {
		return rs.Emoji.ID == r.Emoji.ID
	}

	return rs.Emoji.Name == r.Emoji.Name
}

func handleReactionEvent(m *discordgo.Message) {
	if rstore.Check(m.Author, RelationshipBlocked) && cfg.Prop.HideBlocked {
		return
	}

	for i, msg := range messageStore {
		if strings.HasPrefix(msg, fmt.Sprintf("\n"+`["%d"]`, m.ID)) {
			msg := fmt.Sprintf(
				messageFormat+"[::-]",
				m.ID, fmtMessage(m),
			)

			messageStore[i] = msg

			break
		}
	}

	app.QueueUpdateDraw(func() {
		messagesView.SetText(strings.Join(messageStore, ""))
	})

	scrollChat()
}

func formatReactions(rs []*discordgo.MessageReactions) (f string, eM map[string][]string) {
	eM = make(map[string][]string)

	for _, r := range rs {
		f += formatReactionString(r)

		if r.Emoji.ID == 0 {
			continue
		}

		var format = "png"
		if r.Emoji.Animated {
			format = "gif"
		}

		IDstring := fmt.Sprintf("%d", r.Emoji.ID)

		eM[IDstring] = []string{
			r.Emoji.Name,
			`https://cdn.discordapp.com/emojis/` + IDstring + `.` + format,
		}
	}

	return
}

func formatReactionString(r *discordgo.MessageReactions) string {
	if r.Emoji == nil {
		return ""
	}

	var color = unreactColor
	if r.Me {
		color = reactColor
	}

	return fmt.Sprintf(
		formatReactionConstant,
		color, r.Count, strings.TrimSpace(r.Emoji.Name),
	)

}
