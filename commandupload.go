package main

import (
	"os"
	"strings"
)

func uploadFile(args []string) {
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
		ChannelID,
		fileparts[len(fileparts)-1],
		file,
	)

	if err != nil {
		Warn(err.Error())
	}
}
