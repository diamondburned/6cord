package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/rumblefrog/discordgo"
)

func loadChannel(channelID int64) {
	ch, err := d.State.Channel(channelID)
	if err != nil {
		ch, err = d.Channel(ChannelID) // todo: state first
		if err != nil {
			Warn(err.Error())
			return
		}
	}

	switch ch.Type {
	case discordgo.ChannelTypeGuildVoice:
		Message("Voice isn't supported yet :(")
		return
	case discordgo.ChannelTypeGuildCategory:
		return
	}

	ChannelID = ch.ID

	frameTitle := "[#" + ch.Name + "]"

	if ch.Topic != "" {
		frameTitle += " - [" + ch.Topic + "]"
	}

	wrapFrame.SetTitle(frameTitle)
	typing.Reset()

	if us.GetGuildID() != ch.GuildID {
		us.Reset(ch.GuildID)
	}

	msgs, err := d.ChannelMessages(ChannelID, 35, 0, 0, 0)
	if err != nil {
		Warn(err.Error())
		return
	}

	if len(msgs) < 1 {
		// Drop out early if no messages
		return
	}

	// reverse
	for i := len(msgs)/2 - 1; i >= 0; i-- {
		opp := len(msgs) - 1 - i
		msgs[i], msgs[opp] = msgs[opp], msgs[i]
	}

	go func(c *discordgo.Channel, msgs []*discordgo.Message) {
		var ackMsg *discordgo.Message

		for i := len(msgs) - 1; i >= 0; i-- {
			if msgs[i].Author.ID != d.State.User.ID {
				ackMsg = msgs[i]
				break
			}
		}

		if ackMsg == nil {
			return
		}

		ackMe(c, ackMsg)
		checkReadState()
	}(ch, msgs)

	//var wg sync.WaitGroup
	messageStore = []string{}

	for i, m := range msgs {
		//wg.Add(1)
		//go func(m *discordgo.Message, i int) {
		//defer wg.Done()

		if rstore.Check(m.Author, RelationshipBlocked) {
			continue
		}

		if !isRegularMessage(m) {
			continue
		}

		sentTime, err := m.Timestamp.Parse()
		if err != nil {
			sentTime = time.Now()
		}

		if i > 0 && msgs[i-1].Author.ID != m.Author.ID {
			username, color := us.DiscordThis(m)

			messageStore = append(messageStore, fmt.Sprintf(
				authorFormat,
				color, username,
				sentTime.Format(time.Stamp),
			))
		}

		messageStore = append(messageStore, fmt.Sprintf(
			messageFormat,
			m.ID, fmtMessage(m),
		))

		//}(m, i)
	}

	//wg.Wait()

	messagesView.Clear()
	messagesView.Write([]byte(
		strings.Join(messageStore, ""),
	))

	app.Draw()

	setLastAuthor(msgs[len(msgs)-1].Author.ID)

	messagesView.ScrollToEnd()

	app.SetFocus(input)

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
				m.User.ID,
				m.User.Username,
				m.Nick,
				m.User.Discriminator,
				color,
			)
		}
	}()
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

func scrollChat() {
	if !messagesView.HasFocus() {
		messagesView.ScrollToEnd()
	}
}
