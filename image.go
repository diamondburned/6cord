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

	"github.com/diamondburned/discordgo"
	img "gitlab.com/diamondburned/6cord/image"
)

type imageCtx struct {
	state  img.Backend
	index  int
	assets []*imageAsset
}

type imageAsset struct {
	url  string
	w, h int

	sizedURL string
}

var imageClient = &http.Client{
	Timeout: 30 * time.Second,
}

func newDiscordImageContext(m *discordgo.Message) *imageCtx {
	ctx := &imageCtx{
		assets: make([]*imageAsset, 0, len(m.Attachments)+len(m.Embeds)),
	}

	for _, a := range m.Attachments {
		ctx.assets = append(ctx.assets, &imageAsset{
			url: a.ProxyURL,
			w:   a.Width,
			h:   a.Height,
		})
	}

	for _, e := range m.Embeds {
		if t := e.Thumbnail; t != nil {
			ctx.assets = append(ctx.assets, &imageAsset{
				url: t.ProxyURL,
				w:   t.Width,
				h:   t.Height,
			})
		}
	}

	if len(ctx.assets) == 0 {
		return nil
	}

	if err := ctx.showImage(ctx.assets[0]); err != nil {
		return nil
	}

	return ctx
}

func (ctx *imageCtx) nextImage() error {
	ctx.index++
	if ctx.index >= len(ctx.assets) {
		ctx.index = 0
	}

	return ctx.showImage(ctx.assets[ctx.index])
}

func (ctx *imageCtx) prevImage() error {
	ctx.index--
	if ctx.index < 0 {
		ctx.index = len(ctx.assets) - 1
	}

	return ctx.showImage(ctx.assets[ctx.index])
}

func (ctx *imageCtx) showImage(a *imageAsset) error {
	var (
		resizeW int
		resizeH int
	)

	if a.w > a.h {
		resizeH = cfg.Prop.ImageHeight
		resizeW = cfg.Prop.ImageHeight * a.w / a.h
	} else {
		resizeW = cfg.Prop.ImageWidth
		resizeH = cfg.Prop.ImageWidth * a.h / a.w
	}

	if a.sizedURL == "" {
		a.sizedURL = strings.Split(a.url, "?")[0] + fmt.Sprintf(
			"?width=%d&height=%d",
			resizeW, resizeH,
		)
	}

	r, err := imageClient.Get(a.sizedURL)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	i, _, err := image.Decode(r.Body)
	if err != nil {
		return err
	}

	if ctx.state != nil {
		ctx.state.Delete()
	}

	c, err := img.New(i)
	if err != nil {
		return err
	}

	ctx.state = c

	return err
}

func (c *imageCtx) Delete() {
	c.state.Delete()
}
