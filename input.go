package main

import (
	"bytes"
	"log"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func resetInputBehavior() {
	app.QueueUpdate(func() {
		input.SetLabel(cfg.CommandPrefix)
		input.SetPlaceholder(cfg.DefaultStatus)
		input.SetLabelColor(tcell.Color(cfg.BackgroundColor))
		input.SetBackgroundColor(tcell.Color(cfg.BackgroundColor))
		input.SetFieldBackgroundColor(tcell.Color(cfg.BackgroundColor))
		input.SetPlaceholderTextColor(tcell.ColorDarkCyan)
		input.SetText("")

		clearList()

		stateResetter()
		toEditMessage = 0
	})
}

func inputKeyHandler(ev *tcell.EventKey) *tcell.EventKey {
	switch ev.Key() {
	case tcell.KeyEscape:
		resetInputBehavior()

	case tcell.KeyCtrlV:
		cb, err := clipboard.ReadAll()
		if err != nil {
			log.Println("Couldn't get clipboard:", err)
			return nil
		}

		b := []byte(cb)

		if IsFile(b) {
			modal := tview.NewModal()
			modal.AddButtons([]string{"Cancel", "Yes"})
			modal.SetText("Upload file in clipboard?")
			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				switch buttonLabel {
				case "Yes":
					go func() {
						if Channel == nil {
							Warn("Not in a channel.")
						}

						input.SetPlaceholder("Uploading file...")

						br := bytes.NewReader(b)
						_, err = d.ChannelFileSend(
							Channel.ID,
							"clipboard.png",
							br,
						)

						input.SetPlaceholder(cfg.DefaultStatus)

						if err != nil {
							Warn(err.Error())
						}
					}()
				}

				app.SetRoot(appflex, true).SetFocus(input)
			})

			app.SetRoot(modal, false).SetFocus(modal)

		} else {
			input.SetText(input.GetText() + cb)
		}

	case tcell.KeyLeft:
		if input.GetText() != "" {
			return ev
		}

		app.SetFocus(guildView)
		return nil

	case tcell.KeyUp:
		if autocomp.GetItemCount() < 1 {
			app.SetFocus(messagesView)
		} else {
			if autocomp.GetCurrentItem() == 0 {
				newitem := autocomp.GetItemCount() - 1
				autocomp.SetCurrentItem(newitem)
			}

			app.SetFocus(autocomp)
		}

	case tcell.KeyDown:
		var newitem = autocomp.GetCurrentItem() + 1

		switch {
		case autocomp.GetItemCount() == 0:
			return ev
		case newitem > autocomp.GetItemCount()-1:
			newitem = 0
		}

		autocomp.SetCurrentItem(newitem)
		app.SetFocus(autocomp)

	case tcell.KeyEnter:
		if ev.Modifiers() == tcell.ModAlt {
			input.SetText(input.GetText() + "\n")
			return nil
		}

		if autocomp.GetItemCount() > 0 {
			autofillfunc(0)
			return nil
		}

		// log.Println(ev.Name())

		// if ev.Name() == "Shift+Enter" {
		// 	input.SetText(input.GetText() + "\\n")
		// 	return nil
		// }

		switch input.GetLabel() {
		case EditMessageLabel:
			go editHandler()

		default:
			CommandHandler()
		}

		return nil
	}

	return ev
}
