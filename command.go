package main

import (
	"fmt"
	"strings"
)

var (
	senderRegex = strings.NewReplacer()
	cmdHistory  = make([]string, 0, 256)
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
		Command:     "/editor",
		Function:    commandEditor,
		Description: "Pop up $EDITOR to send a message <C-e>",
	},
	Command{
		Command:     "/mentions",
		Function:    commandMentions,
		Description: "shows the last mentions",
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
		Description: "[n:int optional] - edits the latest n message of yours",
	},
	Command{
		Command:     "/delete",
		Function:    deleteMessage,
		Description: "[messageID:int] - deletes the message",
	},
	Command{
		Command:     "/presence",
		Function:    setGame,
		Description: "[string] - sets your \"Playing\" or \"Listening to\" presence, empty to reset",
	},
	Command{
		Command:     "/react",
		Function:    reactMessage,
		Description: "[messageID:int] [emoji:string] - toggle reaction on a message",
	},
	Command{
		Command:     "/upload",
		Function:    uploadFile,
		Description: "[file path] - uploads file",
	},
	Command{
		Command:     "/heated",
		Function:    cmdHeated,
		Description: "warns you when a message is sent, regardless of settings",
	},
	Command{
		Command:     "/copy",
		Function:    matchCopyMessage,
		Description: "[n:int] - copies the entire last n message",
	},
	Command{
		Command:     "/highlight",
		Function:    highlightMessage,
		Description: "[ID:int64] - highlights the message ID if possible",
	},
	Command{
		Command:     "/dm",
		Function:    makeDirectMessage,
		Description: "[@mention] - starts a new direct message",
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
		Command:     "/debug",
		Function:    commandDebug,
		Description: "prints extra debug info",
	},
	Command{
		Command:     "/quit",
		Function:    commandExit,
		Description: "quits <C-c>",
	},
}

func commandExit(text []string) {
	app.Stop()
}

// CommandHandler .
func CommandHandler() {
	text := input.GetText()
	if text == "" {
		return
	}

	defer input.SetText("")

	if len(cmdHistory) >= 256 {
		cmdHistory = cmdHistory[255:]
	}

	cmdHistory = append(
		cmdHistory, text,
	)

	switch {
	case strings.HasPrefix(text, "s/"):
		go editMessageRegex(text)

	case strings.HasPrefix(text, "/"):
		f := strings.Fields(text)
		if len(f) < 0 {
			return
		}

		for _, cmd := range commands {
			if f[0] == cmd.Command && cmd.Function != nil {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							Warn(fmt.Sprintf("%v", r))
						}
					}()

					cmd.Function(f)
				}()

				return
			}
		}

		fallthrough
	default:
		// Trim literal backslash, in case "\/actual message"
		text = strings.TrimPrefix(text, `\`)

		if Channel == nil {
			Message("You're not in a channel!")
			return
		}

		go func(text string) {
			_, err := d.ChannelMessageSend(Channel.ID, processString(text))
			if err != nil {
				Warn("Failed to send message:\n" + text + "\nError: " + err.Error())
			}

			messagesView.ScrollToEnd()
			app.Draw()
		}(text)
	}
}
