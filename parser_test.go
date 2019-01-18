package main

import (
	"testing"

	"github.com/rivo/tview"
)

func TestParser(t *testing.T) {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetText(parseMD(`**bold *italics*** *italics only*`))

	if err := app.SetRoot(textView, true).SetFocus(textView).Run(); err != nil {
		panic(err)
	}
}
