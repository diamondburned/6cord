package main

import (
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/rumblefrog/discordgo"
)

func matchCopyMessage(text []string) {
	if len(text) != 2 {
		Message("Invalid args! Refer to description.")
		return
	}

	residue, err := strconv.Atoi(text[1])
	if err != nil {
		Message(err.Error())
		return
	}

	var message *discordgo.Message

	m, err := d.State.Message(ChannelID, int64(residue))
	if err == nil {
		message = m
	} else {
		for i := len(messageStore) - 1; i >= 0; i-- {
			if ID := getIDfromindex(i); ID != 0 {
				m, err := d.State.Message(ChannelID, ID)
				if err != nil {
					continue
				}

				if residue == 0 {
					message = m
					break
				}

				residue--
			}
		}
	}

	if message == nil {
		Message("Can't find any message to copy.")
		return
	}

	if err := clipboard.WriteAll(message.Content); err != nil {
		Warn(err.Error())
	} else {
		Message("Copied message from " + message.Author.Username)
	}

}
