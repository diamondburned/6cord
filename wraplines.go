package main

import (
	"strings"

	"github.com/eidolon/wordwrap"
	"gitlab.com/diamondburned/6cord/md"
)

// WordWrapper makes a global wrapper for embed use
var WordWrapper = wordwrap.Wrapper(EmbedColLimit, false)

// 2nd arg ::-
// 3rd arg -::
func splitEmbedLine(e string, customMarkup ...string) (spl []string) {
	lines := strings.Split(e, "\n")

	// Todo: clean this up ETA never

	var (
		cm = ""
		ce = ""
	)

	if len(customMarkup) > 0 {
		cm = customMarkup[0]
		ce = "[::-]"
	}

	if len(customMarkup) > 1 {
		cm += customMarkup[1]
		ce += "[-::]"
	}

	for _, l := range lines {
		splwrap := strings.Split(md.Parse(WordWrapper(l)), "\n")

		for _, s := range splwrap {
			spl = append(spl, cm+s+ce)
		}
	}

	return
}
