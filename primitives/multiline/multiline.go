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
func NewMultiline() *Multiline {
	p := &Multiline{
		Buffer: [][]rune{
			[]rune{},
		},
	}

	p.focus = p
	p.Style = tcell.Style(0).Background(tcell.ColorBlack)

	return p
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
		case tcell.KeyDEL: // backspace
			m.delRune(false)
			return

		case tcell.KeyDelete: // delete
			m.delRune(true)

			return

		case tcell.KeyEscape, tcell.KeyEnter:
			if m.done != nil {
				m.done(*event, strings.Join(m.getLines(), "\n"))
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

		case tcell.KeyUp, tcell.KeyDown:
			switch key {
			case tcell.KeyUp:
				if m.cursorY > 0 {
					m.cursorY--
				}

			case tcell.KeyDown:
				if m.cursorY < len(m.state)-1 {
					m.cursorY++
				}
			}

			if m.cursorX > len(m.state[m.cursorY]) {
				m.cursorX = len(m.state[m.cursorY])
			}

			return
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

func (m *Multiline) getLines() []string {
	var s = make([]string, 0, len(m.Buffer))
	for _, l := range m.Buffer {
		s = append(s, string(l))
	}

	return s
}

func (m *Multiline) newLine() {
	m.Buffer = append(m.Buffer, []rune{})
	m.cursorX = 0
	m.cursorY++
}

func (m *Multiline) addRune(r rune) {
	newBuf := m.Buffer[m.cursorY]

	if m.cursorX < len(newBuf) {
		newBuf = append(newBuf[:m.cursorX+1], newBuf[m.cursorX:]...)
		newBuf[m.cursorX] = r
	} else {
		newBuf = append(newBuf, r)
	}

	m.Buffer[m.cursorY] = newBuf

	m.cursorX++
}

func (m *Multiline) delRune(reverse bool) {
	if len(m.Buffer[m.cursorY]) > 0 {
		if reverse {
			if m.cursorX < len(m.Buffer[m.cursorY]) {
				m.Buffer[m.cursorY] = append(
					m.Buffer[m.cursorY][:m.cursorX], m.Buffer[m.cursorY][m.cursorX+1:]...,
				)
			}

			return
		}

		m.Buffer[m.cursorY] = append(
			m.Buffer[m.cursorY][:m.cursorX-1], m.Buffer[m.cursorY][m.cursorX:]...,
		)

		m.cursorX--
	} else if len(m.Buffer) > 1 {
		// Delete the empty new line
		m.cursorY--
		m.Buffer = m.Buffer[:m.cursorY+1]
		m.cursorX = len(m.Buffer[m.cursorY]) // wrap cursor to EOL
	}
}

func (m *Multiline) Insert(s string) {
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		for _, r := range []rune(lines[0]) {
			m.addRune(r)
		}
	}

	if len(lines) > 1 {
		for _, l := range lines[1:] {
			m.newLine()

			for _, r := range []rune(l) {
				m.addRune(r)
			}
		}
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
