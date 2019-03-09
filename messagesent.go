package main

import (
	"fmt"
	"sync"
)

var sentMessages = &sentMessageStruct{}

type sentMessageStruct struct {
	sync.RWMutex
	store []int64
}

func (sm *sentMessageStruct) In(is int64) bool {
	sm.RLock()
	defer sm.RUnlock()

	Message(fmt.Sprintln(sm.store))

	for _, id := range sm.store {
		if id == is {
			return true
		}
	}

	return false
}

func (sm *sentMessageStruct) Add(what int64) {
	sm.Lock()
	defer sm.Unlock()

	Message("Prepending")

	sm.store = append(
		[]int64{what},
		sm.store...,
	)
}
