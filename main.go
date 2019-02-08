package main

import (
	"flag"
	"log"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/RumbleFrog/discordgo"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	keyring "github.com/zalando/go-keyring"
)

const (
	AppName = "6cord"
)

var (
	app           = tview.NewApplication()
	guildView     = tview.NewTreeView()
	messagesView  = tview.NewTextView()
	messagesFrame = tview.NewFrame(messagesView)
	input         = tview.NewInputField()

	ChannelID int64

	LastAuthor int64

	d *discordgo.Session
)

func init() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.Stop()
		}

		return event
	})

	guildView.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		// workaround to prevent crash when no root in tree
		return nil
	})

	messagesView.SetRegions(true)
	messagesView.SetWrap(true)
	messagesView.SetWordWrap(true)
	messagesView.SetScrollable(true)
	messagesView.SetDynamicColors(true)

	token := flag.String("t", "", "Discord token (1)")

	username := flag.String("u", "", "Username/Email (2)")
	password := flag.String("p", "", "Password (2)")

	flag.Parse()

	var login []interface{}

	k, err := keyring.Get(AppName, "token")
	if err != nil {
		if err != keyring.ErrNotFound {
			log.Println(err.Error())
		}

		switch {
		case *token != "":
			login = append(login, *token)
		case *username != "", *password != "":
			login = append(login, *username)
			login = append(login, *password)

			if *token != "" {
				login = append(login, *token)
			}
		default:
			log.Fatalln("Token OR username + password missing! Refer to -h")
		}
	} else {
		login = append(login, k)
	}

	d, err = discordgo.New(login...)
	if err != nil {
		log.Panicln(err)
	}

	d.State.MaxMessageCount = 50

	appflex := tview.NewFlex().SetDirection(tview.FlexColumn)

	{ // Left container
		appflex.AddItem(guildView, 0, 1, true)
	}

	{ // Right container
		flex := tview.NewFlex().SetDirection(tview.FlexRow)
		flex.SetBackgroundColor(tcell.ColorDefault)

		input.SetBackgroundColor(tcell.ColorAqua)

		input.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				text := input.GetText()
				if text == "" {
					return
				}

				go func(text string) {
					if _, err := d.ChannelMessageSend(ChannelID, text); err != nil {
						log.Println(err)
					}
				}(text)
			}

			input.SetText("")
		})

		input.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
			if input.GetText() != "" {
				return ev
			}

			switch ev.Key() {
			case tcell.KeyLeft:
				app.SetFocus(guildView)
				return nil
			}

			return ev
		})

		messagesFrame.SetBorders(0, 0, 0, 1, 0, 0)

		flex.AddItem(messagesFrame, 0, 1, false)
		flex.AddItem(input, 1, 1, true)

		appflex.AddItem(flex, 0, 3, true)
	}

	app.SetRoot(appflex, true)
}

func main() {
	logFile, err := os.OpenFile("/tmp/6cord.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0664)
	if err != nil {
		panic(err)
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// if len(os.Args) > 1 {
	// 	switch os.Args[1] {
	// 	case "rmkeyring":
	// 		switch err := keyring.Delete(AppName, "token"); err {
	// 		case nil:
	// 			log.Println("Keyring deleted.")
	// 			return
	// 		default:
	// 			log.Panicln(err)
	// 		}
	// 	}
	// }

	d.AddHandler(onReady)
	d.AddHandler(messageCreate)
	d.AddHandler(messageUpdate)

	if err := d.Open(); err != nil {
		log.Fatalln("Failed to connect to Discord", err.Error())
	}

	defer d.Close()
	defer app.Stop()

	log.Println("Storing token inside keyring...")
	if err := keyring.Set(AppName, "token", d.Token); err != nil {
		log.Println("Failed to set keyring! Continuing anyway...", err.Error())
	}

	if err := app.Run(); err != nil {
		log.Panicln(err)
	}
}
