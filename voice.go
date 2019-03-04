package main

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/rumblefrog/discordgo"
)

func voiceStateUpdate(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	refreshVoiceStates(vsu.VoiceState)
}

func getVoiceChannel(guildID, channelID int64) (vcs []*discordgo.VoiceState) {
	g, err := d.State.Guild(guildID)
	if err != nil {
		return
	}

	for _, vc := range g.VoiceStates {
		if vc.ChannelID == channelID {
			vcs = append(vcs, vc)
		}
	}

	return
}

func canIhearthem(vc *discordgo.VoiceState) bool {
	return !(vc.SelfMute || vc.Mute || vc.Suppress)
}

func refreshVoiceStates(vc *discordgo.VoiceState) {
	defer func() {
		if r := recover(); r != nil {
			Warn(fmt.Sprintln(r))
			return
		}
	}()

	root := guildView.GetRoot()
	if root == nil {
		return
	}

	if vc == nil {
		return
	}

	root.Walk(func(node, parent *tview.TreeNode) bool {
		if parent == nil || node == nil {
			return true
		}

		reference := node.GetReference()
		if reference == nil {
			return true
		}

		id, ok := reference.(int64)
		if !ok {
			return true
		}

		// user left voice chat
		if vc.ChannelID == 0 {
			// checks for ID should match the userID instead,
			// as the user left voice chat
			if id == vc.UserID {
				// parent node at this point should be the voice
				// channel
				var nodes []*tview.TreeNode
				for _, ch := range parent.GetChildren() {
					if node != ch {
						// we add everything except for the user
						// that left by adding everything else back
						nodes = append(nodes, ch)
					}
				}

				app.QueueUpdateDraw(func() {
					parent.SetChildren(nodes)
				})

				return false
			}

			return true
		}

		// user joined a voice channel

		if id != vc.ChannelID {
			return true
		}

		// checks should all pass to confirm this is
		// the right voice channel

		refreshVoiceTreeNode(node, vc.GuildID, vc.ChannelID)
		return false
	})
}

func refreshVoiceTreeNode(node *tview.TreeNode, guildID, channelID int64) {
	var (
		nodes []*tview.TreeNode
		vcs   = getVoiceChannel(guildID, channelID)
	)

	for _, vc := range vcs {
		vcNode := generateVoiceNode(vc)
		if vcNode == nil {
			continue
		}

		nodes = append(nodes, vcNode)
	}

	app.QueueUpdateDraw(func() {
		node.SetChildren(nodes)
	})
}

func generateVoiceNode(vc *discordgo.VoiceState) *tview.TreeNode {
	var color = "d"

	// Reserved for onSpeak
	//if !canIhearthem(vc) || vc.SelfDeaf || vc.Deaf {
	//	color = "d"
	//}

	u, err := d.State.Member(vc.GuildID, vc.UserID)
	if err != nil {
		return nil
	}

	if u.User == nil {
		return nil
	}

	var name = u.User.Username
	if u.Nick != "" {
		name = u.Nick
	}

	var suffix string

	if vc.SelfMute || vc.Mute {
		suffix += " [gray][M[][-]"
	}

	if vc.SelfDeaf || vc.Deaf {
		suffix += " [gray][D[][-]"
	}

	if vc.Suppress {
		suffix += " [red][Suppressed[][-]"
	}

	vcNode := tview.NewTreeNode(
		"[::" + color + "]" + name + "[::-]" + suffix,
	)

	vcNode.SetSelectable(false)
	vcNode.SetReference(vc.UserID)

	return vcNode
}
