package main

import (
	"github.com/rumblefrog/discordgo"
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

func getChannelFromGuildSettings(chID int64, s *discordgo.UserGuildSettings,
) *discordgo.UserGuildSettingsChannelOverride {

	if s == nil {
		return nil
	}

	for _, c := range s.ChannelOverrides {
		if c.ChannelID == chID {
			return c
		}
	}

	return nil
}

func settingChannelIsMuted(
	cho *discordgo.UserGuildSettingsChannelOverride,
	s *discordgo.UserGuildSettings) (m bool) {

	if cho == nil {
		return false
	}

	m = cho.Muted

	c, err := d.State.Channel(cho.ChannelID)
	if err != nil {
		return
	}

	if c.ParentID != 0 {
		if cs := getChannelFromGuildSettings(c.ParentID, s); cs != nil {
			return cs.Muted
		}
	}

	return
}
