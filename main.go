package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tcell"
	"github.com/diamondburned/tview"
	"github.com/valyala/fasttemplate"
	keyring "github.com/zalando/go-keyring"
	"gitlab.com/diamondburned/6cord/image"
	"gitlab.com/diamondburned/6cord/shortener"
)

const (
	// AppName used for keyrings
	AppName = "6cord"
)

var (
	app           = tview.NewApplication()
	appflex       = tview.NewFlex()
	rightflex     = tview.NewFlex()
	guildView     = tview.NewTreeView()
	messagesView  = tview.NewTextView()
	messagesFrame = tview.NewFrame(messagesView)
	wrapFrame     *tview.Frame
	input         = tview.NewInputField()
	autocomp      = tview.NewList()

	// Channel stores the current channel's pointer
	Channel *discordgo.Channel

	// LastAuthor stores for appending messages
	// TODO: migrate to table + lastRow
	LastAuthor int64

	d *discordgo.Session

	prefixTpl *fasttemplate.Template
)

func init() {
	// less aggressive garbage collector
	// debug.SetGCPercent(200)

	if err := loadCfg(); err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	t, err := fasttemplate.NewTemplate(cfg.Prop.CommandPrefix, "${", "}")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	prefixTpl = t

	app.SetBeforeDrawFunc(func(s tcell.Screen) bool {
		if cfg.Prop.BackgroundColor == -1 {
			s.Clear()
		}

		if redrawDisabled {
			return true
		}

		return false
	})

	commands = append(commands, CustomCommands...)
}

