package multiline

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Multiline is the primitive that draws w3m
type Multiline struct {
	x, y, width, height int

	focus  tview.Focusable
	focusB bool

	Placeholder      string
	PlaceholderColor tcell.Style

	cursorX, cursorY int

	Buffer  [][]rune
	current int      // current line
	state   []string // populated on last drawn

	isEnter bool

	Style tcell.Style

	done func(key tcell.EventKey, content string)
}

// NewMultiline makes a new picture
func NewMultiline() (*Multiline, error) {
	p := &Multiline{
		Buffer: [][]rune{},
	}

	p.focus = p
	p.Style = tcell.Style(0).Background(tcell.ColorBlack)

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
				m.newLine()
			}

			return
		}

		switch key {
		case tcell.KeyDEL:
			if len(m.Buffer) > 0 {
				m.addRune(-1)
				m.cursorX--
			}

			return

		case tcell.KeyEscape, tcell.KeyEnter:
			if m.done != nil {
				var s = make([]string, 0, len(m.Buffer))
				for _, l := range m.Buffer {
					s = append(s, string(l))
				}

				m.done(*event, strings.Join(s, "\n"))
			}

			if key == tcell.KeyEnter {
				m.newLine()
			}

			return

		case tcell.KeyLeft:
			if m.cursorX > 0 {
				m.cursorX--
			}

		case tcell.KeyRight:
			if m.cursorX < len(m.state[m.cursorY]) {
				m.cursorX++
			}
		}

		if r := event.Rune(); r != 0 {
			// Hack for Shift+Enter, which sends 'O' + 'M'
			if r == 'O' && event.Modifiers() == 4 {
				m.isEnter = true
			} else if r == 'M' && m.isEnter {
				m.newLine()
			} else {
				m.addRune(r)
			}

			return
		}
	}
}

func (m *Multiline) newLine() {
	m.Buffer = append(m.Buffer, []rune{})
	m.cursorX = 0
	m.cursorY++
}

func (m *Multiline) addRune(r rune) {
	m.cursorX++

	last := m.Buffer[m.cursorY][m.cursorX:]

	if r > 0 {
		m.Buffer[m.cursorY] = append(m.Buffer[m.cursorY][:m.cursorX], r)
		m.Buffer[m.cursorY] = append(m.Buffer[m.cursorY], last...)
	} else {
		m.Buffer[m.cursorY] = append(m.Buffer[m.cursorY][:m.cursorX], last...)
	}
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
