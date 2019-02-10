package main

import "sync"

var (
	lastAuthor      int64
	lastAuthorMutex = &sync.RWMutex{}
)

func getLastAuthor() int64 {
	lastAuthorMutex.RLock()
	defer lastAuthorMutex.RUnlock()

	return lastAuthor
}

func setLastAuthor(i64 int64) {
	lastAuthorMutex.Lock()
	defer lastAuthorMutex.Unlock()

	lastAuthor = i64
}
