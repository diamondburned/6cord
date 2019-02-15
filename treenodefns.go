package main

import "github.com/rivo/tview"

// CollapseAll collapses all tree nodes
func CollapseAll(gn *tview.TreeNode) {
	for _, c := range gn.GetChildren() {
		c.Collapse()
	}
}
