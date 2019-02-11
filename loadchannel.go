package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/RumbleFrog/discordgo"
)

func loadChannel() {
	ch, err := d.Channel(ChannelID) // todo: state first
	if err != nil {
		log.Panicln(err)
	}

	wrapFrame.SetTitle("#" + ch.Name)
	typing.Reset()

	if us.GetGuildID() != ch.GuildID {
		us.Reset(ch.GuildID)
	}

	msgs, err := d.ChannelMessages(ChannelID, 35, 0, 0, 0)
	if err != nil {
		log.Panicln(err)
	}

	// reverse
	for i := len(msgs)/2 - 1; i >= 0; i-- {
		opp := len(msgs) - 1 - i
		msgs[i], msgs[opp] = msgs[opp], msgs[i]
	}

	var (
		messages = make([]string, len(msgs))
		wg       sync.WaitGroup
	)

	for i, m := range msgs {
		wg.Add(1)
		go func(m *discordgo.Message, i int) {
			defer wg.Done()

			if rstore.Check(m.Author, RelationshipBlocked) {
				return
			}

			sentTime, err := m.Timestamp.Parse()
			if err != nil {
				sentTime = time.Now()
			}

			var msg string
			if i > 0 && msgs[i-1].Author.ID != m.Author.ID {
				username, color := us.DiscordThis(m)

				msg = fmt.Sprintf(
					authorFormat,
					color, username,
					sentTime.Format(time.Stamp),
				)
			}

			msg += fmt.Sprintf(
				messageFormat,
				m.ID, fmtMessage(m),
			)

			messages[i] = msg
		}(m, i)
	}

	wg.Wait()

	messagesView.Clear()
	messagesView.Write([]byte(
		strings.Join(messages, ""),
	))

	setLastAuthor(msgs[len(msgs)-1].Author.ID)

	messagesView.ScrollToEnd()

	app.SetFocus(input)

	go func() {
		members := &([]*discordgo.Member{})
		recurseMembers(members, ch.GuildID, 0)

		guild, err := d.State.Guild(ch.GuildID)
		if err != nil {
			if guild, err = d.Guild(ch.GuildID); err != nil {
				log.Println(err)
				return
			}
		}

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

			us.AddUser(
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
