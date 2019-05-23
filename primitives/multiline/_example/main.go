package main

import (
	"github.com/rivo/tview"
	"gitlab.com/diamondburned/6cord/primitives/multiline"
)

func main() {
	tv := tview.NewTextView()
	tv.SetText(`heavy test`)

	m := multiline.NewMultiline()
	m.Placeholder = "Placeholder test"

	f := tview.NewFlex()
	f.SetDirection(tview.FlexRow)
	f.AddItem(tv, 0, 1, false)
	f.AddItem(m, 3, 1, true)

	if err := tview.NewApplication().SetRoot(f, true).SetFocus(f).Run(); err != nil {
		panic(err)
	}
}
