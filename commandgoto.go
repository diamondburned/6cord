package main

import (
	"strconv"
	"strings"

	"github.com/rivo/tview"
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

	id := parseChannelID(text[1])
	if id == 0 {
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

			refr, ok := node.GetReference().(int64)
			if !ok {
				return true
			}

			if id != refr {
				return true
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
