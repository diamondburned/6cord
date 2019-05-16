package image

import (
	"errors"

	"github.com/mattn/go-tty"
)

var t *tty.TTY

var (
	PixelW int
	PixelH int
)

func Listen() (err error) {
	t, err = tty.Open()
	if err != nil {
		return err
	}

	var (
		fn func(*tty.TTY) (int, int)
		w  int
		h  int
	)

	if w, h = ttySize(t); w > 0 && h > 0 {
		fn = ttySize
	} else if w, h = xSize(t); w > 0 && h > 0 {
		fn = xSize
	} else if w, h = escapeSize(t); w > 0 && h > 0 {
		fn = escapeSize
	}

	if fn == nil {
		return errors.New("No method of getting terminal size avilable")
	}

	PixelW = w
	PixelH = h

	go func() {
		winch := t.SIGWINCH()
		for range winch {
			PixelW, PixelH = fn(t)
		}
	}()

	return nil
}
