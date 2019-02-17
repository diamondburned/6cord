package main

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/RumbleFrog/discordgo"
	"github.com/davecgh/go-spew/spew"
)

// TypingUsers is a store for all typing users
type TypingUsers struct {
	Store []TypingUser
	lock  sync.RWMutex
}

// TypingUser one user
type TypingUser struct {
	ID   int64
	Time time.Time
}

var typing = &TypingUsers{}

func onTyping(s *discordgo.Session, ts *discordgo.TypingStart) {
	log.Println(spew.Sdump(ts))

	if ts.ChannelID != ChannelID {
		return
	}

	log.Println(ts.UserID, ts.Timestamp)
	go typing.AddUser(ts.UserID, time.Now())
}

func renderCallback(tu *TypingUsers) {
	ch, err := d.State.Channel(ChannelID)
	if err != nil {
		Warn(err.Error())
		return
	}

	var mems []string

	for _, st := range tu.Store {
		m, err := d.State.Member(ch.GuildID, st.ID)
		if err != nil {
			Warn(err.Error())
			continue
		}

		if m.Nick != "" {
			mems = append(mems, m.Nick)
		} else {
			mems = append(mems, m.User.Username)
		}
	}

	text := HumanizeStrings(mems)
	switch {
	case len(mems) < 1:
		text = "Send a message or input a command"
	case len(mems) > 3:
		text = "Several people are typing···"
	case len(mems) == 1:
		text += " is typing···"
	case len(mems) > 1:
		text += " are typing···"
	}

	input.SetPlaceholder(text)
}

// Reset resets the store
func (tu *TypingUsers) Reset() {
	tu.lock.Lock()
	defer tu.lock.Unlock()

	tu.Store = []TypingUser{}
}

// AddUser this function needs to run in a goroutine
func (tu *TypingUsers) AddUser(id int64, t time.Time) {
	tu.lock.Lock()

	tu.Store = append(tu.Store, TypingUser{
		ID:   id,
		Time: t,
	})

	// Might be overkill
	sort.Slice(tu.Store, func(i, j int) bool {
		return tu.Store[i].Time.UnixNano() <
			tu.Store[j].Time.UnixNano()
	})

	tu.lock.Unlock()

	go renderCallback(tu)

	// 6 seconds according to djs code
	time.Sleep(time.Second * 6)

	tu.lock.Lock()
	defer tu.lock.Unlock()

	tu.RemoveUser(id)

	go renderCallback(tu)
}

// RemoveUser removes a user from a store array
func (tu *TypingUsers) RemoveUser(id int64) {
	var index int

	for i, d := range tu.Store {
		if d.ID == id {
			index = i
			goto Remove
		}
	}

	return

Remove:
	tu.Store = append(
		tu.Store[:index],
		tu.Store[index+1:]...,
	)
}
