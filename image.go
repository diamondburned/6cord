package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	img "gitlab.com/diamondburned/6cord/image"
)

type imageCtx struct {
	state img.Backend
}

var imageClient = &http.Client{
	Timeout: 30 * time.Second,
}

func newDiscordImageContext(url string, w, h int) *imageCtx {
	var (
		resizeW int
		resizeH int
	)

	if w > h {
		resizeH = cfg.Prop.ImageHeight
		resizeW = cfg.Prop.ImageHeight * w / h
	} else {
		resizeW = cfg.Prop.ImageWidth
		resizeH = cfg.Prop.ImageWidth * h / w
	}

	url = strings.Split(url, "?")[0] + fmt.Sprintf(
		"?width=%d&height=%d",
		resizeW, resizeH,
	)

	r, err := imageClient.Get(url)
	if err != nil {
		return nil
	}

	defer r.Body.Close()

	i, _, err := image.Decode(r.Body)
	if err != nil {
		return nil
	}

	c, err := img.New(i)
	if err != nil {
		return nil
	}

	return &imageCtx{
		state: c,
	}
}

func (c *imageCtx) Delete() {
	c.state.Delete()
}
