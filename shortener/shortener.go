package shortener

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

var (
	shortenerState = map[string]string{}
	shortenerMutex = &sync.RWMutex{}

	incr int

	// Enabled once StartHTTP is run
	Enabled = false
)

func Handler(w http.ResponseWriter, r *http.Request) {
	shortenerMutex.RLock()
	defer shortenerMutex.RUnlock()

	ou, ok := shortenerState[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, ou, http.StatusTemporaryRedirect)
}

func ShortenURL(targetURL string) string {
	if !Enabled {
		return targetURL
	}

	shortenerMutex.Lock()
	defer shortenerMutex.Unlock()

	u, err := url.Parse(targetURL)
	if err != nil {
		return targetURL
	}

	var id string
	var ext = GetExtension(u.Path)

	var fileshards []string
	if u.Path != "" {
		fileshards = strings.Split(u.Path, "/")
	} else {
		fileshards = []string{u.Host}
	}

	filename := fileshards[len(fileshards)-1]
	filename = filename[:max(len(filename)-len(ext), 0)]
	filename = filename[:min(len(filename), 8)]

	slug := filename + "-" + increment()
	id = "/" + slug + ext

	shortenerState[id] = targetURL
	return "http://" + URL + id
}

func min(i, j int) int {
	if i < j {
		return i
	}

	return j
}

func max(i, j int) int {
	if i > j {
		return i
	}

	return j
}
func increment() string {
	incr++
	return strconv.Itoa(incr)
}
