package main

import (
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
	root := guildView.GetRoot()
	if root == nil {
		return
	}

	if vc == nil {
		return
	}

	root.Walk(func(node, parent *tview.TreeNode) bool {
		if parent == nil {
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

		if vc.ChannelID == 0 {
			if id == vc.UserID {
				refreshVoiceTreeNode(parent, vc.GuildID, vc.ChannelID)
				return false
			}
		}

		if id != vc.ChannelID {
			return true
		}

		// Checks should all pass to confirm this is
		// a voice channel

		refreshVoiceTreeNode(node, vc.GuildID, vc.ChannelID)
		return false

	})

	app.Draw()
}

func refreshVoiceTreeNode(node *tview.TreeNode, guildID, channelID int64) {
	node.ClearChildren()

	vcs := getVoiceChannel(guildID, channelID)
	for _, vc := range vcs {
		vcNode := generateVoiceNode(vc)
		if vcNode == nil {
			continue
		}

		node.AddChild(vcNode)
	}
}

func generateVoiceNode(vc *discordgo.VoiceState) *tview.TreeNode {
	var color = "-"
	if canIhearthem(vc) || vc.SelfDeaf || vc.Deaf {
		color = "d"
	}

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
