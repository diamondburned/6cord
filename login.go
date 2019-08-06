package main

import (
	"fmt"

	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tcell"
	"github.com/diamondburned/tview/v2"
	"gitlab.com/diamondburned/6cord/center"
)

func promptLogin(l *discordgo.Login, text string, mfa bool) (ok bool) {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(tview.NewTextView().SetText(text), 1, 1, false)
	flex.SetBackgroundColor(-1)

	f := tview.NewForm()
	f.SetBackgroundColor(tcell.Color237)
	f.SetButtonBackgroundColor(tcell.Color255)
	f.SetButtonTextColor(tcell.Color237)

	// Field
	f.SetFieldBackgroundColor(tcell.Color248)
	f.SetFieldTextColor(tcell.Color237)

	f.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		switch ev.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
		}

		return nil
	})

	if !mfa {
		f.AddInputField(
			"Email   ", l.Email, 59, nil,
			func(s string) { l.Email = s },
		)

		f.AddPasswordField(
			"Password", l.Password, 59, '*',
			func(s string) { l.Password = s },
		)
	} else {
		f.AddInputField(
			"MFA     ", l.MFA, 59, nil,
			func(s string) { l.MFA = s },
		)
	}

	f.AddButton("Login", func() {
		ok = true
		app.Stop()
	})

	f.SetCancelFunc(func() {
		app.Stop()
	})

	flex.AddItem(f, 0, 1, true)

	center := center.New(flex)
	center.MaxWidth = 70
	center.MaxHeight = 10

	app.SetRoot(center, true)
	if err := app.Run(); err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	return
}
