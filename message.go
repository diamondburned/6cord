package main

import (
	"strings"
)

// hopefully avoids messagesView going out of range after a
// while
func cleanupBuffer() {
	if len(messageStore) > prefetchMessageCount*6 {
		messageStore = messageStore[prefetchMessageCount*4:]

		messagesView.SetText(
			strings.Join(messageStore, ""),
		)
	}
}

func scrollChat() (clear bool) {
	// If the message box is not focused and the input is empty
	if !messagesView.HasFocus() && input.GetText() == "" {
		clear = true
		cleanupBuffer()
	}

	current, lines := getLineStatus()
	if lines-current > 5 {
		return false
	}

	if clear {
		messagesView.ScrollToEnd()
	}

	if Channel == nil {
		wrapFrame.SetTitle(generateTitle(Channel))
	}

	return
}
