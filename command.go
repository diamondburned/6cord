package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rumblefrog/discordgo"
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
		// var (
		// 	ReplaceRegex string
		// 	ReplaceWith  string
		// 	MessageOrder = 1
		// )

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

			lastMsg *discordgo.Message
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

		for i := len(messageStore) - 1; i >= 0; i-- {
			if ID := getIDfromindex(i); ID != 0 {
				m, err := d.State.Message(ChannelID, ID)
				if err != nil {
					continue
				}

				if m.Author.ID == d.State.User.ID {
					if messageN == 0 {
						lastMsg = m
						break
					}

					messageN--
				}
			}
		}

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
