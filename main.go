package main

import (
	"flag"
	"log"
	"os"
	"strconv"

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
	messagesView  = tview.NewTextView()
	messagesFrame = tview.NewFrame(messagesView)

	ChannelID int64 = 0

	LastAuthor int64 = 0

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

	messagesView.SetWrap(true)
	messagesView.SetWordWrap(true)
	messagesView.SetScrollable(false)
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

	if len(flag.Args()) != 1 {
		panic("Invalid args! First arg should be ChannelID!")
	}

	ChannelID, err = strconv.ParseInt(flag.Args()[0], 10, 64)
	if err != nil {
		panic(err)
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBackgroundColor(tcell.ColorDefault)

	input := tview.NewInputField()

	input.SetBackgroundColor(tcell.ColorDefault)
	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			d.ChannelMessageSend(ChannelID, input.GetText())
		}

		input.SetText("")
	})

	messagesFrame.SetBorders(2, 2, 0, 0, 2, 2)
	flex.AddItem(messagesFrame, 0, 1, false)
	flex.AddItem(input, 1, 1, true)

	app.SetRoot(flex, true)
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

	// go func(ml *tview.List) {
	// 	time.Sleep(time.Second * 2)
	// 	f, _ := os.Open("/home/diamond/Pictures/emoji.jpg")
	// 	defer f.Close()

	// 	img, _, _ := image.Decode(f)

	// 	tmp := image.NewNRGBA64(image.Rect(0, 0, int(20), int(20)))
	// 	_ = graphics.Scale(tmp, img)

	// 	img = tmp

	// 	app.QueueUpdateDraw(func() {
	// 		ml.AddItem(
	// 			"ym555#5555",
	// 			"uwu im gay",
	// 			' ', nil,
	// 		)
	// 	})
	// }(messagesList)

}
