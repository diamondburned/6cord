package main

import (
	"fmt"
	"strings"
	"time"
)

var allMessages []string

func fuzzyMessages(text string) {
	var fuzzied []string

	if len(allMessages) == 0 && Channel != nil {
		for i := len(messageStore) - 1; i >= 0; i-- {
			if ID := getIDfromindex(i); ID != 0 {
				m, err := d.State.Message(Channel.ID, ID)
				if err != nil {
					continue
				}

				username, color := us.DiscordThis(m)

				sentTime, err := m.Timestamp.Parse()
				if err != nil {
					sentTime = time.Now()
				}

				allMessages = append(
					allMessages,
					fmt.Sprintf(
						"%d - [#%06X]%s[-] [::d]- %s[::-]",
						ID, color, username,
						sentTime.Local().Format(time.Stamp),
					),
				)
			}
		}
	}

	if len(text) > 1 {
		text := strings.TrimPrefix(text, "~")
		fuzzied = make([]string, 0, len(allMessages))

		for _, m := range allMessages {
			if strings.Contains(m, text) {
				fuzzied = append(fuzzied, m)
			}
		}
	} else {
		fuzzied = allMessages
	}

	clearList()

	if len(fuzzied) > 0 {
		for _, u := range fuzzied {
			autocomp.InsertItem(
				0, u,
				"", 0, nil,
			)
		}

		autocomp.SetCurrentItem(-1)

		rightflex.ResizeItem(autocomp, min(len(fuzzied), 10), 1)

		autofillfunc = func(i int) {
			words := strings.Fields(input.GetText())

			withoutlast := words[:len(words)-1]
			withoutlast = append(
				withoutlast,
				strings.Split(fuzzied[i], " - ")[0],
			)

			input.SetText(strings.Join(withoutlast, " ") + " ")

			clearList()
			app.SetFocus(input)
		}

	} else {
		rightflex.ResizeItem(autocomp, 1, 1)
	}

	app.Draw()

	autocomp.SetChangedFunc(func(i int, t string, st string, s rune) {
		ID := strings.Split(t, " - ")[0]

		checkForImage(ID)

		messagesView.Highlight(ID)
		messagesView.ScrollToHighlight()
	})
}