func main() {
	tview.Borders.HorizontalFocus = tview.Borders.Horizontal
	tview.Borders.VerticalFocus = tview.Borders.Vertical

	tview.Borders.TopLeftFocus = tview.Borders.TopLeft
	tview.Borders.TopRightFocus = tview.Borders.TopRight
	tview.Borders.BottomLeftFocus = tview.Borders.BottomLeft
	tview.Borders.BottomRightFocus = tview.Borders.BottomRight

	tview.Borders.Horizontal = ' '
	tview.Borders.Vertical = ' '

	tview.Borders.TopLeft = ' '
	tview.Borders.TopRight = ' '
	tview.Borders.BottomLeft = ' '
	tview.Borders.BottomRight = ' '

	guildView.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if ev.Rune() == '/' {
			app.SetFocus(input)
			input.SetText("/")
		}

		// workaround to prevent crash when no root in tree
		return nil
	})

	messagesView.SetRegions(true)
	messagesView.SetWrap(true)
	messagesView.SetWordWrap(false)
	messagesView.SetScrollable(true)
	messagesView.SetDynamicColors(true)
	messagesView.SetTextColor(tcell.Color(cfg.Prop.ForegroundColor))
	messagesView.SetBackgroundColor(tcell.Color(cfg.Prop.BackgroundColor))
	messagesView.SetText(` [::b]Quick Start[::-]
    - Right arrow or C-l from the server list to focus to input
    - Left arrow or C-h from input to focus to the server list
    - Up arrow from input to go to autocomplete/message scrollback
    - Tab to show/hide channels
    - /goto [#channel] jumps to that channel
	- Page Up/Down jumps between the server entries`)

	var (
		login []interface{}
		err   error
	)

	switch {
	case cfg.Token != "":
		login = append(login, cfg.Token)

		if err := keyring.Delete(AppName, "token"); err == nil {
			log.Println("Keyring deleted.")
		}

	case cfg.Username != "" && cfg.Password != "":
		login = append(login, cfg.Username)
		login = append(login, cfg.Password)

		if cfg.Token != "" {
			login = append(login, cfg.Token)
		}

	default:
		k, err := keyring.Get(AppName, "token")
		if err != nil {
			fmt.Println("Missing token OR username and password!")
			panic(err)
		}

		login = append(login, k)
	}

	d, err = discordgo.New(login...)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	d.UserAgent = `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3534.4 Safari/537.36`

	d.State.MaxMessageCount = 50

	// Main app page

	appflex.SetDirection(tview.FlexColumn)
	appflex.SetBackgroundColor(tcell.Color(cfg.Prop.BackgroundColor))

	{ // Left container
		guildView.SetPrefixes([]string{"", ""})
		guildView.SetTopLevel(1)
		guildView.SetAlign(false)
		guildView.SetBorder(true)
		guildView.SetBorderAttributes(tcell.AttrDim)
		guildView.SetBorderPadding(0, 0, 1, 0)
		guildView.SetBorderColor(tcell.Color(cfg.Prop.ForegroundColor))
		guildView.SetTitle("[Servers[]")
		guildView.SetTitleAlign(tview.AlignLeft)

		guildView.SetBackgroundColor(tcell.Color(cfg.Prop.BackgroundColor))
		guildView.SetGraphicsColor(tcell.Color(cfg.Prop.ForegroundColor))
		guildView.SetTitleColor(tcell.Color(cfg.Prop.ForegroundColor))

		guildView.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
			return nil
		})
	}

	{ // Right container
		rightflex.SetDirection(tview.FlexRow)
		rightflex.SetBackgroundColor(tcell.Color(cfg.Prop.BackgroundColor))

		wrapFrame = tview.NewFrame(rightflex)
		wrapFrame.SetBorder(true)
		wrapFrame.SetBorderAttributes(tcell.AttrDim)
		wrapFrame.SetBorders(0, 0, 0, 0, 0, 0)
		wrapFrame.SetBorderColor(tcell.Color(cfg.Prop.ForegroundColor))
		wrapFrame.SetTitle("")
		wrapFrame.SetTitleAlign(tview.AlignLeft)
		wrapFrame.SetTitleColor(tcell.Color(cfg.Prop.ForegroundColor))
		wrapFrame.SetBackgroundColor(tcell.Color(cfg.Prop.BackgroundColor))

		autocomp.ShowSecondaryText(false)
		autocomp.SetHighlightFullLine(true)
		autocomp.SetBackgroundColor(tcell.Color(cfg.Prop.BackgroundColor))
		autocomp.SetMainTextColor(tcell.Color(cfg.Prop.ForegroundColor))
		autocomp.SetSelectedTextColor(tcell.Color(15 - cfg.Prop.ForegroundColor))
		autocomp.SetSelectedBackgroundColor(tcell.Color(cfg.Prop.ForegroundColor))
		autocomp.SetShortcutColor(tcell.Color(cfg.Prop.ForegroundColor))
		autocomp.SetInputCapture(autocompHandler)

		resetInputBehavior()
		input.SetInputCapture(inputKeyHandler)
		input.SetFieldTextColor(tcell.Color(cfg.Prop.ForegroundColor))

		input.SetChangedFunc(func(text string) {
			if len(text) == 0 {
				clearList()
				stateResetter()
				return
			}

			if text == "/" {
				fuzzyCommands(text)
				return
			}

			if string(text[len(text)-1]) == " " {
				clearList()
				stateResetter()
				return
			}

			words := strings.Fields(text)

			if len(words) < 1 {
				clearList()
				stateResetter()
				return
			}

			switch last := words[len(words)-1]; {
			case strings.HasPrefix(last, "@"):
				fuzzyMentions(last)
			case strings.HasPrefix(last, "#"):
				fuzzyChannels(last)
			case strings.HasPrefix(last, ":"):
				fuzzyEmojis(last)
			case strings.HasPrefix(last, "~"):
				fuzzyMessages(last)
			case strings.HasPrefix(text, "/upload "):
				fuzzyUpload(text)
			case strings.HasPrefix(text, "/"):
				if len(words) == 1 {
					fuzzyCommands(text)
				}
			default:
				typingTrigger()
				clearList()
				stateResetter()
			}
		})

		messagesFrame.SetBorders(0, 0, 0, 0, 0, 0)
		messagesFrame.SetBackgroundColor(tcell.Color(cfg.Prop.BackgroundColor))

		rightflex.AddItem(messagesFrame, 0, 1, false)
		rightflex.AddItem(autocomp, 1, 1, true)
		rightflex.AddItem(input, 1, 1, true)
		rightflex.SetBackgroundColor(tcell.Color(cfg.Prop.BackgroundColor))
	}

	messagesView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyPgDn, tcell.KeyPgUp, tcell.KeyUp, tcell.KeyDown:
			handleScroll()
			return event

		case tcell.KeyLeft:
			app.SetFocus(guildView)
			return nil

		case tcell.KeyEnd:
			messagesView.ScrollToEnd()
			return nil
		}

		switch event.Rune() {
		case 'j', 'k':
			handleScroll()
			return event

		case 'g', 'G':
			return event
		}

		resetInputBehavior()
		app.SetFocus(input)
		return nil
	})

	autocomp.SetSelectedFunc(func(i int, a, b string, c rune) {
		autofillfunc(i)
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF5:
			fmt.Print("\033[2J")
			app.Screen.Sync()
			app.Draw()

		case tcell.KeyCtrlC:
			app.Stop()

		case tcell.KeyTab:
			if autocomp.GetItemCount() < 1 {
				toggleChannels()
				app.ForceDraw()

				return nil
			}
		}

		return event
	})

	app.SetRoot(appflex, true)

	toggleChannels()

	// image
	defer image.Close()
	defer imageRendererPipeline.clean()

	logFile, err := os.OpenFile(
		os.TempDir()+"/6cord.log",
		os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC,
		0664,
	)

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	discordgo.Logger = func(msgL, caller int, format string, a ...interface{}) {
		log.Println("Discordgo:", msgL, caller, format, a)

		if cfg.Debug {
			// Unsure if I should have spew as a dependency
			log.Println(spew.Sdump(a))
		}
	}

	d.AddHandler(onReady)
	d.AddHandler(messageCreate)
	d.AddHandler(messageUpdate)
	d.AddHandler(messageDelete)
	d.AddHandler(messageDeleteBulk)
	d.AddHandler(reactionAdd)
	d.AddHandler(reactionRemove)
	d.AddHandler(reactionRemoveAll)
	d.AddHandler(onTyping)
	d.AddHandler(messageAck)
	d.AddHandler(voiceStateUpdate)
	d.AddHandler(userSettingsUpdate)
	d.AddHandler(relationshipAdd)
	d.AddHandler(relationshipRemove)
	d.AddHandler(guildMemberAdd)
	d.AddHandler(guildMemberUpdate)
	d.AddHandler(guildMemberRemove)

	if cfg.Debug {
		d.AddHandler(onTyping)

		d.AddHandler(func(s *discordgo.Session, r *discordgo.Resumed) {
			log.Println(spew.Sdump(r))
		})

		d.AddHandler(func(s *discordgo.Session, dc *discordgo.Disconnect) {
			log.Println(spew.Sdump(dc))
		})

		d.Debug = true
		d.LogLevel = discordgo.LogDebug

		// d.AddHandler(func(s *discordgo.Session, i interface{}) {
		// 	log.Println(spew.Sdump(i))
		// })
	}

	// d.AddHandler(func(s *discordgo.Session, ev *discordgo.Event) {
	// 	log.Println(spew.Sdump(ev))
	// })

	d.StateEnabled = true
	d.State.MaxMessageCount = 35
	d.State.TrackChannels = true
	d.State.TrackEmojis = true
	d.State.TrackMembers = true
	d.State.TrackRoles = true
	d.State.TrackVoice = true
	d.State.TrackPresences = true

	if err := d.Open(); err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	defer d.Close()
	defer app.Stop()

	go messageRenderer()
	go renderCallback()

	imageRendererPipeline = startImageRendererPipeline()

	log.Println("Storing token inside keyring...")
	if err := keyring.Set(AppName, "token", d.Token); err != nil {
		log.Println("Failed to set keyring! Continuing anyway...", err.Error())
	}

	if cfg.Prop.ShortenURL {
		if err := shortener.StartHTTP("127.0.0.1"); err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	}

	// Stored in syscall.go, only does something when target OS is Linux
	syscallSilenceStderr(logFile)

	if err := app.Run(); err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
