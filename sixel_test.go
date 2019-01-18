package main

import (
	"fmt"
	"strings"
	"testing"
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
