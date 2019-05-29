package shortener

import "strings"

func GetExtension(name string) string {
	parts := strings.Split(name, "/")
	ss := strings.Split(parts[len(parts)-1], ".")

	if len(ss) < 2 {
		return ""
	}

	return "." + ss[len(ss)-1]
}
