package main

import (
	"fmt"

	"github.com/diamondburned/tcell"
	"github.com/diamondburned/tview"
)

func promptLogin(text string, mfa bool) (u, p, mfa string, ok bool) {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(tview.NewTextView().SetText(text), 1, 1, false)

	f := tview.NewForm()
	f.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
		}

		return nil
	})

	f.AddInputField(
		"Email", "",
		35, nil, func(s string) { u = s },
	)

	f.AddPasswordField(
		"Password", "",
		35, '*', func(s string) { p = s },
	)

	f.AddButton("Login", func() {
		ok = true
		app.Stop()
	})

	f.SetCancelFunc(func() {
		app.Stop()
	})

	flex.AddItem(f, 0, 1, true)

	frame := tview.NewFrame(flex)
	frame.SetBorders(5, 5, 0, 0, 20, 20)

	if err := app.SetRoot(frame, true).Run(); err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	return
}
