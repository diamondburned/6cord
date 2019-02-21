package main

import "strconv"

func messageGetTopID() int64 {
	if len(messageStore) < 1 {
		return 0
	}

	var msg string

	for i := 0; i < len(messageStore); i++ {
		if len(messageStore[i]) < 23 {
			continue
		}

		switch {
		case messageStore[i][1] != '[':
			continue
		case messageStore[i][2] != '"':
			continue
		}

		msg = messageStore[i]
		break
	}

	if msg == "" {
		return 0
	}

	var idRune string

	for i := 3; i < len(msg); i++ {
		if msg[i] == '"' {
			break
		}

		idRune += string(msg[i])
	}

	i, _ := strconv.ParseInt(idRune, 10, 64)
	return i
}
