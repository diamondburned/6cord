package main

import (
	"strings"

	"github.com/rivo/tview"
)

var (
//StrongRegex  = regexp.MustCompile(`(\*\*|__) (?=\S) (.+?[*_]*) (?<=\S) \1`)
//ItalicsRegex = regexp.MustCompile(`(\*|_) (?=\S) (.+?) (?<=\S) \1`)

)

// ParseMessageContent parses bold and stuff into markup
func ParseMessageContent(m string) string {
	var (
		output = []string{}
	)

	for _, w := range strings.Fields(m) {
		switch {
		case strings.HasPrefix(w, "<:"), strings.HasPrefix(w, "<a:"):
			parts := strings.Split(w, ":")
			if len(parts) < 3 || len(w) <= 18 {
				goto Skip
			}

			id := strings.TrimSuffix(parts[2], ">")

			var URL string
			if strings.HasPrefix(w, "<a:") {
				URL = "https://cdn.discordapp.com/emojis/" + id + ".gif"
			} else {
				URL = "https://cdn.discordapp.com/emojis/" + id + ".png"
			}

			output = append(output, URL)
			goto Skip
		}

		continue
	Skip:
		output = append(output, tview.Escape(w))
	}

	return strings.Join(output, " ")
}

func parseMD(content string) string {
	//content = StrongRegex.ReplaceAllString(content, "[::b]$0[::-]")
	//return ItalicsRegex.ReplaceAllString(content, `\e[3m$0\e[0m`)
	return ""
}
