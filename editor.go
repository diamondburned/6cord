package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func summonEditor() (b []byte, err error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	f, err := ioutil.TempFile("", "6cord-editor-*.md")
	if err != nil {
		return
	}

	defer os.Remove(f.Name())

	for {
		cmd := exec.Command(editor, f.Name())
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		app.Suspend(func() error {
			fmt.Println("Opening", editor+"...")
			err = cmd.Run()

			return nil
		})

		if err != nil {
			return
		}

		b, err = ioutil.ReadAll(f)
		if err != nil {
			return
		}

		if len(b) > 2000 {
			Warn("Content too long! The limit is 2000 characters!")
		}

		break
	}

	return
}
