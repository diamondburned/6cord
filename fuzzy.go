package main

import (
	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tcell"
	"github.com/diamondburned/tview/v2"
	"github.com/sahilm/fuzzy"
)

func stateResetter() {
	channelFuzzyCache = allChannels([]fuzzyReadState{})
	discordEmojis = DiscordEmojis([]*discordgo.Emoji{})
	allMessages = make([]*tview.ListItem, 0, len(messageStore))
	autocomp.SetChangedFunc(nil)
	messagesView.Highlight()

	current, total := getLineStatus()

	// If the scroll offset is < 20
	if total-current < 20 {
		scrollChat()
	}

	if imageRendererPipeline != nil {
		imageRendererPipeline.clean()
	}
}

func clearList() {
	rightflex.ResizeItem(autocomp, 1, 1)

	if autocomp.GetItemCount() != 0 {
		autocomp.Clear()
	}
}

func formatNeedle(m fuzzy.Match) (f string) {
	isHL := false

	for i := 0; i < len(m.Str); i++ {
		if fuzzyHasNeedle(i, m.MatchedIndexes) {
			f += "[::u]" + string(m.Str[i])
			isHL = true
		} else {
			if isHL {
				f += "[::-]"
			}

			f += string(m.Str[i])
		}
	}

	return
}

func fuzzyHasNeedle(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}

func autocompHandler(ev *tcell.EventKey) *tcell.EventKey {
	i := autocomp.GetCurrentItem()

	switch ev.Key() {
	case tcell.KeyDown:
		if i+1 == autocomp.GetItemCount() {
			app.SetFocus(input)
			return nil
		}

		return ev

	case tcell.KeyUp:
		if i == 0 {
			app.SetFocus(input)
			return nil
		}

		return ev

	case tcell.KeyLeft:
		imageRendererPipeline.prev()
		return nil

	case tcell.KeyRight:
		imageRendererPipeline.next()
		return nil

	case tcell.KeyEnter:
		return ev
	}

	if ev.Rune() >= 0x31 && ev.Rune() <= 0x122 {
		return ev
	}

	app.SetFocus(input)
	return nil
}

var autofillfunc func(i int)
var onhoverfn func(i int)
