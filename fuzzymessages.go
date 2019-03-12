package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/sahilm/fuzzy"
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
		matches := fuzzy.Find(
			strings.TrimPrefix(text, "~"),
			allMessages,
		)

		for _, m := range matches {
			fuzzied = append(
				fuzzied,
				m.Str,
			)
		}

	} else {
		fuzzied = allMessages
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, u := range fuzzied {
			autocomp.InsertItem(
				i, u,
				"", 0, nil,
			)
		}

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
		messagesView.Highlight(
			strings.Split(t, " - ")[0],
		)

		messagesView.ScrollToHighlight()
	})
}
