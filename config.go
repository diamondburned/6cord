package main

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
// For shell-like argument splitting, join the array and run it through
// a CSV reader, delimiter ' '.
func cmdShrug(text []string) {
	if Channel == nil {
		// Error handling in case nil crashes the entire app
		Message("You're not in a channel!")
	}

	// Channel is a global variable indicating the current channel.
	// Writing to this variable will screw _everthing_ up.
	if _, err := d.ChannelMessageSend(Channel.ID, `¯\_(ツ)_/¯`); err != nil {
		Warn(err.Error())
	}
}
