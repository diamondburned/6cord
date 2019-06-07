package main

import (
	"sync"

	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tview"
)

// User is used for one user
type User struct {
	ID      int64
	Discrim string
	Name    string
	Nick    string
	Color   int
	Bot     bool
}

// UserStore stores multiple users
type UserStore struct {
	sync.RWMutex
	Guilds map[int64]UserStoreArray
}

// UserStoreArray is an array
type UserStoreArray []*User

var us = &UserStore{
	Guilds: map[int64]UserStoreArray{},
}

// Populated returns a bool on whether or not the array
// alraedy is populated
func (s *UserStore) Populated(guildID int64) bool {
	if s == nil {
		return false
	}

	if guildID == 0 {
		return true
	}

	return len(s.Guilds[guildID]) > 0
}

// InStore checks if a user is in the store
func (s *UserStore) InStore(guildID, id int64) bool {
	if s == nil {
		return false
	}

	if _, u := s.GetUser(guildID, id); u != nil {
		return true
	}

	return false
}

// DiscordThis interfaces with DiscordGo
func (s *UserStore) DiscordThis(m *discordgo.Message) (n string, c int) {
	n = "invalid user"
	c = defaultNameColor

	if m.Author == nil || s == nil {
		return
	}

	defer func() {
		if m.Author.Bot {
			n += " [#7289da][BOT[][-::-]"
		}
	}()

	if m.GuildID == 0 {
		channel, err := d.State.Channel(m.ChannelID)
		if err != nil {
			return
		}

		m.GuildID = channel.GuildID
	}

	_, user := s.GetUser(m.GuildID, m.Author.ID)
	if user != nil {
		n = user.Name
		c = user.Color

		if user.Nick != "" {
			n = user.Nick
		}

		return
	}

	nick, color := getUserData(m.Author, m.ChannelID)

	u := s.UpdateUser(
		m.GuildID,
		m.Author.ID,
		m.Author.Username,
		nick,
		m.Author.Discriminator,
		color,
	)

	n = u.Name
	c = u.Color

	if u.Nick != "" {
		n = u.Nick
	}

	return
}

// GetUser returns the index and user for that ID
func (s *UserStore) GetUser(guildID, id int64) (int, *User) {
	s.RLock()
	defer s.RUnlock()

	if v, ok := s.Guilds[guildID]; ok {
		for i, u := range v {
			if u.ID == id {
				return i, u
			}
		}
	}

	return 0, nil
}

// RemoveUser removes the user from the store
func (s *UserStore) RemoveUser(guildID, id int64) {
	var index int

	s.Lock()
	defer s.Unlock()

	if v, ok := s.Guilds[guildID]; ok {
		for i, u := range v {
			if u.ID == id {
				index = i
				goto Remove
			}
		}
	}

	return

Remove:
	var st = s.Guilds[guildID]

	st[len(st)-1], st[index] = st[index], st[len(st)-1]
	s.Guilds[guildID] = st[:len(st)-1]
}

// UpdateUser updates an user
func (s *UserStore) UpdateUser(guildID, id int64, name, nick, discrim string, color int) *User {
	if s == nil {
		return nil
	}

	if i, u := s.GetUser(guildID, id); u != nil {
		if name != "" {
			u.Name = tview.Escape(name)
		}

		if nick != "" {
			u.Nick = tview.Escape(nick)
		}

		if discrim != "" {
			u.Discrim = discrim
		}

		if color > 0 {
			u.Color = color
		}

		s.Lock()
		defer s.Unlock()

		s.Guilds[guildID][i] = u
		return u
	}

	s.Lock()
	defer s.Unlock()

	u := &User{
		ID:      id,
		Discrim: discrim,
		Name:    tview.Escape(name),
		Nick:    tview.Escape(nick),
		Color:   color,
	}

	s.Guilds[guildID] = append(s.Guilds[guildID], u)
	return u
}
