package main

import (
	"strconv"
	"time"

	"github.com/diamondburned/tview"
)

func commandMentions(text []string) {
	input.SetPlaceholder("Loading mentions...")
	defer input.SetPlaceholder(cfg.Prop.DefaultStatus)

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
			authorTmpl.ExecuteString(map[string]interface{}{
				"color": fmtHex(color),
				"name":  tview.Escape(username),
				"time":  sentTime.Format(time.Stamp),
			}),
		))

		messagesView.Write([]byte(
			messageTmpl.ExecuteString(map[string]interface{}{
				"ID":      strconv.FormatInt(m.ID, 10),
				"content": fmtMessage(m),
			}),
		))
	}

	wrapFrame.SetTitle("[Mentions[]")
	input.SetPlaceholder("Done.")
	app.Draw()

	time.Sleep(time.Second * 5)
}
