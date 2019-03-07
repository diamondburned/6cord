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
			if Channel == nil {
				return
			}

			go func() {
				err := d.ChannelTyping(Channel.ID)
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
