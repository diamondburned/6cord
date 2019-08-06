package main

import "github.com/diamondburned/tview/v2"

// CollapseAll collapses all tree nodes
func CollapseAll(gn *tview.TreeNode) {
	for _, c := range gn.GetChildren() {
		c.Collapse()
	}
}
