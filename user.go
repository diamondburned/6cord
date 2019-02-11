package main

import (
	"log"
	"sort"

	"github.com/rivo/tview"

	"github.com/RumbleFrog/discordgo"
)

func safeAuthor(u *discordgo.User) (string, int64) {
	if u != nil {
		return u.Username, u.ID
	}

	return "invalid user", 0
}

func getUserData(u *discordgo.User, chID int64) (name string, color int) {
	color = 16711422
	name, id := safeAuthor(u)

	if d == nil {
		return
	}

	channel, err := d.State.Channel(chID)
	if err != nil {
		if channel, err = d.Channel(chID); err != nil {
			log.Println(err)
			return
		}
	}

	if channel.GuildID == 0 {
		return
	}

	guild, err := d.State.Guild(channel.GuildID)
	if err != nil {
		if guild, err = d.Guild(channel.GuildID); err != nil {
			log.Println(err)
			return
		}
	}

	member, err := d.State.Member(guild.ID, id)
	if err != nil {
		if member, err = d.GuildMember(channel.GuildID, id); err != nil {
			log.Println(err)
			return
		}
	}

	if member.Nick != "" {
		name = tview.Escape(member.Nick)
	}

	roles := guild.Roles
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Position > roles[j].Position
	})

	for _, role := range roles {
		for _, roleID := range member.Roles {
			if role.ID == roleID && role.Color != 0 {
				color = role.Color
				return
			}
		}
	}

	return
}
