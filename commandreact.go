package main

import (
	"strconv"
	"strings"
)

func reactMessage(text []string) {
	if Channel == nil {
		Message("You're not in a channel!")
		return
	}

	if len(text) < 3 {
		Message("Invalid arguments! Refer to description.")
		return
	}

	messageID, err := strconv.ParseInt(text[1], 10, 64)
	if err != nil {
		Message("Failed to find the message.")
		return
	}

	message, err := d.State.Message(Channel.ID, messageID)
	if err != nil {
		Message("Failed to find the message.")
		return
	}

	var (
		emoji   = make([]string, len(text)-2)
		reacted = make([]bool, len(text)-2)
	)

	for i := 0; i < len(text)-2; i++ {
		regres := EmojiRegex.FindAllStringSubmatch(text[i+2], -1)
		if len(regres) > 0 && len(regres[0]) == 4 {
			emoji[i] = regres[0][2] + ":" + regres[0][3]

			for _, r := range message.Reactions {
				if r.Emoji == nil {
					continue
				}

				if strconv.FormatInt(r.Emoji.ID, 10) == regres[0][3] {
					reacted[i] = r.Me
					break
				}
			}
		} else {
			emoji[i] = strings.TrimSpace(text[i+2])

			for _, r := range message.Reactions {
				if r.Emoji == nil {
					continue
				}

				if r.Emoji.Name == text[i+2] {
					reacted[i] = r.Me
					break
				}
			}
		}
	}

	for i := 0; i < len(emoji); i++ {
		if reacted[i] {
			err = d.MessageReactionRemoveMe(
				Channel.ID,
				message.ID,
				emoji[i],
			)
		} else {
			err = d.MessageReactionAdd(
				Channel.ID,
				message.ID,
				emoji[i],
			)
		}

		if err != nil {
			Warn(err.Error())
			return
		}
	}

}
