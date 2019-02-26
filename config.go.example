package main

const (
	// ChatPadding determines the width the message is from the border
	ChatPadding = 2

	// HideBlocked when true, ignores all blocked users
	HideBlocked = true

	// BackgroundColor self explanatory
	// Acceptable values: tcell.Color*, -1, 0-255 (terminal colors)
	// Check https://jonasjacek.github.io/colors/ for reference
	BackgroundColor = -1

	// CommandPrefix is prefix, like $PS1
	CommandPrefix = "[white]> "

	// SyntaxHighlightColorscheme is the color scheme for syntax highlighting
	// https://xyproto.github.io/splash/docs/all.html
	SyntaxHighlightColorscheme = "vs"

	// ShowEmojiURLs converts emojis into clickable URLs if true
	ShowEmojiURLs = true
)

// CustomCommands is for user-made commands
var CustomCommands = []Command{
	Command{
		Command:     "/shrug",
		Function:    cmdShrug,
		Description: `¯\_(ツ)_/¯`,
	},
}

// `text` is the chat argument, split into arrays.
// This is done with strings.Fields(messageContent).
// For shell-like argument splitting, join the array
// and run it through a CSV reader, delimiter ' '.
func cmdShrug(text []string) {
	if _, err := d.ChannelMessageSend(ChannelID, `¯\_(ツ)_/¯`); err != nil {
		Warn(err.Error())
	}
}
