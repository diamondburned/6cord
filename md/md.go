package md

import "github.com/rivo/tview"

// Parse parses md into tview string
func Parse(s string) string {
	s = tview.Escape(s)

	// this should have higher prio than italic
	s = UnderlineRegex.ReplaceAllString(s, "[::u]$1[::-]")

	// Bold for super bold (reverse bg/fg + bold)
	s = BoldRegex.ReplaceAllString(s, "[::br]$1[::-]")

	// tview doesn't have italics
	// We're treating them like bold
	for _, r := range ItalicRegexes {
		s = r.ReplaceAllString(s, "[::b]$1[::-]")
	}

	// Dim
	s = SpoilerRegex.ReplaceAllString(s, "[::d]$1[::-]")
	s = StrikethroughRegex.ReplaceAllString(s, "[::d]$1[::-]")

	return s
}
