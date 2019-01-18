package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
//StrongRegex  = regexp.MustCompile(`(\*\*|__) (?=\S) (.+?[*_]*) (?<=\S) \1`)
//ItalicsRegex = regexp.MustCompile(`(\*|_) (?=\S) (.+?) (?<=\S) \1`)
)

// ParseMessageContent parses bold and stuff into markup
func ParseMessageContent(m *discordgo.Message) string {
	return ""
}

func parseMD(content string) string {
	//content = StrongRegex.ReplaceAllString(content, "[::b]$0[::-]")
	//return ItalicsRegex.ReplaceAllString(content, `\e[3m$0\e[0m`)
	return ""
}
