package main

import (
	"image"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/rivo/tview"
)

var (
	messagesList = tview.NewList()
)

func main() {
	app := tview.NewApplication()

	go func(ml *tview.List) {
		time.Sleep(time.Second * 2)
		f, _ := os.Open("/home/diamond/Pictures/emoji.jpg")
		defer f.Close()

		img, _, _ := image.Decode(f)

		tmp := image.NewNRGBA64(image.Rect(0, 0, int(20), int(20)))
		_ = graphics.Scale(tmp, img)

		img = tmp

		app.QueueUpdateDraw(func() {
			ml.AddItem(
				"ym555#5555",
				"uwu im gay",
				' ', nil,
			)
		})
	}(messagesList)

	if err := app.SetRoot(messagesList, true).Run(); err != nil {
		panic(err)
	}
}
