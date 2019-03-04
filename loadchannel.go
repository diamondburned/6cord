package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/rivo/tview"
	"github.com/rumblefrog/discordgo"
)

func loadChannel(channelID int64) {
	wrapFrame.SetTitle("[Loading...]")

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
		Message("Voice is currently not working D:")
		//go toggleVoiceJoin(ch.ID)
		return
	case discordgo.ChannelTypeGuildCategory:
		return
	}

	ChannelID = ch.ID

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

	wrapFrame.SetTitle(tview.Escape(frameTitle))
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

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].ID < msgs[j].ID
	})

	go func(c *discordgo.Channel, msgs []*discordgo.Message) {
		ackMe(msgs[len(msgs)-1])
		checkReadState(msgs[0].ChannelID)
	}(ch, msgs)

	//var wg sync.WaitGroup
	messageStore = []string{}

	for i := 0; i < len(msgs); i++ {
		m := msgs[i]

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
				color, tview.Escape(username),
				sentTime.Format(time.Stamp),
			))
		}

		messageStore = append(messageStore, fmt.Sprintf(
			messageFormat,
			m.ID, fmtMessage(m),
		))

		d.State.MessageAdd(m)

		//}(m, i)
	}

	//wg.Wait()

	setLastAuthor(msgs[len(msgs)-1].Author.ID)

	app.QueueUpdateDraw(func() {
		messagesView.SetText(
			strings.Join(messageStore, ""),
		)
	})

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
