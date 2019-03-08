package main

import (
	"log"
	"sort"

	"github.com/rivo/tview"

	"github.com/diamondburned/discordgo"
)

func userSettingsUpdate(s *discordgo.Session, settings *discordgo.UserSettingsUpdate) {
	if settings == nil {
		return
	}

	_settings := *settings

	if status, ok := _settings["status"]; ok {
		if str, ok := status.(string); ok {
			st := discordgo.Status(str)
			d.State.Settings.Status = st
		}
	}
}

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

	member, err := d.State.Member(channel.GuildID, id)
	if err != nil {
		if member, err = d.GuildMember(channel.GuildID, id); err != nil {
			log.Println(err)
			return
		}
	}

	if member.Nick != "" {
		name = tview.Escape(member.Nick)
	}

	color = getUserColor(channel.GuildID, member.Roles)

	return
}

func getUserColor(guildID int64, rls discordgo.IDSlice) int {
	g, err := d.State.Guild(guildID)
	if err != nil {
		if g, err = d.Guild(guildID); err != nil {
			log.Println(err)
			return 16711422
		}
	}

	roles := g.Roles
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Position > roles[j].Position
	})

	for _, role := range roles {
		for _, roleID := range rls {
			if role.ID == roleID && role.Color != 0 {
				return role.Color
			}
		}
	}

	return 16711422
}
