package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/rivo/tview"
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
			app.SetRoot(appflex, true).SetFocus(input)
		}
	})

	app.SetRoot(modal, false).SetFocus(modal)
	app.Draw()
}

// Message prints a system message
func Message(m string) {
	msg := fmt.Sprintf(
		authorFormat,
		16777215, "<!6cord bot>",
		time.Now().Format(time.Stamp),
	)

	var (
		l = strings.Split(m, "\n")
		c []string
	)

	for i := 0; i < len(l); i++ {
		c = append(c, "\t"+l[i])
	}

	msg += fmt.Sprintf(
		messageFormat+"[::-]",
		0, strings.Join(c, "\n"),
	)

	messagesView.Write([]byte(msg))
	messageStore = append(messageStore, msg)

	scrollChat()

	setLastAuthor(0)
}
