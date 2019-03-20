// +build linux

package main

import (
	"log"
	"os"
	"syscall"
)

func syscallSilenceStderr(f *os.File) {
	if err := syscall.Dup3(int(f.Fd()), 2, 0); err != nil {
		log.Println("Can't steal stderr, instabilities may occur")
	}
}
