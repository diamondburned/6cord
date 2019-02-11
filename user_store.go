package main

import (
	"sync"

	"github.com/RumbleFrog/discordgo"
)

// User is used for one user
type User struct {
	Name  string
	Nick  string
	Color int
}

// UserStore stores multiple users
type UserStore struct {
	Data map[int64]User
	Lock sync.RWMutex
}

var us = &UserStore{
	Data: make(map[int64]User),
}

// InStore checks if a user is in the store
func (s *UserStore) InStore(id int64) bool {
	if s == nil {
		return false
	}

	s.Lock.RLock()
	defer s.Lock.RUnlock()

	_, ok := s.Data[id]
	return ok
}

// AddUser adds an user into the store
func (s *UserStore) AddUser(id int64, name, nick string, color int) {
	if s.InStore(id) {
		return
	}

	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Data[id] = User{
		Name:  name,
		Nick:  nick,
		Color: color,
	}
}

// DiscordThis interfaces with DiscordGo
func (s *UserStore) DiscordThis(m *discordgo.Message) (n string, c int) {
	n = "invalid user"
	c = 16777215

	if m.Author == nil || s == nil {
		return
	}

	user := s.GetUser(m.Author.ID)
	if user != nil {
		n = user.Name
		c = user.Color

		if user.Nick != "" {
			n = user.Nick
		}

		return
	}

	nick, color := getUserData(m.Author, m.ChannelID)
	s.AddUser(
		m.Author.ID,
		m.Author.Username,
		nick,
		color,
	)

	n = m.Author.Username
	c = color

	if nick != "" {
		n = nick
	}

	return
}

// GetUser returns the user for that ID
func (s *UserStore) GetUser(id int64) *User {
	s.Lock.RLock()
	defer s.Lock.RUnlock()

	if u, ok := s.Data[id]; ok {
		return &u
	}

	return nil
}

// UpdateUser updates an user
func (s *UserStore) UpdateUser(id int64, name, nick string, color int) {
	if s == nil {
		return
	}

	s.Lock.Lock()
	defer s.Lock.Unlock()

	if u, ok := s.Data[id]; ok {
		switch {
		case name != "":
			u.Name = name
		case nick != "":
			u.Nick = nick
		case color > 0:
			u.Color = color
		}

		s.Data[id] = u
	}
}
