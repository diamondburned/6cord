package main

import (
	"strconv"
	"strings"
)

func makeDirectMessage(text []string) {
	if len(text) != 2 {
		Message("No channels given!")
		return
	}

	i := 0
	trim := strings.TrimFunc(text[1], func(r rune) bool {
		switch {
		case r == '<', r == '>', r == '@':
			i++
			return true
		default:
			return false
		}
	})

	id, err := strconv.ParseInt(trim, 10, 64)
	if i != 3 || err != nil {
		Message("Invalid user mention!")
		return
	}

	ch, err := d.UserChannelCreate(id)
	if err != nil {
		Message(err.Error())
	}

	loadChannel(ch.ID)
}
