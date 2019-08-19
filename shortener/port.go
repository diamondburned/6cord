package shortener

import (
	"net"
	"net/http"
	"strconv"
)

// URL is the current URL
var URL = ""

// GetOpenPort finds an open port
func GetOpenPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func StartHTTP(ip string) error {
	p, err := GetOpenPort()
	if err != nil {
		return err
	}

	URL = ip + ":" + strconv.Itoa(p)

	http.HandleFunc("/", Handler)
	go http.ListenAndServe(URL, nil)

	Enabled = true

	return nil
}
