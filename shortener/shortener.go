package shortener

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	shortenerState = map[string]string{}
	shortenerMutex = &sync.RWMutex{}
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

func ShortenURL(url, customSlug, suffix string) string {
	shortenerMutex.Lock()
	defer shortenerMutex.Unlock()

	var id string

	if customSlug != "" {
		customSlug += "-"
	}

	for {
		id = "/" + customSlug + getTime() + suffix
		if _, ok := shortenerState[id]; !ok {
			break
		}
	}

	shortenerState[id] = url
	return "http://" + URL + id
}

func getTime() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
