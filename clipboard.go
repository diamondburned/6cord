package main

import (
	"net/http"
)

// IsFile checks if f clipboard is plain text
func IsFile(b []byte) bool {
	return http.DetectContentType(b) != "text/plain; charset=utf-8"
}
