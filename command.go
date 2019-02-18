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

// Commands contains multiple commands
type Commands []Command

// Command contains a command's info
type Command struct {
	Command     string
	Function    func([]string)
	Description string
}

var commands = Commands{
	Command{
		Command:     "/goto",
		Function:    gotoChannel,
		Description: "[channel name] - jumps to a channel",
	},
	Command{
		Command:     "/status",
		Function:    setStatus,
		Description: "[online|busy|away|invisible] - sets your status",
	},
	Command{
		Command:     "/upload",
		Function:    uploadFile,
		Description: "[file path] - uploads file",
	},
	Command{
		Command:     "/block",
		Function:    blockUser,
		Description: "[@mention] - blocks someone",
	},
	Command{
		Command:     "/unblock",
		Function:    unblockUser,
		Description: "[@mention] - unblocks someone",
	},
	Command{
		Command:     "/nick",
		Function:    changeSelfNick,
		Description: "[nickname] - changes nickname for the current guild",
	},
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

		for _, cmd := range commands {
			if f[0] == cmd.Command && cmd.Function != nil {
				cmd.Function(f)
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
