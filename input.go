package main

import (
	"bytes"
	"log"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/diamondburned/tcell"
	"github.com/diamondburned/tview"
	"gitlab.com/diamondburned/6cord/antitele"
)

var toReplaceRuneMap = map[byte]string{
	'n': "\n",
	't': "    ",
}

var (
	currentHistoryItem = -1
	currentMessage     string
)

func resetInputBehavior() {
	app.QueueUpdate(func() {
		input.SetLabel(templatePrefix())
		input.SetLabelColor(tcell.Color(cfg.Prop.BackgroundColor))
		input.SetBackgroundColor(tcell.Color(cfg.Prop.BackgroundColor))
		input.SetFieldBackgroundColor(tcell.Color(cfg.Prop.BackgroundColor))
		input.SetPlaceholderTextColor(tcell.ColorDarkCyan)
		input.SetText("")
	})
	clearList()

	stateResetter()
	toEditMessage = 0
}

func templatePrefix() string {
	var (
		channelTpl  = "#nil"
		guildTpl    = "nil"
		selfnameTpl = "nil"
		discrimTpl  = "0000"
	)

	if Channel != nil {
		switch {
		case Channel.Name == "":
			channelTpl = "#dm"
		default:
			channelTpl = "#" + Channel.Name
		}

		g, _ := d.State.Guild(Channel.GuildID)
		if g != nil {
			guildTpl = g.Name
		}
	}

	if d != nil && d.State.User != nil {
		selfnameTpl = d.State.User.Username
		discrimTpl = d.State.User.Discriminator
	}

	return prefixTpl.ExecuteString(map[string]interface{}{
		"CHANNEL":  channelTpl,
		"GUILD":    guildTpl,
		"USERNAME": selfnameTpl,
		"DISCRIM":  discrimTpl,
	})
}

func processString(input string) string {
	var output = strings.Builder{}

RuneWalk:
	for i := 0; i < len(input); i++ {
		for match, with := range toReplaceRuneMap {
			if input[i] == '\\' && (i < len(input)-1 && input[i+1] == match) {
				if i == 0 || input[i-1] != '\\' {
					output.WriteString(with)
					if i < len(input)-2 && input[i+2] == ' ' {
						i++
					}

					i++
					continue RuneWalk
				} else {
					i++
				}
			}
		}

		output.WriteByte(input[i])
	}

	if cfg.Prop.ObfuscateWords {
		return antitele.Insert(output.String())
	}

	return output.String()
}

var store bool

func handleHistoryItem() {
	if currentHistoryItem < 0 {
		input.SetText(currentMessage)
		store = true
	} else {
		if store {
			currentMessage = input.GetText()
			store = false
		}

		input.SetText(cmdHistory[currentHistoryItem])
	}
}

func inputKeyHandler(ev *tcell.EventKey) *tcell.EventKey {
	if ev.Modifiers() == tcell.ModAlt {
		switch {
		case ev.Key() == tcell.KeyEnter:
			input.SetText(input.GetText() + "\n")
			return nil

		case ev.Key() == tcell.KeyDown || ev.Rune() == 'j':
			currentHistoryItem++
			if currentHistoryItem > len(cmdHistory)-1 {
				currentHistoryItem = len(cmdHistory) - 1
			}

			handleHistoryItem()
			return nil

		case ev.Key() == tcell.KeyUp || ev.Rune() == 'k':
			currentHistoryItem--
			if currentHistoryItem < -1 {
				currentHistoryItem = -1
			}

			handleHistoryItem()
			return nil
		}
	} else {
		acItem, _ := autocomp.GetCurrentItem()

		switch ev.Key() {
		case tcell.KeyEscape:
			resetInputBehavior()
			if showChannels {
				app.QueueUpdateDraw(func() {
					toggleChannels()
				})
			}

			app.SetFocus(guildView)
			return nil

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

							input.SetPlaceholder(cfg.Prop.DefaultStatus)

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
		case tcell.KeyLeft, tcell.KeyCtrlH:
			if input.GetText() != "" {
				return ev
			}

			app.SetFocus(guildView)
			return nil

		case tcell.KeyUp:
			if autocomp.GetItemCount() < 1 {
				app.SetFocus(messagesView)
			} else {
				if acItem == 0 {
					newitem := autocomp.GetItemCount() - 1
					autocomp.SetCurrentItem(newitem)
				}

				app.SetFocus(autocomp)
			}

		case tcell.KeyDown:
			var newitem = acItem + 1

			switch {
			case autocomp.GetItemCount() == 0:
				return ev
			case newitem > autocomp.GetItemCount()-1:
				newitem = 0
			}

			autocomp.SetCurrentItem(newitem)
			app.SetFocus(autocomp)

		case tcell.KeyTab:
			if autocomp.GetItemCount() > 0 {
				autofillfunc(acItem)
				return nil
			}

			return ev

		case tcell.KeyEnter:
			if autocomp.GetItemCount() > 0 {
				autofillfunc(acItem)
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
	}

	return ev
}
