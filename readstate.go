package main

import (
	"log"
	"strings"
	"sync"

	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tview"
)

const readChannelColorPrefix = "[#808080::]"
const unreadChannelColorPrefix = "[::b]"

func messageAck(s *discordgo.Session, a *discordgo.MessageAck) {
	// Sets ReadState to the message you read
	for _, c := range d.State.ReadState {
		if c.ID == a.ChannelID && c.LastMessageID != 0 {
			c.LastMessageID = a.MessageID
		}
	}

	c, err := d.State.Channel(a.ChannelID)
	if err != nil {
		return
	}

	if c.GuildID == 0 {
		ackMeUI(c)
	} else {
		g, err := d.State.Guild(c.GuildID)
		if err != nil {
			return
		}

		checkGuild(g)
	}
}

// "[::b]actual string[::-]"
func stripFormat(a string) string {
	if len(a) <= 10 {
		return a
	}

	if strings.HasPrefix(a, readChannelColorPrefix) {
		a = a[len(readChannelColorPrefix):]
	}

	return strings.TrimSuffix(a, "[-::-]")
}

var (
	guildSettingsMuted   = map[int64]bool{}
	channelSettingsMuted = map[int64]bool{}
	settingsCacheMutex   = &sync.Mutex{}
)

func guildMuted(g *discordgo.Guild) bool {
	if g == nil {
		return false
	}

	settingsCacheMutex.Lock()
	defer settingsCacheMutex.Unlock()

	guMuted, ok := guildSettingsMuted[g.ID]
	if !ok {
		gs := getGuildFromSettings(g.ID)

		guMuted = settingGuildIsMuted(gs)
		guildSettingsMuted[g.ID] = guMuted
	}

	return guMuted
}

// true if channelID has unread msgs
func isUnread(ch *discordgo.Channel) bool {
	var gs *discordgo.UserGuildSettings

	settingsCacheMutex.Lock()

	chMuted, ok := channelSettingsMuted[ch.ID]
	if !ok {
		if gs == nil {
			gs = getGuildFromSettings(ch.GuildID)
		}

		cs := getChannelFromGuildSettings(ch.ID, gs)

		chMuted = settingChannelIsMuted(cs, gs)
		channelSettingsMuted[ch.ID] = chMuted
	}

	var guMuted = false

	if ch.GuildID != 0 {
		guMuted, ok = guildSettingsMuted[ch.GuildID]
		if !ok {
			if gs == nil {
				gs = getGuildFromSettings(ch.GuildID)
			}

			guMuted = settingGuildIsMuted(gs)
			guildSettingsMuted[ch.GuildID] = guMuted
		}
	}

	settingsCacheMutex.Unlock()

	if chMuted {
		return false
	}

	if ch.LastMessageID == 0 {
		return false
	}

	for _, c := range d.State.ReadState {
		if c.ID == ch.ID {
			return c.LastMessageID != ch.LastMessageID
		}
	}

	return false
}

func markUnread(m *discordgo.Message) {
	var unread bool

	c, err := d.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	if c.GuildID == 0 {
		// If the latest DM message is not the current message,
		// it's unread.
		for _, r := range d.State.ReadState {
			if r.ID == c.ID {
				unread = (m.ID != r.LastMessageID)
				break
			}
		}
	} else {
		// If neither the channel nor the guild is muted, it's
		// unread.
		gs := getGuildFromSettings(c.GuildID)
		chSettings := getChannelFromGuildSettings(c.ID, gs)

		var (
			chMuted = settingChannelIsMuted(chSettings, gs)
			guMuted = settingGuildIsMuted(gs)
		)

		unread = !(guMuted || chMuted)
	}

	if !unread {
		return
	}

	root := guildView.GetRoot()
	if root == nil {
		return
	}

	root.Walk(func(node, parent *tview.TreeNode) bool {
		if parent == nil {
			return true
		}

		switch reference := node.GetReference().(type) {
		case *discordgo.Guild:
			if reference.ID != m.GuildID {
				return false
			}
		case *discordgo.Channel:
			if reference.ID != m.ChannelID {
				return false
			}

			if reference.GuildID == 0 {
				node.SetText(unreadChannelColorPrefix + reference.Name + "[-::-]")
			} else {
				node.SetText(unreadChannelColorPrefix + "#" + reference.Name + "[-::-]")
			}

			return false
		default:
			return true
		}

		return true
	})

	app.Draw()
}

var lastAck string

func ackMe(chID, ID int64) {
	c, err := d.State.Channel(chID)
	if err != nil {
		return
	}

	if isUnread(c) {
		// triggers messageAck
		a, err := d.ChannelMessageAck(c.ID, ID, lastAck)

		if err != nil {
			log.Println(err)
			return
		}

		lastAck = a.Token
	}

	if c.GuildID == 0 {
		ackMeUI(c)
	} else {
		g, err := d.State.Guild(c.GuildID)
		if err != nil {
			return
		}

		checkGuild(g)
	}
}

func ackMeUI(ch *discordgo.Channel) {
	root := guildView.GetRoot()
	if root == nil {
		return
	}

	root.Walk(func(node, parent *tview.TreeNode) bool {
		if parent == nil {
			return true
		}

		switch reference := node.GetReference().(type) {
		case *discordgo.Guild:
			if reference.ID != ch.GuildID {
				return false
			}
		case *discordgo.Channel:
			if reference.ID != ch.ID {
				return false
			}

			var name = makeDMName(reference)
			if reference.GuildID != 0 {
				name = "#" + name
			}

			node.SetText(readChannelColorPrefix + name + "[-::-]")
		default:
			return true
		}

		return true
	})

	app.Draw()
}

func checkGuild(g *discordgo.Guild) {
	root := guildView.GetRoot()
	if root == nil {
		return
	}

	for _, n := range root.GetChildren() {
		gd, ok := n.GetReference().(*discordgo.Guild)
		if !ok {
			continue
		}

		if gd.ID != g.ID {
			continue
		}

		checkGuildNode(g, n)
		return
	}
}

func checkGuildNode(g *discordgo.Guild, n *tview.TreeNode) {
	var unreads = make([]*discordgo.Channel, 0, len(g.Channels))
	for _, c := range g.Channels {
		if isUnread(c) {
			unreads = append(unreads, c)
		}
	}

	if len(unreads) == 0 {
		n.SetText(readChannelColorPrefix + g.Name + "[-::-]")
	} else if !guildMuted(g) {
		n.SetText(unreadChannelColorPrefix + g.Name + "[-::-]")
	}

Main:
	for _, node := range n.GetChildren() {
		reference := node.GetReference()
		if reference == nil {
			continue
		}

		ch, ok := reference.(*discordgo.Channel)
		if !ok {
			continue
		}

		for _, u := range unreads {
			if u.ID == ch.ID {
				node.SetText(unreadChannelColorPrefix + "#" + ch.Name + "[-::-]")
				continue Main
			}
		}

		node.SetText(readChannelColorPrefix + "#" + ch.Name + "[-::-]")
	}

	app.Draw()
}

// TODO: Check if guild has unread channel
