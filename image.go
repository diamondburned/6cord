package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

var (
	lastImgCtx *imageCtx
	lastImgMu  = &sync.Mutex{}
)

func checkForImage(ID string) {
	lastImgMu.Lock()
	defer lastImgMu.Unlock()

	if lastImgCtx != nil {
		lastImgCtx.Delete()
	}

	if Channel == nil {
		return
	}

	id, _ := strconv.ParseInt(ID, 10, 64)
	if id == 0 {
		return
	}

	m, err := d.State.Message(Channel.ID, id)
	if err != nil {
		return
	}

	go func() {
		lastImgCtx = newDiscordImageContext(m)
	}()
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

		maxW = min(img.PixelW, cfg.Prop.ImageWidth)
		maxH = min(img.PixelH, cfg.Prop.ImageHeight)
	)

	if a.w < a.h {
		resizeH = maxH
		resizeW = maxH * a.w / a.h
	} else {
		resizeW = maxW
		resizeH = maxW * a.h / a.w
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

	ctx.Delete()

	c, err := img.New(i)
	if err != nil {
		return err
	}

	ctx.state = c

	return err
}

func (c *imageCtx) Delete() {
	if c.state != nil {
		c.state.Delete()
	}
}
