package main

// Source: https://gist.github.com/rumblefrog/c9ebd9fb84a8955495d4fb7983345530

import (
	"sort"

	"github.com/rumblefrog/discordgo"
)

type ChannelGeneric struct {
	Underlying *discordgo.Channel

	Children []*discordgo.Channel
}

func SortChannels(cs []*discordgo.Channel) (out []*discordgo.Channel) {
	p := make(map[int64]*ChannelGeneric)

	for _, c := range cs {
		if c.Type != discordgo.ChannelTypeGuildCategory && c.ParentID != 0 {
			v, ok := p[c.ParentID]

			if ok {
				v.Children = append(v.Children, c)
			} else {
				p[c.ParentID] = &ChannelGeneric{
					Children: []*discordgo.Channel{c},
				}
			}

			continue
		}

		if c.Type == discordgo.ChannelTypeGuildCategory {
			v, ok := p[c.ID]

			if ok {
				v.Underlying = c
			} else {
				p[c.ID] = &ChannelGeneric{
					Underlying: c,
				}
			}

			continue
		}

		p[c.ID] = &ChannelGeneric{
			Underlying: c,
		}
	}

	a := make([]*ChannelGeneric, 0, len(p))

	for _, v := range p {
		if v.Children != nil {
			sort.Slice(v.Children, func(i, j int) bool {
				return v.Children[i].Position < v.Children[j].Position
			})

			sort.SliceStable(v.Children, func(i, j int) bool {
				if v.Children[i].Type == discordgo.ChannelTypeGuildVoice {
					return true
				}

				if v.Children[j].Type == discordgo.ChannelTypeGuildVoice {
					return true
				}

				return false
			})
		}

		a = append(a, v)
	}

	sort.Slice(a, func(i, j int) bool {
		return a[i].Underlying.Position < a[j].Underlying.Position
	})

	for _, v := range a {
		out = append(out, v.Underlying)

		if v.Children != nil {
			for _, k := range v.Children {
				out = append(out, k)
			}
		}
	}

	return
}
