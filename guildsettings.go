package main

import (
	"github.com/RumbleFrog/discordgo"
)

func getGuildFromSettings(guildID int64) *discordgo.UserGuildSettings {
	for _, ugs := range d.State.UserGuildSettings {
		if ugs.GuildID == guildID {
			return ugs
		}
	}

	return nil
}

func settingGuildIsMuted(s *discordgo.UserGuildSettings) bool {
	if s == nil {
		return false
	}

	return s.Muted
}

func getChannelFromGuildSettings(
	chID int64, s *discordgo.UserGuildSettings,
) *discordgo.UserGuildSettingsChannelOverride {
	for _, c := range s.ChannelOverrides {
		if c.ChannelID == chID {
			return c
		}
	}

	return nil
}

func settingChannelIsMuted(cho *discordgo.UserGuildSettingsChannelOverride) bool {
	if cho == nil {
		return false
	}

	return cho.Muted
}
