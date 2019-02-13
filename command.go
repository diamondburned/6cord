package main

import (
	"encoding/csv"
	"log"
	"strings"
)

var (
	senderRegex = strings.NewReplacer(
		`\n`, "\n",
		`\t`, "\t",
	)
)

// Command contains a command's info
type Command struct {
	Function    func(string)
	Description string
}

// Commands ..
var Commands = map[string]Command{
	"/status": Command{
		Function:    setStatus,
		Description: "[online|busy|away|invisible] - sets your status",
	},
}

func setStatus(input string) {

}

func CommandHandler() {
	text := input.GetText()
	if text == "" {
		return
	}

	switch {
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

	input.SetText("")
}
