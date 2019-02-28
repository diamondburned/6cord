package md

import (
	"strings"
)

func isFormatEnter(e bool, p string) string {
	if e {
		return "[::" + p + "]"
	}

	return "[::-]"
}

/*
	fixQuotes turns

	">string\nwhatever"

	into

	">string\n\nwhatever"
*/

func fixQuotes(s string) string {
	lines := strings.Split(s, "\n")

	for i := 1; i < len(lines); i++ {
		if (len(lines[i-1]) > 0 && len(lines[i]) > 0) &&
			lines[i-1][0] == '>' && lines[i][0] != '>' {

			lines[i] += "\n" + lines[i]
		}
	}

	return strings.Join(lines, "\n")
}
