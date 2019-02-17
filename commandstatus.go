package main

import (
	"strings"

	"github.com/RumbleFrog/discordgo"
)

func setStatus(input []string) {
	if d.State.Settings == nil {
		Message("Settings are uninitialized")
		return
	}

	s := d.State.Settings.Status

	if len(input) < 2 {
		switch s {
		case discordgo.StatusOnline:
			Message("Status: Online")
		case discordgo.StatusIdle:
			Message("Status: Idle")
		case discordgo.StatusDoNotDisturb:
			Message("Status: Do not disturb")
		case discordgo.StatusInvisible:
			Message("Status: Invisible")
		default:
			Message(string(s))
		}

		return
	}

	switch strings.Join(input[1:], " ") {
	case string(discordgo.StatusOnline), "Online":
		s = discordgo.StatusOnline

	case string(discordgo.StatusIdle), "Idle",
		"Away", "away":
		s = discordgo.StatusIdle

	case string(discordgo.StatusDoNotDisturb),
		"do not disturb", "Do not disturb", "Do Not Disturb",
		"Busy", "busy":
		s = discordgo.StatusDoNotDisturb

	case string(discordgo.StatusInvisible), "invis", "Invisible":
		s = discordgo.StatusInvisible

	default:
		Message("Unknown status to set, check description")
		return
	}

	if _, err := d.UserUpdateStatus(s); err != nil {
		Warn(err.Error())
		return
	}

	Message("Set status to " + string(s))
}
