package main

import (
	"strings"

	"github.com/RumbleFrog/discordgo"
	"github.com/rivo/tview"
	"gitlab.com/diamondburned/6cord/md"
)

func fmtMessage(m *discordgo.Message) string {
	var (
		ct = md.Parse(
			m.ContentWithMentionsReplaced(),
		)

		edited string
		c      []string
		l      = strings.Split(ct, "\n")
	)

	for i := 0; i < len(l); i++ {
		c = append(c, "\t"+l[i])
	}

	if len(m.Attachments) > 0 {
		for _, a := range m.Attachments {
			c = append(c, "\t"+tview.Escape(a.URL))
		}
	}

	if m.EditedTimestamp != "" {
		edited = " [::d](edited)[::-]"
	}

	return strings.Join(c, "\n") + edited
}
