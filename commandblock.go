package main

import (
	"strconv"
)

func parseUserID(text []string) int64 {
	mention := text[1]

	mention = mention[2:]
	mention = mention[:len(mention)-1]

	id, err := strconv.ParseInt(mention, 10, 64)
	if err != nil {
		Message(err.Error())
		return 0
	}

	return id
}

func blockUser(text []string) {
	id := parseUserID(text)
	if id == 0 {
		return
	}

	if err := d.RelationshipUserBlock(id); err != nil {
		Warn(err.Error())
		return
	}
}

func unblockUser(text []string) {
	id := parseUserID(text)
	if id == 0 {
		return
	}

	if err := d.RelationshipDelete(id); err != nil {
		Warn(err.Error())
		return
	}
}
