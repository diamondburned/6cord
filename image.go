package main

import (
	"net/http"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/diamondburned/discordgo"
	img "gitlab.com/diamondburned/6cord/image"
)

type imageCtx struct {
	state img.Backend
	index int
}

type imageRendererPipelineStruct struct {
	event   chan interface{}
	state   img.Backend
	message int64
	index   int

	cache  *imageCacheStruct
	assets []*imageCacheAsset
}

const (
	imagePipelineNextEvent int = iota
	imagePipelinePrevEvent
)

var imageRendererPipeline *imageRendererPipelineStruct

func startImageRendererPipeline() *imageRendererPipelineStruct {
	p := &imageRendererPipelineStruct{
		event: make(chan interface{}, prefetchMessageCount),
		cache: &imageCacheStruct{
			Age: 5 * time.Minute,
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
			store: map[int64]*imageCacheStore{},
		},
	}

	go func() {
		for i := range p.event {
		Switch:
			switch i := i.(type) {
			case *discordgo.Message:
				p.message = i.ID

				a := p.cache.get(i.ID)
				if a == nil || a.state != imageFetched {
					p.clean()
					break
				}

				p.assets = a.assets
				p.clean()

				if p.assets == nil {
					break Switch
				}

				p.show()

			case int:
				if p.assets == nil {
					break Switch
				}

				switch i {
				case imagePipelineNextEvent:
					p.index++
					if p.index >= len(p.assets) {
						p.index = 0
					}
				case imagePipelinePrevEvent:
					p.index--
					if p.index < 0 {
						p.index = len(p.assets) - 1
					}
				default:
					break Switch
				}

				p.show()

			default:
				break Switch
			}
		}
	}()

	return p
}

func (p *imageRendererPipelineStruct) add(m *discordgo.Message) {
	p.event <- m
}

func (p *imageRendererPipelineStruct) next() {
	p.event <- imagePipelineNextEvent
}

func (p *imageRendererPipelineStruct) prev() {
	p.event <- imagePipelinePrevEvent
}

func (p *imageRendererPipelineStruct) clean() {
	if p != nil {
		if p.state != nil {
			p.state.Delete()
		}

		p.cache.gc()
		p.index = 0
	}
}

func (p *imageRendererPipelineStruct) show() (err error) {
	p.clean()

	if p.assets == nil {
		return nil
	}

	if p.assets[p.index].i == nil {
		return nil
	}

	p.state, err = img.New(p.assets[p.index].i)
	if err != nil {
		return err
	}

	return nil
}
