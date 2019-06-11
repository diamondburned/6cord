package main

import (
	"os"
	"os/exec"
)

func summonEditor() (err error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	f, err := TempFile("", "6cord-editor-")
	if err != nil {
		return err
	}

	defer os.Remove(f.Name())

	cmd := exec.Command(editor, f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	app.Suspend(func() { err = cmd.Run() })
	return
}
