package main

func cmdHeated(text []string) {
	if Channel == nil {
		Message("You're not in a channel!")
		return
	}

	if heatedChannelsToggle(Channel.ID) {
		Message("Added this channel. We'll warn you when there's a message.")
	} else {
		Message("Removed this channel.")
	}

	// Heated servers are checked in notify.go
	// Check function is in heated.go
}
