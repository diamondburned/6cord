package main

import (
	"github.com/rumblefrog/discordgo"
)

func messagePingable(m *discordgo.Message, gID int64) bool {
	c, err := d.State.Channel(m.ChannelID)
	if err != nil {
		return false
	}

	if c.GuildID == 0 {
		if m.Author.ID != d.State.User.ID {
			return true
		}

		return false
	}

	s := getGuildFromSettings(c.GuildID)

	if m.MentionEveryone {
		//log.Println(2)
		return settingGuildAllowEveryone(s)
	}

	for _, mention := range m.Mentions {
		if mention.ID == d.State.User.ID {
			//log.Println(4)
			return true
		}
	}

	/*	if s == nil {
			switch g.DefaultMessageNotifications {
			case 0:
				log.Println(3)
				return true
			case 2:
			default:
			}
		} else {
			var notify = s.MessageNotifications
			if c := getChannelFromGuildSettings(m.ChannelID, s); c != nil {
				notify = c.MessageNotifications
			}

			switch notify {
			case 0: // all messages
				log.Println(3)
				return true
			case 2:
				return false
			default: // case 1 - mentions only
			}

			for _, mention := range m.Mentions {
				if mention.ID == d.State.User.ID {
					log.Println(4)
					return true
				}
			}
		}
	*/

	return false
}

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

func settingGuildAllowEveryone(s *discordgo.UserGuildSettings) bool {
	if s == nil {
		return true
	}

	return !s.SupressEveryone
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
