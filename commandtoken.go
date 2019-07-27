package main

import "github.com/atotto/clipboard"

func commandCopyToken(text []string) {
	if err := clipboard.WriteAll(d.Token); err != nil {
		Warn(err.Error())
		return
	}

	Message("Token copied to clipboard.")
}
