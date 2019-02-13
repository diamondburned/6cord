package w3m

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	w3mpath = GetExecPath()

	// ErrNotFound is returned when w3m can't be found
	ErrNotFound = errors.New("w3m not found")
)

// Arguments is the struct for w3m arguments
// All fields are required
type Arguments struct {
	Xoffset  int // default: 0
	Yoffset  int // default: 0
	Width    int
	Height   int
	Filename string
}

func Spawn(a Arguments) error {
	if w3mpath == "" {
		return ErrNotFound
	}

	var (
		cmd  = exec.Command(w3mpath)
		args = fmt.Sprintf(
			"0;1;%d;%d;%d;%d;;;;;%s\n3;\n4\n",
			a.Xoffset, a.Yoffset,
			a.Width, a.Height,
			a.Filename,
		)
	)

	reader := strings.NewReader(args)

	cmd.Stdin = reader

	cmd.Run()

	return nil
}
