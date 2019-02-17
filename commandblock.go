package main

import (
	"strconv"
)

func blockUser(text []string) {
	Message("Feature is blocked until Fishy updates his lib")
	return

	mention := text[1]

	mention = mention[2:]
	mention = mention[:len(mention)-1]

	id, err := strconv.ParseInt(mention, 10, 64)
	if err != nil {
		Message(err.Error())
		return
	}

	if err := d.RelationshipUserBlock(id); err != nil {
		Warn(err.Error())
		return
	}
}
