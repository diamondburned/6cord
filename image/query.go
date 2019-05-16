package image

import (
	"strconv"
	"time"

	"github.com/mattn/go-tty"
)

// All these functions are taken from
// https://github.com/gizak/termui/pull/233/files#diff-61ca5d3d7b39f5b633e6774d6a31aea5R91
// and is modified promptly for this library.
//
// All credits reserved.

func queryTerm(qs string, t *tty.TTY) (ret [][]rune) {
	var b []rune

	ch := make(chan bool, 1)

	go func() {
	Main:
		for {
			r, err := t.ReadRune()
			if err != nil {
				return
			}
			// handle key event
			switch r {
			case 'c', 't':
				ret = append(ret, b)
				break Main
			case '?', ';':
				ret = append(ret, b)
				b = []rune{}
			default:
				b = append(b, r)
			}
		}

		ch <- true
	}()

	timer := time.NewTimer(100 * time.Microsecond)
	defer timer.Stop()

	select {
	case <-ch:
		defer close(ch)
	case <-timer.C:
	}

	return
}

func getTermSize(t *tty.TTY) (w, h int) {
	q := queryTerm("\033[14t", t)
	if len(q) != 3 {
		return
	}

	if yy, err := strconv.Atoi(string(q[1])); err == nil {
		if xx, err := strconv.Atoi(string(q[2])); err == nil {
			w = xx
			h = yy
		}
	}

	return
}
