package main

import "strconv"

func deleteMessage(text []string) {
	if len(text) < 2 {
		Message("Missing arguments! Use Shift + '`' to select messages.")
		return
	}

	for _, a := range text[1:] {
		m, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			Message("Failed to find the message.")
			return
		}

		if err := d.ChannelMessageDelete(Channel.ID, m); err != nil {
			Warn(err.Error())
			return
		}
	}
}
