package image

import (
	"github.com/mattn/go-tty"
	"gitlab.com/diamondburned/ueberzug-go"
)

func ttySize(t *tty.TTY) (int, int) {
	_, _, w, h, _ := t.SizePixel()
	return w, h
}

func escapeSize(t *tty.TTY) (int, int) {
	return getTermSize(t)
}

func xSize(t *tty.TTY) (int, int) {
	w, h, err := ueberzug.GetParentSize()
	if err != nil {
		return 0, 0
	}

	return w, h
}
