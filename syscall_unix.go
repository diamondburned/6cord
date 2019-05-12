// +build unix,amd64

package main

import (
	"log"
	"os"
	"syscall"
)

func syscallSilenceStderr(f *os.File) {
	if err := syscall.Dup2(int(f.Fd()), 2, 0); err != nil {
		log.Println("Can't steal stderr, instabilities may occur")
	}
}