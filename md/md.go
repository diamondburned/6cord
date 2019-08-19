package md

import (
	"regexp"
	"strings"

	"github.com/diamondburned/tview/v2"
	"gitlab.com/diamondburned/6cord/shortener"
)

const mdRegex = `(?m)(?:^\x60\x60\x60 *(\w*)\n([\s\S]*?)\n\x60\x60\x60$)|((?:(?:^|\n)>\s+.*)+)|(?:(?:^|\n)(?:[>*+-]|\d+\.)\s+.*)+|(?:\x60([^\x60].*?)\x60)|(__|\*\*\*|\*\*|[_*]|~~|\|\|)|(https?:\/\S+(?:\.|:)\S+)`

var r1 = regexp.MustCompile(mdRegex)

type match struct {
	from, to int
	str      string
}

type mdState struct {
	strings.Builder
	matches [][]match
	last    int
	chunk   string
	prev    string
	context []string
}

func (s *mdState) tag(token string) string {
	var tags [2]string

	switch token {
	case "*":
		tags[0] = "[::i]"
		tags[1] = "[::-]"
	case "_":
		tags[0] = "[::i]"
		tags[1] = "[::-]"
	case "**":
		tags[0] = "[::b]"
		tags[1] = "[::-]"
	case "__":
		tags[0] = "[::u]"
		tags[1] = "[::-]"
	case "***":
		tags[0] = "[::ib]"
		tags[1] = "[::-]"
	case "~~":
		tags[0] = "[::s]"
		tags[1] = "[::-]"
	case "||":
		tags[0] = "[#777777]"
		tags[1] = "[-]"
	default:
		return token
	}

	var index = -1
	for i, t := range s.context {
		if t == token {
			index = i
			break
		}
	}

	if index >= 0 { // len(context) > 0 always
		s.context = append(s.context[:index], s.context[index+1:]...)
		return tags[1]
	} else {
		s.context = append(s.context, token)
		return tags[0]
	}
}

func (s mdState) getLastIndex(currentIndex int) int {
	if currentIndex >= len(s.matches) {
		return 0
	}

	return s.matches[currentIndex][0].to
}

func Parse(md string) string {
	var s mdState
	s.matches = submatch(r1, md)

	for i := 0; i < len(s.matches); i++ {
		s.prev = md[s.last:s.matches[i][0].from]
		s.last = s.getLastIndex(i)
		s.chunk = "" // reset chunk

		switch {
		case strings.Count(s.prev, "\\")%2 != 0:
			// escaped, print raw
			s.chunk = tview.Escape(s.matches[i][0].str)
		case s.matches[i][2].str != "":
			// codeblock
			s.chunk = RenderCodeBlock(
				tview.Escape(s.matches[i][1].str),
				tview.Escape(s.matches[i][2].str),
			)
		case s.matches[i][3].str != "":
			// blockquotes, greentext
			s.chunk = "\n[#789922]" +
				tview.Escape(strings.TrimPrefix(s.matches[i][3].str, "\n")) +
				"[-]"
		case s.matches[i][4].str != "":
			// inline code
			s.chunk = "[:#4f4f4f:]" + tview.Escape(s.matches[i][4].str) + "[:-:]"
		case s.matches[i][5].str != "":
			// inline stuff
			s.chunk = s.tag(s.matches[i][5].str)
		case s.matches[i][6].str != "":
			s.chunk = shortener.ShortenURL(s.matches[i][6].str)
		default:
			s.chunk = tview.Escape(s.matches[i][0].str)
		}

		s.WriteString(tview.Escape(s.prev))
		s.WriteString(s.chunk)
	}

	s.WriteString(md[s.last:])

	for len(s.context) > 0 {
		s.WriteString(s.tag(s.context[len(s.context)-1]))
	}

	return strings.TrimSpace(s.String())
}

func submatch(r *regexp.Regexp, s string) [][]match {
	found := r.FindAllStringSubmatchIndex(s, -1)
	indices := make([][]match, len(found))

	var m = match{-1, -1, ""}

	for i := range found {
		indices[i] = make([]match, len(found[i])/2)

		for a, b := range found[i] {
			if a%2 == 0 { // first pair
				m.from = b
			} else {
				m.to = b

				if m.from >= 0 && m.to >= 0 {
					m.str = s[m.from:m.to]
				} else {
					m.str = ""
				}

				indices[i][a/2] = m
			}
		}
	}

	return indices
}
