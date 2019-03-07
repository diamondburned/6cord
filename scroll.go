package main

import (
	"fmt"
	"strings"
)

func handleScroll() {
	lines := len(
		strings.Split(
			strings.Join(messageStore, ""),
			"\n",
		),
	)

	if lines <= 1 {
		lines = len(messagesView.GetText(false))
	}

	var (
		toplinepos, _   = messagesView.GetScrollOffset()
		_, _, _, height = messagesView.GetInnerRect()
	)

	if toplinepos == 0 {
		height = 0
		go loadMore()
	}

	current := toplinepos + height

	input.SetPlaceholder(fmt.Sprintf(
		"%d/%d %d%%",
		current, lines, min(current*100/lines, 100),
	))
}
