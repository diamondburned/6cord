package main

import (
	"strconv"
	"strings"

	"github.com/diamondburned/discordgo"
	"github.com/diamondburned/tview/v2"
)

func parseChannelID(input string) int64 {
	chID := strings.TrimSpace(input)

	chID = chID[2:]
	chID = chID[:len(chID)-1]

	id, err := strconv.ParseInt(chID, 10, 64)
	if err != nil {
		Message(err.Error())
		return 0
	}

	return id
}

func gotoChannel(text []string) {
	if len(text) != 2 {
		Message("No channels given!")
		return
	}

	var id int64

	switch {
	case strings.HasPrefix(text[1], "<#"):
		id = parseChannelID(text[1])
	case strings.HasPrefix(text[1], "<@"):
		ch, err := d.UserChannelCreate(parseUserMention(text[1]))
		if err != nil {
			Warn(err.Error())
			return
		}

		id = ch.ID
	}

	if id == 0 {
		Message("No channels given!")
		return
	}

	go func() {
		root := guildView.GetRoot()
		if root == nil {
			return
		}

		root.Walk(func(node, parent *tview.TreeNode) bool {
			if parent == nil {
				CollapseAll(node)
				return true
			}

			refr, ok := node.GetReference().(*discordgo.Channel)
			if !ok {
				return true
			}

			if id != refr.ID {
				return false
			}

			node.Expand()
			parent.Expand()
			guildView.SetCurrentNode(node)

			return false
		})
	}()

	resetInputBehavior()

	loadChannel(id)
}
