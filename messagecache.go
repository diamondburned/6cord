package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/discordgo"
)

type imageCacheStruct struct {
	sync.RWMutex
	Age time.Duration

	client *http.Client
	store  map[int64]*imageCacheStore
	lastCh int64
}

type imageCacheStore struct {
	assets []*imageCacheAsset
	time   time.Time
	state  imageFetchState
}

type imageCacheAsset struct {
	url      string
	w, h     int
	sizedURL string

	i image.Image
}

type imageFetchState string

const (
	imageFetching imageFetchState = "[grey]"
)

func (c *imageCacheStruct) get(m int64) []*imageCacheAsset {
	c.RLock()
	defer c.RUnlock()

	if a, ok := c.store[m]; ok {
		return a.assets
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

// set checks cache as well
func (c *imageCacheStruct) upd(m *discordgo.Message) ([]*imageCacheAsset, error) {
	if assets := c.get(m.ID); assets != nil {
		return assets, nil
	}

	s := &imageCacheStore{
		assets: make(
			[]*imageCacheAsset,
			0, len(m.Attachments)+len(m.Embeds),
		),
		time: time.Now(),
	}

	for _, a := range m.Attachments {
		if a.Width == 0 || a.Height == 0 {
			continue
		}

		a := &imageCacheAsset{
			url: a.ProxyURL,
			w:   a.Width,
			h:   a.Height,
		}

		c.calcURL(a)
		s.assets = append(s.assets, a)
	}

	for _, e := range m.Embeds {
		if t := e.Thumbnail; t != nil {
			a := &imageCacheAsset{
				url: t.ProxyURL,
				w:   t.Width,
				h:   t.Height,
			}

			c.calcURL(a)
			s.assets = append(s.assets, a)
		}
	}

	if len(s.assets) == 0 {
		return nil, nil
	}

	for _, a := range s.assets {
		r, err := c.client.Get(a.sizedURL)
		if err != nil {
			return nil, err
		}

		i, _, err := image.Decode(r.Body)
		if err == nil {
			a.i = i
		} else {
			// Error is ignored, as skipping a non-supported
			// image is fine
			log.Println("Error on", a.sizedURL, "\n"+err.Error())
		}

		r.Body.Close()
	}

	c.Lock()
	defer c.Unlock()

	c.store[m.ID] = s

	return s.assets, nil
}

func (c *imageCacheStruct) reset() {
	c.Lock()
	defer c.Unlock()

	c.store = map[int64]*imageCacheStore{}
}

func (c *imageCacheStruct) gc() {
	c.Lock()
	defer c.Unlock()

	for k, store := range c.store {
		if time.Now().Sub(store.time) > c.Age {
			delete(c.store, k)
		}
	}
}
