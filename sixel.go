package main

import (
	"bytes"
	"io"
	"net/http"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/BurntSushi/graphics-go/graphics"
	sixel "github.com/mattn/go-sixel"
)

// SixelImage takes in a raw image, parse it and return a sixel bufio
func SixelImage(r io.Reader, w, h uint) (wr bytes.Buffer, err error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return
	}

	if w > 0 || h > 0 {
		rx := float64(img.Bounds().Dx()) / float64(w)
		ry := float64(img.Bounds().Dy()) / float64(h)
		if rx < ry {
			w = uint(float64(img.Bounds().Dx()) / ry)
		} else {
			h = uint(float64(img.Bounds().Dy()) / rx)
		}

		tmp := image.NewNRGBA64(image.Rect(0, 0, int(w), int(h)))

		if err := graphics.Scale(tmp, img); err != nil {
			return wr, err
		}

		img = tmp
	}

	enc := sixel.NewEncoder(&wr)
	enc.Dither = false

	return wr, enc.Encode(img)
}

func SixelFromURL(url string, w, h uint) (wr bytes.Buffer, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	return SixelImage(resp.Body, w, h)

}

func Printable(wr bytes.Buffer) string {
	return wr.String()
}
