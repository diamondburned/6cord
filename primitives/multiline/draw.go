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

	// Conditions for drawing the placeholder
	if len(m.Buffer) == 1 && len(m.Buffer[0]) == 0 {
		m.state = strings.Split(m.Placeholder, "\n")
	} else {
		m.state = m.getLines()
	}

	for y := 0; y < m.height; y++ {
		if y < len(m.state) {
			runes := []rune(m.state[y])

			for i, r := range runes {
				s.SetContent(m.x+i, m.y+y, r, nil, m.Style)
			}

			for i := len(runes); i < m.width; i++ {
				s.SetContent(m.x+i, m.y+y, ' ', nil, m.Style)
			}
		} else {
			for x := 0; x < m.width; x++ {
				s.SetContent(m.x+x, m.y+y, ' ', nil, m.Style)
			}
		}
	}

	if m.focusB {
		if len(m.Buffer) > 0 {
			s.ShowCursor(
				m.x+m.cursorX,
				m.y+m.cursorY,
			)
		} else {
			s.ShowCursor(m.x, m.y)
		}
	}
}
