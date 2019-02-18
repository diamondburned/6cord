package main

import (
	"github.com/rumblefrog/discordgo"
)

// User is used for one user
type User struct {
	ID      int64
	Discrim string
	Name    string
	Nick    string
	Color   int
}

// UserStore stores multiple users
type UserStore struct {
	GuildID int64
	Data    UserStoreArray
}

// UserStoreArray is an array
type UserStoreArray []User

var us = &UserStore{}

// Reset resets the store
func (s *UserStore) Reset(guildID int64) {
	if s == nil {
		return
	}

	s.GuildID = guildID
	s.Data = []User{}
}

// InStore checks if a user is in the store
func (s *UserStore) InStore(id int64) bool {
	if s == nil {
		return false
	}

	if _, u := s.GetUser(id); u != nil {
		return true
	}

	return false
}

// AddUser adds an user into the store
func (s *UserStore) AddUser(id int64, name, nick, discrim string, color int) {
	if s.InStore(id) {
		return
	}

	s.Data = append(s.Data, User{
		ID:      id,
		Discrim: discrim,
		Name:    name,
		Nick:    nick,
		Color:   color,
	})
}

// DiscordThis interfaces with DiscordGo
func (s *UserStore) DiscordThis(m *discordgo.Message) (n string, c int) {
	n = "invalid user"
	c = 16777215

	if m.Author == nil || s == nil {
		return
	}

	_, user := s.GetUser(m.Author.ID)
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
		m.Author.Discriminator,
		color,
	)

	n = m.Author.Username
	c = color

	if nick != "" {
		n = nick
	}

	return
}

// GetUser returns the index of the array and the user for that ID
func (s *UserStore) GetUser(id int64) (int, *User) {
	for i, u := range s.Data {
		if u.ID == id {
			return i, &u
		}
	}

	return -1, nil
}

// GetGuildID returns the guildID for the store
func (s *UserStore) GetGuildID() int64 {
	return s.GuildID
}

// UpdateUser updates an user
func (s *UserStore) UpdateUser(id int64, name, nick string, color int) {
	if s == nil {
		return
	}

	if i, u := s.GetUser(id); u != nil {
		switch {
		case name != "":
			u.Name = name
		case nick != "":
			u.Nick = nick
		case color > 0:
			u.Color = color
		}

		s.Data[i] = *u
	}
}
