package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/tview/v2"
)

// [0]:format [1]:ID
var allMessages []*tview.ListItem

func fuzzyMessages(text string) {
	var fuzzied []*tview.ListItem

	if len(allMessages) == 0 && Channel != nil {
		for i := len(messageStore) - 1; i >= 0; i-- {
			ID := getIDfromindex(i)
			if ID == 0 {
				continue
			}

			allMessages = append(allMessages, makeMessageItem(ID))
		}
	}

	if len(text) > 1 {
		text := strings.TrimPrefix(text, "~")
		fuzzied = make([]*tview.ListItem, 0, len(allMessages))

		for _, m := range allMessages {
			if strings.Contains(m.MainText, text) {
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

		autocomp.SetItems(fuzzied)
		autocomp.SetCurrentItem(-1)

		rightflex.ResizeItem(autocomp, min(len(fuzzied), 10), 1)

		autofillfunc = func(i int) {
			words := strings.Fields(input.GetText())

			withoutlast := words[:len(words)-1]
			withoutlast = append(withoutlast, fuzzied[i].SecondaryText)

			input.SetText(strings.Join(withoutlast, " ") + " ")

			clearList()
			app.SetFocus(input)
		}

	} else {
		rightflex.ResizeItem(autocomp, 1, 1)
	}

	app.Draw()

	autocomp.SetChangedFunc(func(i int, t string, st string, s rune) {
		if i >= len(fuzzied) {
			return
		}

		ID := fuzzied[i].SecondaryText

		if Channel != nil {
			id, _ := strconv.ParseInt(ID, 10, 64)
			if id != 0 {
				m, err := d.State.Message(Channel.ID, id)
				if err == nil {
					// Update the list entry
					item := makeMessageItem(id)
					if item.MainText != t {
						fuzzied[i].MainText = item.MainText
					}

					imageRendererPipeline.add(m)
				}
			}
		}

		messagesView.Highlight(ID)
		messagesView.ScrollToHighlight()
	})
}

func makeMessageItem(ID int64) *tview.ListItem {
	id := strconv.FormatInt(ID, 10)

	m, err := d.State.Message(Channel.ID, ID)
	if err != nil {
		return &tview.ListItem{
			MainText:      id + " - ???",
			SecondaryText: id,
			Shortcut:      0,
			Selected:      nil,
		}
	}

	username, color := us.DiscordThis(m)

	sentTime, err := m.Timestamp.Parse()
	if err != nil {
		sentTime = time.Now()
	}

	var fetchedColor = readChannelColorPrefix
	if s := imageRendererPipeline.cache.get(m.ID); s != nil {
		fetchedColor = string(s.state)
	}

	return &tview.ListItem{
		MainText: fmt.Sprintf(
			"%s%s[-] - [#%06X]%s[-] [::d]- %s[::-]",
			fetchedColor, id, color, username,
			sentTime.Local().Format(time.Stamp),
		),
		SecondaryText: id,
		Shortcut:      0,
		Selected:      nil,
	}
}
