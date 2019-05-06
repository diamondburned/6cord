package multiline

import (
	"strings"

	"github.com/gdamore/tcell"
)

var cursorColor = tcell.StyleDefault.Reverse(true)

// Draw draws the w3m image constantly
func (m *Multiline) Draw(s tcell.Screen) {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	/*
		lines := strings.Split(wrap.Wrap(
			m.buffer.String(), m.width,
		), "\n")
	*/

	lines := strings.Split(string(m.buffer), "\n")

	for y := 0; y < m.height; y++ {
		if len(m.buffer) != 0 && y < len(lines) {
			runes := []rune(lines[y])

			for i, r := range runes {
				s.SetContent(m.x+i, m.y+y, r, nil, m.bg)
			}

			for i := len(runes); i < m.width; i++ {
				s.SetContent(m.x+i, m.y+y, ' ', nil, m.bg)
			}
		} else {
			for x := 0; x < m.width; x++ {
				s.SetContent(m.x+x, m.y+y, ' ', nil, m.bg)
			}
		}
	}

	if m.focusB {
		s.ShowCursor(len([]rune(lines[len(lines)-1])), len(lines)-1)
	}
}
