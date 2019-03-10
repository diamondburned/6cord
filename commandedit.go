package main

import (
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/diamondburned/tcell"
	"github.com/diamondburned/discordgo"
)

const (
	// EditMessageLabel is used to detect which function to use
	EditMessageLabel = "Edit "
)

var toEditMessage int64

func editMessage(text []string) {
	if Channel == nil {
		Message("You're not in a channel!")
		return
	}

	var messageN int

	if len(text) == 2 {
		i, err := strconv.Atoi(text[1])
		if err != nil {
			Message(err.Error())
			return
		}

		messageN = i
	}

	lastMsg := matchMyMessage(messageN)
	if lastMsg == nil {
		Message("Can't find your last message :(")
		return
	}

	toEditMessage = lastMsg.ID

	input.SetBackgroundColor(tcell.ColorBlue)
	input.SetFieldBackgroundColor(tcell.ColorBlue)
	input.SetPlaceholderTextColor(tcell.ColorWhite)
	input.SetLabel("Edit ")
	input.SetPlaceholder("Send this empty message to delete it")
	input.SetText(lastMsg.Content)
}

func editHandler() {
	var (
		i   = input.GetText()
		err error
	)

	if i != "" {
		_, err = d.ChannelMessageEdit(
			Channel.ID, toEditMessage, i,
		)
	} else {
		err = d.ChannelMessageDelete(
			Channel.ID, toEditMessage,
		)
	}

	toEditMessage = 0

	if err != nil {
		Warn(err.Error())
	}

	resetInputBehavior()
}

func editMessageRegex(text string) {
	if Channel == nil {
		Message("You're not in a channel!")
	}

	input := csv.NewReader(strings.NewReader(text))
	input.Comma = '/' // delimiter
	args, err := input.Read()
	if err != nil {
		Warn(err.Error())
		return
	}

	if len(args) != 3 && len(args) != 4 {
		Message(fmt.Sprintf("Invalid arguments! %d", len(args)))
		return
	}

	var (
		regexArg = args[1]
		withArg  = args[2]
		messageN int
	)

	if len(args) == 4 {
		order := args[3]

		if order != "" && order != "g" {
			messageN, _ = strconv.Atoi(order)
		}
	}

	regex, err := regexp.Compile(regexArg)
	if err != nil {
		Message(err.Error())
		return
	}

	lastMsg := matchMyMessage(messageN)
	if lastMsg == nil {
		Message("Can't find your last message :(")
		return
	}

	repl := regex.ReplaceAllString(lastMsg.Content, withArg)

	_, err = d.ChannelMessageEdit(
		lastMsg.ChannelID,
		lastMsg.ID,
		repl,
	)

	if err != nil {
		Warn(err.Error())
	}
}

func matchMyMessage(residue int) *discordgo.Message {
	m, err := d.State.Message(Channel.ID, int64(residue))
	if err == nil && m.Author.ID == d.State.User.ID {
		return m
	}

	for i := len(messageStore) - 1; i >= 0; i-- {
		if ID := getIDfromindex(i); ID != 0 {
			m, err := d.State.Message(Channel.ID, ID)
			if err != nil {
				continue
			}

			if m.Author.ID == d.State.User.ID {
				if residue == 0 {
					return m
				}

				residue--
			}
		}
	}

	return nil
}
