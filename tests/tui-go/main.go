package main

import (
	"image"
	"os"
	"time"

	_ "image/png"

	tui "github.com/marcusolsson/tui-go"
	sixel "github.com/mattn/go-sixel"
)

func main() {
	path := "/home/diamond/Pictures/peni_ava.png"

	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	tty, err := os.Open("/dev/tty")
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(time.Second * 1)
		if err := sixel.NewEncoder(tty).Encode(img); err != nil {
			panic(err)
		}
	}()

	text := tui.NewTextEdit()
	text.SetText("homo tard")

	box := tui.NewVBox(text)

	ui, err := tui.New(box)
	if err != nil {
		panic(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
