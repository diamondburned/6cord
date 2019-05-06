package main

import (
	"github.com/rivo/tview"
	"gitlab.com/diamondburned/6cord/primitives/multiline"
)

func main() {
	m, _ := multiline.NewMultiline()
	m.Placeholder = "Placeholder test"

	if err := tview.NewApplication().SetRoot(m, true).SetFocus(m).Run(); err != nil {
		panic(err)
	}
}
