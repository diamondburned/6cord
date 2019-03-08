package main

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

func commandMentions(text []string) {
	input.SetPlaceholder("Loading mentions...")
	defer input.SetPlaceholder(DefaultStatus)

	mentions, err := getMentions()
	if err != nil {
		Warn(err.Error())
		return
	}

	Channel = nil
	messageStore = []string{}
	messagesView.Clear()

	for i := len(mentions) - 1; i >= 0; i-- {
		m := mentions[i]

		username, color := us.DiscordThis(m)

		sentTime, err := m.Timestamp.Parse()
		if err != nil {
			sentTime = time.Now()
		}

		messagesView.Write([]byte(
			fmt.Sprintf(
				authorFormat,
				color, tview.Escape(username),
				sentTime.Format(time.Stamp),
			),
		))

		messagesView.Write([]byte(
			fmt.Sprintf(
				messageFormat+"[::-]",
				m.ID, fmtMessage(m),
			),
		))
	}

	wrapFrame.SetTitle("[Mentions[]")
	input.SetPlaceholder("Done.")
	app.Draw()

	time.Sleep(time.Second * 5)
}
