package main

import (
	"strconv"
	"strings"
)

func parseUserMention(m string) int64 {
	i := 0
	trim := strings.TrimFunc(m, func(r rune) bool {
		switch {
		case r == '<', r == '>', r == '@':
			i++
			return true
		default:
			return false
		}
	})

	id, _ := strconv.ParseInt(trim, 10, 64)
	if i == 0 {
		return 0
	}

	return id
}

func makeDirectMessage(text []string) {
	if len(text) != 2 {
		Message("No channels given!")
		return
	}

	id := parseUserMention(text[1])
	if id == 0 {
		Message("Invalid user mention!")
		return
	}

	ch, err := d.UserChannelCreate(id)
	if err != nil {
		Message(err.Error())
	}

	loadChannel(ch.ID)
}
