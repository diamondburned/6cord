package main

import "time"

var (
	typingDelay = time.Duration(time.Second * 8)
	typingTimer = time.NewTimer(typingDelay)
)

func typingTrigger() {
	select {
	case <-typingTimer.C:
		if cfg.Prop.TriggerTyping {
			if Channel == nil {
				return
			}

			go d.ChannelTyping(Channel.ID)
		}
	default:
		return
	}

	typingTimer.Reset(typingDelay)
}
