package main

import "sync"

var (
	heatedChannels = struct {
		sync.Mutex
		buffer map[int64]struct{}
	}{
		buffer: map[int64]struct{}{},
	}
)

// true == added, false == removed
func heatedChannelsToggle(channelID int64) bool {
	heatedChannels.Lock()
	defer heatedChannels.Unlock()

	if _, ok := heatedChannels.buffer[channelID]; ok {
		delete(heatedChannels.buffer, channelID)
		return false
	}

	heatedChannels.buffer[channelID] = struct{}{}
	return true
}

func heatedChannelsExists(channelID int64) bool {
	heatedChannels.Lock()
	defer heatedChannels.Unlock()

	if _, ok := heatedChannels.buffer[channelID]; ok {
		delete(heatedChannels.buffer, channelID)
		return true
	}

	return false
}
