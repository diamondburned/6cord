package main

var showChannels = cfg.ShowChannelsOnStartup

func toggleChannels() {
	if showChannels {
		wrapFrame.SetBorder(true)
		appflex.RemoveItem(guildView)
		appflex.RemoveItem(wrapFrame)

		wrapFrame.SetBorders(0, 0, 0, 0, 1, 1)

		appflex.AddItem(guildView, 0, 1, true)
		appflex.AddItem(wrapFrame, 0, 2, true)

		app.SetFocus(guildView)
	} else {
		wrapFrame.SetBorder(false)
		appflex.RemoveItem(guildView)

		wrapFrame.SetBorders(0, 0, 0, 0, 0, 0)

		app.SetFocus(input)
	}

	showChannels = !showChannels
}
