package main

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
	"github.com/mattn/go-sixel"
)

var redrawDisabled bool

func commandSixelTest(text []string) {
	r, err := http.Get(text[1])
	if err != nil {
		Warn(err.Error())
		return
	}

	defer r.Body.Close()

	img, _, err := image.Decode(r.Body)
	if err != nil {
		Warn(err.Error())
		return
	}

	if img.Bounds().Dx() > 200 {
		img = imaging.Resize(img, 200, 0, imaging.Linear)
	}

	var b bytes.Buffer

	enc := sixel.NewEncoder(&b)
	enc.Dither = false

	if err := enc.Encode(img); err != nil {
		Warn(err.Error())
		return
	}

	redrawDisabled = true
	fmt.Fprint(os.Stdout, string(b.Bytes()))
}
