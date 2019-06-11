package main

func commandEditor(text []string) {
	if Channel == nil {
		Message("You're not in a channel!")
		return
	}

	b, err := summonEditor()
	if err != nil {
		Warn(err.Error())
		return
	}

	if _, err := d.ChannelMessageSend(Channel.ID, string(b)); err != nil {
		Warn(err.Error())
		return
	}
}
