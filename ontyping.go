package main

import (
	"sync"
	"time"

	"github.com/diamondburned/discordgo"
)

// TypingUsers is a store for all typing users
type TypingUsers struct {
	sync.RWMutex
	Store []*typingEvent
}

type typingEvent struct {
	*discordgo.TypingStart
	Meta *typingMeta
}

type typingMeta struct {
	Name string
	Time time.Time
}

type typingType int

const (
	typingStart typingType = iota
	typingStop
)

var (
	typing       = &TypingUsers{}
	updateTyping = make(chan struct{})
)

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

func getTypingMeta(typing *discordgo.TypingStart) *typingMeta {
	if typing.GuildID != 0 {
		_, user := us.GetUser(
			typing.GuildID, typing.UserID,
		)

		var name string

		if user != nil {
			if user.Nick == "" {
				name = user.Name
			} else {
				name = user.Nick
			}
		} else {
			m, err := d.State.Member(typing.GuildID, typing.UserID)
			if err != nil {
				return nil
			}

			if m.Nick == "" {
				name = m.User.Username
			} else {
				name = m.Nick
			}
		}

		return &typingMeta{
			Name: name,
			Time: time.Now(),
		}
	}

	ch, err := d.State.Channel(Channel.ID)
	if err != nil {
		return nil
	}

	for _, r := range ch.Recipients {
		if r.ID == typing.UserID {
			return &typingMeta{
				Name: r.Username,
				Time: time.Now(),
			}
		}
	}

	return nil
}

func renderCallback() {
	var (
		animation  uint
		laststring string
		tick       = time.Tick(time.Millisecond * 500)
	)

	for {
		var (
			mems []string
			anim string
			text = cfg.Prop.DefaultStatus
		)

		select { // 500ms or instant
		case <-updateTyping:
		case <-tick:
			if len(typing.Store) < 1 {
				animation = 0
			} else {
				animation++
				if animation > 5 {
					animation = 0
				}

				anim = getAnimation(animation)
			}
		}

		typing.RLock()

		for _, t := range typing.Store {
			if t.Meta == nil {
				t.Meta = getTypingMeta(t.TypingStart)
			}

			if t.Meta != nil {
				mems = append(mems, t.Meta.Name)
			}
		}

		typing.RUnlock()

		text = HumanizeStrings(mems)
		switch {
		case len(mems) < 1:
			text = "Send a message or input a command"
		case len(mems) > 3:
			text = "Several people are typing" + anim
		case len(mems) == 1:
			text += " is typing" + anim
		case len(mems) > 1:
			text += " are typing" + anim
		}

		if text != laststring {
			app.QueueUpdateDraw(func() {
				input.SetPlaceholder(text)
			})

			laststring = text
		}
	}
}

func getAnimation(i uint) string {
	switch i {
	case 0:
		return "   "
	case 1:
		return "·  "
	case 2:
		return "·· "
	case 3:
		return "···"
	case 4:
		return " ··"
	case 5:
		return "  ·"
	}

	return "   "
}

// Reset resets the store
func (tu *TypingUsers) Reset() {
	tu.Lock()
	defer tu.Unlock()

	tu.Store = []*typingEvent{}
	updateTyping <- struct{}{}
}

// AddUser this function needs to run in a goroutine
func (tu *TypingUsers) AddUser(ts *discordgo.TypingStart) {
	tu.Lock()
	for _, s := range tu.Store {
		if s.UserID == ts.UserID && s.Meta != nil {
			s.Meta.Time = time.Now()
			tu.Unlock()
			return
		}
	}
	tu.Unlock()

	ev := &typingEvent{
		TypingStart: ts,
		Meta:        getTypingMeta(ts),
	}

	tu.Lock()
	tu.Store = append(tu.Store, ev)
	tu.Unlock()

	updateTyping <- struct{}{}

	time.Sleep(time.Second * 10)

	// should always pass UNLESS there's another AddUser call bumping the
	// time up
	if ev.Meta.Time.Add(10 * time.Second).Before(time.Now()) {
		tu.RemoveUser(ts)
	}

	time.Sleep(time.Second * 1)
}

// RemoveUser removes a user from a store array
// true is returned when a user is found and removed
func (tu *TypingUsers) RemoveUser(ts *discordgo.TypingStart) bool {
	tu.Lock()
	defer tu.Unlock()

	defer func() {
		updateTyping <- struct{}{}
	}()

	if len(tu.Store) == 1 {
		tu.Reset()
		return true
	}

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
