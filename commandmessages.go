package main

func highlightMessage(text []string) {
	if len(text) != 2 {
		messagesView.Highlight()
		return
	}

	messagesView.Highlight(text[1])
	messagesView.ScrollToHighlight()
}
