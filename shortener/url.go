package shortener

import "strings"

func GetExtension(name string) string {
	ss := strings.Split(name, ".")

	if len(ss) < 2 {
		return ""
	}

	return "." + ss[len(ss)-1]
}
