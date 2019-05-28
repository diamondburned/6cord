package main

import (
	"time"

	"github.com/diamondburned/discordgo"
)

// TypingUsers is a store for all typing users
type TypingUsers struct {
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

	go typing.AddUser(ts)
}

func getTypingMeta(typing *discordgo.TypingStart) *typingMeta {
	if typing.GuildID != 0 {
		_, user := us.GetUser(
			typing.GuildID, typing.UserID,
		)

		m := &typingMeta{
			Time: time.Now(),
		}

		if user == nil {
			m, err := d.State.Member(typing.GuildID, typing.UserID)
			if err != nil {
				return nil
			}

			user = us.UpdateUser(
				m.GuildID,
				m.User.ID,
				m.User.Username,
				m.Nick,
				m.User.Discriminator,
				getUserColor(m.GuildID, m.Roles),
			)
		}

		if user.Nick == "" {
			m.Name = user.Name
		} else {
			m.Name = user.Nick
		}

		return m
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
		animation uint
		tick      = time.Tick(time.Second)
		dots      = [6]string{
			"   ", "·  ", "·· ",
			"···", " ··", "  ·",
		}
	)

	for {
		var mems = make([]string, 0, len(typing.Store))

		select { // 500ms or instant
		case <-updateTyping:
		case <-tick:
		}

		if len(typing.Store) < 1 {
			animation = 0
		} else {
			animation++
			if animation > 5 {
				animation = 0
			}
		}

		for _, t := range typing.Store {
			if t.Meta != nil {
				mems = append(mems, t.Meta.Name)
			}
		}

		text := cfg.Prop.DefaultStatus

		switch {
		case len(mems) > 3:
			text = "Several people are typing" + dots[animation]
		case len(mems) == 1:
			text = HumanizeStrings(mems) + " is typing" + dots[animation]
		case len(mems) > 1:
			text = HumanizeStrings(mems) + " are typing" + dots[animation]
		}

		if text != input.GetPlaceholder() && !messagesView.HasFocus() {
			input.SetPlaceholder(text)
			app.Draw()
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
	tu.Store = []*typingEvent{}
	updateTyping <- struct{}{}
}

// AddUser this function needs to run in a goroutine
func (tu *TypingUsers) AddUser(ts *discordgo.TypingStart) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	for _, s := range tu.Store {
		if s.UserID == ts.UserID && s.Meta != nil {
			s.Meta.Time = time.Now()
			return
		}
	}

	ev := &typingEvent{
		TypingStart: ts,
		Meta:        getTypingMeta(ts),
	}

	tu.Store = append(tu.Store, ev)

	updateTyping <- struct{}{}

	time.Sleep(time.Second * 10)

	// should always pass UNLESS there's another AddUser call bumping the
	// time up
	for {
		t := ev.Meta.Time
		if t.Add(10 * time.Second).Before(time.Now()) {
			tu.RemoveUser(ts)
			break
		}

		time.Sleep(time.Second * 1)
	}
}

// RemoveUser removes a user from a store array
// true is returned when a user is found and removed
func (tu *TypingUsers) RemoveUser(ts *discordgo.TypingStart) bool {
	for i, d := range tu.Store {
		if d.UserID == ts.UserID {
			tu.Store = append(
				tu.Store[:i],
				tu.Store[i+1:]...,
			)

			updateTyping <- struct{}{}
			return true
		}
	}

	return false
}
