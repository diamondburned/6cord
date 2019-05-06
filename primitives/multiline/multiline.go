package multiline

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Multiline is the primitive that draws w3m
type Multiline struct {
	x, y, width, height int

	focus  tview.Focusable
	focusB bool

	buffer  []rune
	current int // current line

	bg tcell.Style

	done func(key tcell.EventKey, content string)
}

// NewMultiline makes a new picture
func NewMultiline() (*Multiline, error) {
	p := &Multiline{
		buffer: []rune{},
	}

	p.focus = p
	p.SetBackgroundColor(tcell.ColorBlack)

	return p, nil
}

// GetRect returns the rectangle dimensions
func (m *Multiline) GetRect() (int, int, int, int) {
	return m.x, m.y, m.width, m.height
}

// SetRect sets the rectangle dimensions
func (m *Multiline) SetRect(x, y, width, height int) {
	m.x = x
	m.y = y
	m.width = width
	m.height = height
}

// InputHandler sets no input handler, satisfying Primitive
func (m *Multiline) InputHandler() func(event *tcell.EventKey, setFocus func(m tview.Primitive)) {
	return func(event *tcell.EventKey, _ func(m tview.Primitive)) {
		key := event.Key()

		if event.Modifiers() != 0 && m.done != nil {
			if key == tcell.KeyEnter {
				m.buffer = append(m.buffer, '\n')
			}

			return
		}

		if key == tcell.KeyDEL && len(m.buffer) > 0 {
			m.buffer = m.buffer[:len(m.buffer)-1]
			return
		}

		if key == tcell.KeyEscape || key == tcell.KeyEnter {
			if m.done != nil {
				m.done(*event, string(m.buffer))
			}

			if key == tcell.KeyEnter {
				m.buffer = append(m.buffer, '\n')
			}

			return
		}

		if r := event.Rune(); r != 0 {
			if r == 79 {
				m.buffer = append(m.buffer, '\n')
			} else {
				m.buffer = append(m.buffer, r)
			}

			return
		}
	}
}

// SetBackgroundColor sets the background color
func (m *Multiline) SetBackgroundColor(c tcell.Color) {
	m.bg = tcell.Style(0).Background(c)
}

// Focus does nothing, really.
func (m *Multiline) Focus(delegate func(tview.Primitive)) {
	m.focusB = true
}

// Blur also does nothing.
func (m *Multiline) Blur() {
	m.focusB = false
}

// HasFocus always returns false, as you can't focus on this.
func (m *Multiline) HasFocus() bool {
	return m.focusB
}

// GetFocusable does whatever the fuck I have no idea
func (m *Multiline) GetFocusable() tview.Focusable {
	return m.focus
}

func (m *Multiline) SetOnFocus(func()) {}
func (m *Multiline) SetOnBlur(func())  {}
