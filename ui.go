package main

var showChannels bool

func toggleChannels() {
	if showChannels {
		wrapFrame.SetBorder(true)
		appgrid.Clear()

		appgrid.SetColumns(cfg.Prop.SidebarWidth, 0)
		wrapFrame.SetBorders(0, 0, 0, 0, 1, 1)

		// We want the 2 items to span into 1 row, 2 columns.
		appgrid.AddItem(guildView, 0, 0, 1, 1, 0, 0, false)
		appgrid.AddItem(wrapFrame, 0, 1, 1, 1, 0, 0, false)

		// calculate minimum window width before hiding
		minWidth := cfg.Prop.SidebarWidth * 3

		// We also want to hide the side panel if the window is smaller than 200
		// chars or so.
		appgrid.AddItem(guildView, 0, 0, 0, 0, 0, minWidth, false)
		appgrid.AddItem(wrapFrame, 0, 1, 1, 2, 0, minWidth, false)

		app.SetFocus(guildView)
	} else {

		wrapFrame.SetBorder(false)
		appgrid.Clear()

		appgrid.SetColumns(0)
		wrapFrame.SetBorders(0, 0, 0, 0, 0, 0)

		appgrid.AddItem(wrapFrame, 0, 0, 1, 4, 0, 200, false)

		app.SetFocus(input)
	}

	showChannels = !showChannels
	app.Draw()
}
