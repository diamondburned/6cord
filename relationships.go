package main

// Taken from
// https://gitlab.com/diamondburned/shittercord/blob/master/relationships.go

import (
	"sync"

	"github.com/rumblefrog/discordgo"
)

func relationshipAdd(s *discordgo.Session, ra *discordgo.RelationshipAdd) {
	for i, r := range d.State.Relationships {
		if r.ID == ra.ID {
			d.State.Relationships[i] = ra.Relationship
			return
		}
	}

	d.State.Relationships = append(
		d.State.Relationships,
		ra.Relationship,
	)
}

func relationshipRemove(s *discordgo.Session, rm *discordgo.RelationshipRemove) {
	rs := d.State.Relationships

	for i, r := range rs {
		if r.ID == rm.ID {
			rs = append(rs[:i], rs[i+1:]...)
			return
		}
	}

	d.State.Relationships = rs
}

// RStore contains Discord relationships
type RStore struct {
	Relationships []*discordgo.Relationship
	lock          sync.RWMutex
}

type Relationship int

const (
	RelationshipNone Relationship = iota

	RelationshipFriend                // friend
	RelationshipBlocked               // blocked
	RelationshipIncomingFriendRequest // incoming friend request
	RelationshipSentFriendRequest     // sent friend request
)

var (
	rstore = &RStore{}
)

// Check returns true if user is blocked
func (rs *RStore) Check(u *discordgo.User, relationship Relationship) bool {
	if !HideBlocked && relationship == RelationshipBlocked {
		return false
	}

	rs.lock.RLock()
	defer rs.lock.RUnlock()

	for _, r := range rs.Relationships {
		if r == nil {
			continue
		}

		if r.User == nil {
			continue
		}

		if r.Type == int(relationship) && r.User.ID == u.ID {
			return true
		}
	}

	return false
}

// Get gets the relationship of a user
func (rs *RStore) Get(u *discordgo.User) Relationship {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	for _, r := range rs.Relationships {
		if r.User.ID == u.ID {
			return parseInt(r.Type)
		}
	}

	return RelationshipNone
}

// Remove removes a relationship from the array
func (rs *RStore) Remove(r *discordgo.Relationship) {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	for i, rr := range rs.Relationships {
		if rr == r {
			// arr := *rs
			// arr[i] = arr[len(arr)-1]
			// *rs = arr[:len(arr)-1]

			rs.Relationships = append(rs.Relationships[:i], rs.Relationships[i+1:]...)
		}
	}
}

// Add adds a relationship to the store
func (rs *RStore) Add(r *discordgo.Relationship) {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.Relationships = append(rs.Relationships, r)
}

func parseInt(i int) Relationship {
	switch i {
	case int(RelationshipFriend):
		return RelationshipFriend
	case int(RelationshipBlocked):
		return RelationshipBlocked
	case int(RelationshipIncomingFriendRequest):
		return RelationshipIncomingFriendRequest
	case int(RelationshipSentFriendRequest):
		return RelationshipSentFriendRequest
	default:
		return RelationshipNone
	}
}
