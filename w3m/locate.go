package w3m

import (
	"log"
	"os"
	"path/filepath"
)

// Thanks dylan
var (
	paths = []string{
		"/usr/local/lib/w3m/w3mi*",
		"/usr/local/libexec/w3m/w3mi*",
		"/usr/local/lib64/w3m/w3mi*",
		"/usr/local/libexec64/w3m/w3mi*",
		"/usr/lib/w3m/w3mimgdisplay",
		"/usr/libexec/w3m/w3mi*",
		"/usr/lib64/w3m/w3mimgdisplay",
		"/usr/libexec64/w3m/w3mi*",
	}
)

// GetExecPath finds w3mimgdisplay
func GetExecPath() string {
	// Todo: find a more performant way to do this
	for _, p := range paths {
		m, err := filepath.Glob(p)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, path := range m {
			info, err := os.Stat(path)
			if err != nil {
				log.Println(err)
				continue
			}

			if info.Mode()&0111 != 0 {
				return path
			}
		}
	}

	return ""
}
