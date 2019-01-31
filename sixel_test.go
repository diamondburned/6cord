package main

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/rthornton128/goncurses"
)

func TestSixel(t *testing.T) {
	// Get test
	wr, err := SixelFromURL(
		"https://i.kym-cdn.com/photos/images/original/001/244/891/d1f.png",
		20, 20,
	)

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("tfw you reatred %s for real though, fuck you in the ass all the way to hell\n", strings.TrimSpace(Printable(wr)))
}

func TestTUI(t *testing.T) {
	src, err := goncurses.Init()
	if err != nil {
		log.Fatal("init:", err)
	}
	defer goncurses.End()

	wr, err := SixelFromURL(
		"https://i.kym-cdn.com/photos/images/original/001/244/891/d1f.png",
		20, 20,
	)

	if err != nil {
		t.Error(err)
	}

	src.Print(Printable(wr))

	src.Refresh()

	// file, _ := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	// fmt.Fprintf(file, "\033[1mBOLD\033[0m")
}
