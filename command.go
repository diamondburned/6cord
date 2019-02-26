package main

import (
	"os"
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
		Command:     "/nick",
		Function:    changeSelfNick,
		Description: "[nickname] - changes nickname for the current guild",
	},
	Command{
		Command:     "/status",
		Function:    setStatus,
		Description: "[online|busy|away|invisible] - sets your status",
	},
	Command{
		Command:     "/edit",
		Function:    editMessage,
		Description: "[n:int optional] - edits the latest n message",
	},
	Command{
		Command:     "/presence",
		Function:    setGame,
		Description: "[string] - sets your \"Playing\" or \"Listening to\" presence, empty to reset",
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
		Command:     "/quit",
		Function:    commandExit,
		Description: "quits",
	},
}

func commandExit(text []string) {
	app.Stop()
	os.Exit(0)
}

// CommandHandler .
func CommandHandler() {
	text := input.GetText()
	if text == "" {
		return
	}

	defer input.SetText("")

	switch {
	case strings.HasPrefix(text, "s/"):
		editMessageRegex(text)

	case strings.HasPrefix(text, "/"):
		f := strings.Fields(text)
		if len(f) < 0 {
			return
		}

		for _, cmd := range commands {
			if f[0] == cmd.Command && cmd.Function != nil {
				go cmd.Function(f)
				return
			}
		}

		fallthrough
	default:
		// Trim literal backslash, in case "\/actual message"
		text = strings.TrimPrefix(text, `\`)
		text = senderRegex.Replace(text)

		go func(text string) {
			if _, err := d.ChannelMessageSend(ChannelID, text); err != nil {
				Warn("Failed to send message:\n" + text + "\nError: " + err.Error())
			}
		}(text)
	}
}
