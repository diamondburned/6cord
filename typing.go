package main

import "time"

var (
	typingDelay = time.Duration(time.Second * 8)
	typingTimer = time.NewTimer(typingDelay)
)

func typingTrigger() {
	select {
	case <-typingTimer.C:
		if TriggerTyping {
			if ChannelID == 0 {
				return
			}

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
