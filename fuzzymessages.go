package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sahilm/fuzzy"
)

var allMessages []string

var lastImgCtx *imageCtx

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
		ID := strings.Split(t, " - ")[0]

		go func() {
			if lastImgCtx != nil {
				lastImgCtx.Delete()
			}

			if Channel == nil {
				return
			}

			id, _ := strconv.ParseInt(ID, 10, 64)
			if id == 0 {
				return
			}

			m, err := d.State.Message(Channel.ID, id)
			if err != nil {
				return
			}

			if len(m.Attachments) == 0 {
				return
			}

			lastImgCtx = newDiscordImageContext(
				m.Attachments[0].ProxyURL,
				m.Attachments[0].Width,
				m.Attachments[0].Height,
			)
		}()

		messagesView.Highlight(ID)
		messagesView.ScrollToHighlight()
	})
}
