package shortener

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	shortenerState = map[string]string{}
	shortenerMutex = &sync.RWMutex{}

	incr int
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

func ShortenURL(targetURL, customSlug, suffix string) string {
	shortenerMutex.Lock()
	defer shortenerMutex.Unlock()

	if customSlug != "" {
		customSlug += "-"
	}

	var id string

	fileshards := strings.Split(targetURL, "/")
	filename := fileshards[len(fileshards)-1]
	filename = filename[:max(len(filename)-len(suffix), 0)]
	filename = filename[:min(len(filename), 8)]

	for {
		slug := filename + "-" + increment()
		id = "/" + customSlug + slug + suffix
		if _, ok := shortenerState[id]; !ok {
			break
		}
	}

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
