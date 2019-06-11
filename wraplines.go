package main

import (
	"strings"

	"github.com/diamondburned/tview"
	"gitlab.com/diamondburned/6cord/md"
)

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

	_, _, col, _ := messagesView.GetInnerRect()

	for _, l := range lines {
		splwrap := strings.Split(
			md.Parse(tview.Escape(strings.Join(
				tview.WordWrap(l, min(col-5, 100)),
				"\n",
			))),
			"\n",
		)

		for _, s := range splwrap {
			spl = append(spl, cm+s+ce)
		}
	}

	return
}
