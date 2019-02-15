package main

import (
	"encoding/csv"
	"log"
	"strings"

	"github.com/RumbleFrog/discordgo"
)

var (
	senderRegex = strings.NewReplacer(
		`\n`, "\n",
		`\t`, "\t",
	)
)

// Command contains a command's info
type Command struct {
	Function    func([]string)
	Description string
}

// Commands ..
var Commands = map[string]Command{
	"/status": Command{
		Function:    setStatus,
		Description: "[online|busy|away|invisible] - sets your status",
	},
}

func setStatus(input []string) {
	if len(input) < 2 {
		if d.State.Settings == nil {
			Message("Settings are uninitialized")
			return
		}

		switch s := d.State.Settings.Status; s {
		case discordgo.StatusOnline:
			Message("Status: Online")
		case discordgo.StatusIdle:
			Message("Status: Idle")
		case discordgo.StatusDoNotDisturb:
			Message("Status: Do not disturb")
		case discordgo.StatusInvisible:
			Message("Status: Invisible")
		default:
			Message(string(s))
		}

		return
	}

	switch input[1] {
	}
}

// CommandHandler .
func CommandHandler() {
	text := input.GetText()
	if text == "" {
		return
	}

	defer input.SetText("")

	switch {
	case strings.HasPrefix(text, "/"):
		f := strings.Fields(text)
		if len(f) < 0 {
			return
		}

		for cmd, command := range Commands {
			if f[0] == cmd {
				command.Function(f)
				return
			}
		}

	case strings.HasPrefix(text, "s/"):
		// var (
		// 	ReplaceRegex string
		// 	ReplaceWith  string
		// 	MessageOrder = 1
		// )

		input := csv.NewReader(strings.NewReader(text))
		input.Comma = '/' // delimiter
		args, err := input.Read()
		if err != nil {
			log.Println(err)
			return
		}

		if len(args) < 3 {
			log.Println("")
		}

	default:
		text = senderRegex.Replace(text)

		go func(text string) {
			if _, err := d.ChannelMessageSend(ChannelID, text); err != nil {
				Warn("Failed to send message:\n" + text + "\nError: " + err.Error())
			}
		}(text)
	}
}
