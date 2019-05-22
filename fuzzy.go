package main

import (
	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tcell"
	"github.com/sahilm/fuzzy"
)

func stateResetter() {
	channelFuzzyCache = allChannels([]fuzzyReadState{})
	discordEmojis = DiscordEmojis([]*discordgo.Emoji{})
	allMessages = []string{}
	autocomp.SetChangedFunc(nil)
	messagesView.Highlight()

	imageRendererPipeline.clean()
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

func min(i, j int) int {
	if i < j {
		return i
	}

	return j
}

func autocompHandler(ev *tcell.EventKey) *tcell.EventKey {
	switch ev.Key() {
	case tcell.KeyDown:
		if autocomp.GetCurrentItem()+1 == autocomp.GetItemCount() {
			app.SetFocus(input)
			return nil
		}

		return ev

	case tcell.KeyUp:
		if autocomp.GetCurrentItem() == 0 {
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

// fuck you you fucking tview dev
// you could just fucking globalize the goddamn selection
// function, but no. you didn't. why the fuck didn't you?
// are you fucking retarded in the head? stop trying to
// handle everything by you yourself and your shitty shoddy
// little functions, you fucking stupid asshat
// bloody fucking jesus i fucking hate doing this, but
// this is literally my only fucking choice
var autofillfunc func(i int)
var onhoverfn func(i int)
