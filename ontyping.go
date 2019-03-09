package main

import (
	"sync"
	"time"

	"github.com/diamondburned/discordgo"
)

// TypingUsers is a store for all typing users
type TypingUsers struct {
	sync.RWMutex
	Store []*discordgo.TypingStart
}

var typing = &TypingUsers{}

func onTyping(s *discordgo.Session, ts *discordgo.TypingStart) {
	if Channel == nil {
		return
	}

	if ts.ChannelID != Channel.ID {
		return
	}

	if ts.UserID == d.State.User.ID {
		return
	}

	typing.AddUser(ts)
}

func renderCallback(tu *TypingUsers) {
	var (
		mems []string
		text = cfg.Prop.DefaultStatus
	)

	if len(tu.Store) > 0 {
		for _, st := range tu.Store {
			if st.GuildID != 0 {
				_, user := us.GetUser(
					st.GuildID, st.UserID,
				)

				if user != nil {
					var name = user.Nick
					if name == "" {
						name = user.Name
					}

					mems = append(mems, name)
				} else {
					m, err := d.State.Member(st.GuildID, st.UserID)
					if err != nil {
						continue
					}

					var name = m.Nick
					if name == "" {
						name = m.User.Username
					}

					mems = append(mems, name)
				}
			} else {
				ch, err := d.State.Channel(Channel.ID)
				if err != nil {
					continue
				}

				for _, r := range ch.Recipients {
					if r.ID == st.UserID {
						mems = append(mems, r.Username)
					}
				}
			}
		}

		text = HumanizeStrings(mems)
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
	}

	input.SetPlaceholder(text)
}

// Reset resets the store
func (tu *TypingUsers) Reset() {
	tu.Lock()
	defer tu.Unlock()

	tu.Store = []*discordgo.TypingStart{}
	go renderCallback(tu)
}

// AddUser this function needs to run in a goroutine
func (tu *TypingUsers) AddUser(ts *discordgo.TypingStart) {
	tu.RLock()

	for _, t := range tu.Store {
		if t.UserID == ts.UserID {
			tu.RUnlock()
			return
		}
	}

	tu.RUnlock()

	tu.Lock()

	tu.Store = append(tu.Store, ts)

	// Might be overkill
	/*
		sort.Slice(tu.Store, func(i, j int) bool {
			return tu.Store[i].Time.UnixNano() <
				tu.Store[j].Time.UnixNano()
		})
	*/

	tu.Unlock()

	renderCallback(tu)

	time.Sleep(time.Second * 15)

	if tu.RemoveUser(ts) {
		renderCallback(tu)
	}
}

// RemoveUser removes a user from a store array
// true is returned when a user is found and removed
func (tu *TypingUsers) RemoveUser(ts *discordgo.TypingStart) bool {
	tu.Lock()
	defer tu.Unlock()

	for i, d := range tu.Store {
		if d.UserID == ts.UserID {
			tu.Store = append(
				tu.Store[:i],
				tu.Store[i+1:]...,
			)

			return true
		}
	}

	return false
}
