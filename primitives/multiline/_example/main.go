package main

import (
	"github.com/rivo/tview"
	"gitlab.com/diamondburned/6cord/primitives/multiline"
)

func main() {
	m, _ := multiline.NewMultiline()
	m.Placeholder = "Placeholder test"
	m.Insert(`Something
extra---
ordinary!`)

	if err := tview.NewApplication().SetRoot(m, true).SetFocus(m).Run(); err != nil {
		panic(err)
	}
}
