package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/diamondburned/tview/v2"
)

// Warn ..
func Warn(c string) {
	var content string

	{
		_, fn, line, ok := runtime.Caller(1)
		if ok {
			content += fmt.Sprintf("%s:%d ->", fn, line)
		}
	}

	{
		_, fn, line, ok := runtime.Caller(2)
		if ok {
			content += fmt.Sprintf(" %s:%d -> ", fn, line)
		}
	}

	log.Println(content + c)

	modal := tview.NewModal()
	modal.AddButtons([]string{"mkay"})
	modal.SetText(c)
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "mkay" {
			app.SetRoot(appgrid, true)
			app.SetFocus(input)
		}
	})

	app.SetRoot(modal, false)
	app.SetFocus(modal)
	app.Draw()
}

// Message prints a system message
func Message(m string) {
	messageRender <- m
}
