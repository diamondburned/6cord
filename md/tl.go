package md

import (
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func isFormatEnter(e bool, p string) string {
	if e {
		return "[::" + p + "]"
	}

	return "[::-]"
}

// fixQuotes turns
// ">string\nwhatever"
// into
// ">string\n\nwhatever"
func fixQuotes(s string) string {
	lines := strings.Split(s, "\n")
	if len(lines) < 1 {
		return s
	}

	rebuild := []string{lines[0]}

	for i := 1; i < len(lines); i++ {
		if len(lines[i]) < 1 {
			goto Done
		}

		if lines[i][0] != '>' {
			if len(lines[i-1]) < 1 {
				goto Done
			}

			if lines[i-1][0] == '>' {
				rebuild = append(rebuild, "")
				goto Done
			}
		}

	Done:
		rebuild = append(rebuild, lines[i])
	}

	log.Println(spew.Sdump(rebuild))

	return strings.Join(rebuild, "\n")
}
