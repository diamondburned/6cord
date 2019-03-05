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

	var (
		toplinepos, _   = messagesView.GetScrollOffset()
		_, _, _, height = messagesView.GetInnerRect()
	)

	if toplinepos == 0 {
		go loadMore()
	}

	current := toplinepos + (height - 2)

	input.SetPlaceholder(fmt.Sprintf(
		"%d/%d %d%%",
		current, lines, current*100/lines,
	))
}
