package main

import (
	"encoding/json"
	"fmt"

	"github.com/rumblefrog/discordgo"
)

const mentionsEndpoint = "https://discordapp.com/api/v6/users/@me/mentions?limit=%d&roles=true&everyone=true"

func getMentions() (ms []*discordgo.Message, err error) {
	resp, err := d.Request(
		"GET",
		fmt.Sprintf(mentionsEndpoint, 25),
		nil,
	)

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &ms)
	return
}
