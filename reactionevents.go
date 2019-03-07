package main

import "github.com/rumblefrog/discordgo"

func reactionAdd(s *discordgo.Session, ra *discordgo.MessageReactionAdd) {
	if Channel == nil || ra.ChannelID != Channel.ID {
		return
	}

	m, err := d.State.Message(Channel.ID, ra.MessageID)
	if err != nil {
		return
	}

	d.State.Lock()

	for _, r := range m.Reactions {
		if !isSameEmoji(r, ra.MessageReaction) {
			continue
		}

		r.Count++

		if ra.UserID == d.State.User.ID {
			r.Me = true
		}

		goto Found
	}

	m.Reactions = append(
		m.Reactions,
		&discordgo.MessageReactions{
			Count: 1,
			Me:    ra.UserID == d.State.User.ID,
			Emoji: &ra.Emoji,
		},
	)

Found:
	d.State.Unlock()
	handleReactionEvent(m)
}

func reactionRemove(s *discordgo.Session, rm *discordgo.MessageReactionRemove) {
	if Channel == nil || rm.ChannelID != Channel.ID {
		return
	}

	m, err := d.State.Message(Channel.ID, rm.MessageID)
	if err != nil {
		return
	}

	d.State.Lock()

	for i, r := range m.Reactions {
		if !isSameEmoji(r, rm.MessageReaction) {
			continue
		}

		r.Count--

		if r.Count == 0 {
			m.Reactions = removeAllReactions(
				m.Reactions, i,
			)

			break
		}

		if rm.UserID == d.State.User.ID {
			r.Me = false
		}
	}

	d.State.Unlock()
	handleReactionEvent(m)
}

func reactionRemoveAll(s *discordgo.Session, rm *discordgo.MessageReactionRemoveAll) {
	if Channel == nil || rm.ChannelID != Channel.ID {
		return
	}

	m, err := d.State.Message(Channel.ID, rm.MessageID)
	if err != nil {
		return
	}

	d.State.Lock()

	m.Reactions = []*discordgo.MessageReactions{}

	d.State.Unlock()

	handleReactionEvent(m)
}
