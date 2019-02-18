package main

import "strconv"

func parseChannelID(text []string) int64 {
	chID := text[1]

	chID = chID[2:]
	chID = chID[:len(chID)-1]

	id, err := strconv.ParseInt(chID, 10, 64)
	if err != nil {
		Message(err.Error())
		return 0
	}

	return id
}

func gotoChannel(text []string) {
	id := parseChannelID(text)
	if id == 0 {
		return
	}

	ChannelID = id
	loadChannel()
}
