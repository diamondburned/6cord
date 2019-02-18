package main

import "strings"

func changeSelfNick(text []string) {
	if len(text) < 2 {
		Message("Missing nickname argument!")
		return
	}

	nickname := strings.Join(text[1:], " ")

	ch, err := d.State.Channel(ChannelID)
	if err != nil {
		Warn(err.Error())
		return
	}

	if ch.GuildID == 0 {
		Message("You can't set a nickname in a DM")
		return
	}

	err = d.GuildMemberNicknameMe(
		ch.GuildID,
		nickname,
	)

	if err != nil {
		Message(err.Error())
	} else {
		Message("Changed successfully")
	}
}
