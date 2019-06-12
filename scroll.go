package main

import (
	"fmt"
	"strings"
)

func handleScroll() {
	current, lines := getLineStatus()

	if current == 0 {
		go loadMore()
	}

	input.SetPlaceholder(fmt.Sprintf(
		"%d/%d %d%%",
		current, lines, min(current*100/lines, 100),
	))

	app.Draw()
}

func getLineStatus() (current, total int) {
	total = len(
		strings.Split(
			strings.Join(messageStore, ""),
			"\n",
		),
	)

	if total <= 1 {
		total = len(messagesView.GetText(false))
	}

	var (
		toplinepos, _   = messagesView.GetScrollOffset()
		_, _, _, height = messagesView.GetInnerRect()
	)

	if toplinepos == 0 {
		height = 0
	}

	current = toplinepos + height

	return
}
