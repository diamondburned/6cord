package main

import (
	"fmt"
	"strconv"
)

func deleteMessage(text []string) {
	toDelete := make([]int64, 0, len(text))

	if len(text) == 1 {
		lastMsg := matchMyMessage(0)
		if lastMsg == nil {
			Message("Can't find your last message :(")
			return
		}

		toDelete = append(toDelete, lastMsg.ID)
	} else {
		for i, a := range text[1:] {
			m, err := strconv.Atoi(a)
			if err != nil {
				Message(fmt.Sprintf("Failed to parse argument %d", i-1))
				return
			}

			lastMsg := matchMyMessage(m)
			if lastMsg == nil {
				Message("Can't find your last message :(")
				return
			}

			toDelete = append(toDelete, lastMsg.ID)
		}
	}

	for _, m := range toDelete {
		if err := d.ChannelMessageDelete(Channel.ID, m); err != nil {
			Warn(err.Error())
			return
		}
	}
}
