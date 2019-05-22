package main

import (
	"fmt"
	"image"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/discordgo"
)

type imageCacheStruct struct {
	sync.Mutex

	client *http.Client
	store  map[int64][]*imageCacheAsset
	lastCh int64
}

type imageCacheAsset struct {
	url      string
	w, h     int
	sizedURL string

	i image.Image
}

var imageCache = &imageCacheStruct{
	client: &http.Client{
		Timeout: 30 * time.Second,
	},
}

func (c *imageCacheStruct) get(m int64) []*imageCacheAsset {
	c.Lock()
	defer c.Unlock()

	if a, ok := c.store[m]; ok {
		return a
	}

	return nil
}

func (c *imageCacheStruct) calcURL(a *imageCacheAsset) {
	var (
		resizeW int
		resizeH int
	)

	if a.w < a.h {
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
}

func (c *imageCacheStruct) set(m *discordgo.Message) ([]*imageCacheAsset, error) {
	if c.lastCh != m.ChannelID {
		c.reset()
	}

	assets := make(
		[]*imageCacheAsset,
		0, len(m.Attachments)+len(m.Embeds),
	)

	for _, a := range m.Attachments {
		a := &imageCacheAsset{
			url: a.ProxyURL,
			w:   a.Width,
			h:   a.Height,
		}

		c.calcURL(a)
		assets = append(assets, a)
	}

	for _, e := range m.Embeds {
		if t := e.Thumbnail; t != nil {
			a := &imageCacheAsset{
				url: t.ProxyURL,
				w:   t.Width,
				h:   t.Height,
			}

			c.calcURL(a)
			assets = append(assets, a)
		}
	}

	if len(assets) == 0 {
		return nil, nil
	}

	for _, a := range assets {
		r, err := c.client.Get(a.sizedURL)
		if err != nil {
			return nil, err
		}

		i, _, err := image.Decode(r.Body)
		if err != nil {
			r.Body.Close()
			return nil, err
		}

		a.i = i
		r.Body.Close()
	}

	c.Lock()
	defer c.Unlock()

	c.store[m.ID] = assets

	return assets, nil
}

func (c *imageCacheStruct) reset() {
	c.Lock()
	defer c.Unlock()

	c.store = map[int64][]*imageCacheAsset{}
}
