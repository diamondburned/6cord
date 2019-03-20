// +build !linux,!windows
// +build !unix,amd64

package main

import "os"

func syscallSilenceStderr(f *os.File) {
	// since we're unsure if the platform has dup2, we can just
	// silent all errors
	d.Debug = false
	d.LogLevel = 0
}
