package image

import (
	"fmt"
	"testing"
)

func TestWinch(t *testing.T) {
	if err := Listen(); err != nil {
		t.Fatal(t)
	}

	defer Close()

	if PixelH < 1 || PixelW < 1 {
		t.Fatal("failed")
	}

	fmt.Println(PixelH, PixelW)
}
