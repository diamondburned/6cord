package main

import "sync"

var (
	heatedChannels = struct {
		sync.RWMutex
		buffer map[int64]struct{}
	}{
		buffer: map[int64]struct{}{},
	}
)

func heatedChannelsAdd(channelID int64) {
	heatedChannels.Lock()
	defer heatedChannels.Unlock()

	heatedChannels.buffer[channelID] = struct{}{}
}

func heatedChannelsExists(channelID int64) bool {
	heatedChannels.RLock()

	if _, ok := heatedChannels.buffer[channelID]; ok {
		heatedChannels.RUnlock()
		heatedChannels.Lock()
		defer heatedChannels.Unlock()

		delete(heatedChannels.buffer, channelID)

		return true
	}

	heatedChannels.RUnlock()
	return false
}
