package main

import (
	"os"
	"strings"
)

func uploadFile(args []string) {
	if Channel == nil {
		Message("You're not in a channel!")
		return
	}

	if len(args) < 2 {
		Message("Missing file path!")
		return
	}

	file, err := os.Open(strings.Join(args[1:], " "))
	if err != nil {
		Warn(err.Error())
		return
	}

	fileparts := strings.Split(file.Name(), "/")

	_, err = d.ChannelFileSend(
		Channel.ID,
		fileparts[len(fileparts)-1],
		file,
	)

	if err != nil {
		Warn(err.Error())
	}
}
