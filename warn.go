package main

import (
	"fmt"
	"log"
	"runtime"

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
	modal.SetText(content + c)
}
