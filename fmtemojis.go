package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// EmojiRegex to get emoji IDs
	// thanks ym
	EmojiRegex = regexp.MustCompile(`<(.*?):(.*?):(\d+)>`)
)

// returns map[ID][]{name, url}
func parseEmojis(content string) (fmtted string, emojiMap map[string][]string) {
	emojiMap = make(map[string][]string)
	fmtted = content

	emojiIDs := EmojiRegex.FindAllStringSubmatch(content, -1)
	for _, nameandID := range emojiIDs {
		if len(nameandID) < 4 {
			continue
		}

		if _, ok := emojiMap[nameandID[3]]; !ok {
			var format = "png"
			if nameandID[1] != "" {
				format = "gif"
			}

			fmtted = strings.Replace(
				fmtted,
				strings.TrimSpace(nameandID[0]),
				"["+nameandID[2]+"[]",
				-1,
			)

			if ShowEmojiURLs {
				emojiMap[nameandID[3]] = []string{
					nameandID[2],
					fmt.Sprintf(
						`https://cdn.discordapp.com/emojis/%s.%s`,
						nameandID[3], format,
					),
				}
			}
		}
	}

	return
}
