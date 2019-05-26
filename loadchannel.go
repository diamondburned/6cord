package main

import (
	"log"
	"sort"
	"time"

	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tview"
)

const prefetchMessageCount = 35

func loadChannel(channelID int64) {
	app.QueueUpdateDraw(func() {
		wrapFrame.SetTitle("[Loading...[]")
	})

	go actualLoadChannel(channelID)
}

func actualLoadChannel(channelID int64) {
	ch, err := d.State.Channel(channelID)
	if err != nil {
		ch, err = d.Channel(channelID)
		if err != nil {
			Warn(err.Error())
			return
		}
	}

	switch ch.Type {
	case discordgo.ChannelTypeGuildVoice:
		Message("Voice is currently not working D:")
		return
	case discordgo.ChannelTypeGuildCategory:
		return
	}

	Channel = ch

	typing.Reset()

	if !us.Populated(Channel.GuildID) {
		d.GatewayManager.SubscribeGuild(
			Channel.GuildID, true, true,
		)
	}

	msgs, err := d.ChannelMessages(
		Channel.ID, prefetchMessageCount,
		0, 0, 0,
	)

	if err != nil {
		Warn(err.Error())
		return
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].ID < msgs[j].ID
	})

	if len(msgs) > 0 {
		go func(c *discordgo.Channel, msgs []*discordgo.Message) {
			ackMe(c.ID, msgs[len(msgs)-1].ID)
		}(ch, msgs)
	}

	// Clears the buffer
	messageRender <- nil

	for i := 0; i < len(msgs); i++ {
		m := msgs[i]

		if rstore.Check(m.Author, RelationshipBlocked) {
			continue
		}

		if !isRegularMessage(m) {
			continue
		}

		messageRender <- m
		d.State.MessageAdd(m)
	}

	if len(msgs) > 0 {
		setLastAuthor(msgs[len(msgs)-1].Author.ID)
	} else {
		Message("There's nothing here!")
	}

	resetInputBehavior()
	app.SetFocus(input)

	var frameTitle string

	if ch.Name != "" {
		frameTitle = "[#" + ch.Name + "]"

		if ch.Topic != "" {
			topic, _ := parseEmojis(ch.Topic)
			frameTitle += " - [" + topic + "]"
		}
	} else {
		if len(ch.Recipients) == 1 {
			frameTitle = "[" + ch.Recipients[0].String() + "]"
		} else {
			var names = make([]string, len(ch.Recipients))
			for i, r := range ch.Recipients {
				names[i] = r.Username
			}

			frameTitle = "[" + HumanizeStrings(names) + "]"
		}
	}

	app.QueueUpdateDraw(func() {
		wrapFrame.SetTitle(tview.Escape(frameTitle))
	})

	go func() {
		if ch.GuildID == 0 {
			return
		}

		members := &([]*discordgo.Member{})

		guild, err := d.State.Guild(ch.GuildID)
		if err != nil {
			if guild, err = d.Guild(ch.GuildID); err != nil {
				Warn(err.Error())
				return
			}
		}

		recurseMembers(members, ch.GuildID, 0)

		guild.Members = *members

		roles := guild.Roles
		sort.Slice(roles, func(i, j int) bool {
			return roles[i].Position > roles[j].Position
		})

		for _, m := range *members {
			color := 16711422

		RoleLoop:
			for _, role := range roles {
				for _, roleID := range m.Roles {
					if role.ID == roleID && role.Color != 0 {
						color = role.Color
						break RoleLoop
					}
				}
			}

			us.UpdateUser(
				guild.ID,
				m.User.ID,
				m.User.Username,
				m.Nick,
				m.User.Discriminator,
				color,
			)
		}
	}()
}

func messageisOld(m, l *discordgo.Message) bool {
	if m == nil || l == nil {
		return true
	}

	mt, err := m.Timestamp.Parse()
	if err != nil {
		return true
	}

	lt, err := l.Timestamp.Parse()
	if err != nil {
		return true
	}

	return lt.Add(time.Minute).Before(mt)
}

func recurseMembers(memstore *[]*discordgo.Member, guildID, after int64) {
	members, err := d.GuildMembers(guildID, after, 1000)
	if err != nil {
		log.Println(err)
		return
	}

	if len(members) == 1000 {
		recurseMembers(
			memstore,
			guildID,
			members[999].User.ID,
		)
	}

	*memstore = append(*memstore, members...)

	return
}
