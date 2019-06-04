package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// [0]:format [1]:ID
var allMessages [][2]string

func fuzzyMessages(text string) {
	var fuzzied [][2]string

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

				var fetchedColor = readChannelColorPrefix
				if imageRendererPipeline.cache.get(m.ID) != nil {
					fetchedColor = unreadChannelColorPrefix
				}

				id := strconv.FormatInt(ID, 10)

				allMessages = append(allMessages, [2]string{
					fmt.Sprintf(
						"%s%s[-::] - [#%06X]%s[-] [::d]- %s[::-]",
						fetchedColor, id, color, username,
						sentTime.Local().Format(time.Stamp),
					), id,
				})
			}
		}
	}

	if len(text) > 1 {
		text := strings.TrimPrefix(text, "~")
		fuzzied = make([][2]string, 0, len(allMessages))

		for _, m := range allMessages {
			if strings.Contains(m[0], text) {
				fuzzied = append(fuzzied)
			}
		}
	} else {
		fuzzied = allMessages
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, j := 0, len(fuzzied)-1; i < j; i, j = i+1, j-1 {
			fuzzied[i], fuzzied[j] = fuzzied[j], fuzzied[i]
		}

		for i, u := range fuzzied {
			autocomp.InsertItem(
				i, u[0],
				"", 0, nil,
			)
		}

		autocomp.SetCurrentItem(-1)

		rightflex.ResizeItem(autocomp, min(len(fuzzied), 10), 1)

		autofillfunc = func(i int) {
			words := strings.Fields(input.GetText())

			withoutlast := words[:len(words)-1]
			withoutlast = append(withoutlast, fuzzied[i][1])

			input.SetText(strings.Join(withoutlast, " ") + " ")

			clearList()
			app.SetFocus(input)
		}

	} else {
		rightflex.ResizeItem(autocomp, 1, 1)
	}

	app.Draw()

	autocomp.SetChangedFunc(func(i int, t string, st string, s rune) {
		ID := fuzzied[i][1]

		if Channel != nil {
			id, _ := strconv.ParseInt(ID, 10, 64)
			if id != 0 {
				m, err := d.State.Message(Channel.ID, id)
				if err == nil {
					imageRendererPipeline.add(m)
				}
			}
		}

		messagesView.Highlight(ID)
		messagesView.ScrollToHighlight()
	})
}
