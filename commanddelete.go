package main

import "strconv"

func deleteMessage(text []string) {
	toDelete := make([]int64, 0, len(text)-1)

	if len(text) < 2 {
		lastMsg := matchMyMessage(0)
		if lastMsg == nil {
			Message("Can't find your last message :(")
			return
		}

		toDelete = append(toDelete, lastMsg.ID)
	} else {
		for _, a := range text[1:] {
			m, err := strconv.ParseInt(a, 10, 64)
			if err != nil {
				Message("Failed to find the message.")
				return
			}

			toDelete = append(toDelete, m)
		}
	}

	for _, m := range toDelete {
		if err := d.ChannelMessageDelete(Channel.ID, m); err != nil {
			Warn(err.Error())
			return
		}
	}
}
