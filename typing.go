package main

import "time"

var (
	typingDelay = time.Duration(time.Second * 2)
	typingTimer = time.NewTimer(typingDelay)
)

func typingTrigger() {
	select {
	case <-typingTimer.C:
		if TriggerTyping {
			go func() {
				err := d.ChannelTyping(ChannelID)
				if err != nil {
					Message(err.Error())
				}
			}()
		}
	default:
		return
	}

	typingTimer.Reset(typingDelay)
}
