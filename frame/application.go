package frame

import "time"

func Idle(dur time.Duration, cb func() error) {
	for {
		select {
		case <-time.After(dur):
			if cb != nil {
				if cb() != nil {
					break
				}
			}
		}
	}
}
