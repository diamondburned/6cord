package main

import (
	"strings"
)

// hopefully avoids messagesView going out of range after a
// while
func cleanupBuffer() {
	if len(messageStore) > 512 && !messagesView.HasFocus() {
		messageStore = messageStore[:256]

		app.QueueUpdateDraw(func() {
			messagesView.SetText(
				strings.Join(messageStore, ""),
			)
		})

		messagesView.ScrollToEnd()
	}
}

func scrollChat() bool {
	current, lines := getLineStatus()
	if lines-current > 5 {
		return false
	}

	if Channel == nil {
		wrapFrame.SetTitle(generateTitle(Channel))
	}

	if !messagesView.HasFocus() && !autocomp.HasFocus() {
		messagesView.ScrollToEnd()
		cleanupBuffer()
	}

	return true
}
