package image

import (
	"errors"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"gitlab.com/diamondburned/go-w3m"
)

type W3M struct {
	a *w3m.Arguments
	f *os.File
}

var Errw3mUnavailable = errors.New("w3m unavailable")

func (w *W3M) Available() error {
	if w3m.GetExecPath() == "" {
		return Errw3mUnavailable
	}

	return nil
}

func (w *W3M) Spawn(i image.Image, x, y int) error {
	bounds := i.Bounds()

	w.a = &w3m.Arguments{
		Width:   bounds.Dx(),
		Height:  bounds.Dy(),
		Xoffset: x,
		Yoffset: y,
	}

	t := strconv.FormatInt(time.Now().UnixNano(), 10)

	f, err := ioutil.TempFile(os.TempDir(), t+".png")
	if err != nil {
		return err
	}

	defer f.Close()

	if err := png.Encode(f, i); err != nil {
		return err
	}

	w.f = f

	return w3m.Spawn(w.a, f.Name())
}

func (w *W3M) Delete() error {
	os.Remove(w.f.Name())
	return w3m.Clear(w.a)
}
