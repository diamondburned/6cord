package main

import "strings"

func changeSelfNick(text []string) {
	if Channel == nil {
		Message("You're not in a channel!")
		return
	}

	if len(text) < 2 {
		Message("Missing nickname argument!")
		return
	}

	nickname := strings.Join(text[1:], " ")

	ch, err := d.State.Channel(Channel.ID)
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
		return
	}

	Message("Changed successfully")

	go func() {
		i, u := us.GetUser(ch.GuildID, d.State.User.ID)
		if u != nil {
			us.Guilds[ch.GuildID][i].Nick = nickname
		}
	}()
}
