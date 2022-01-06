package watcher

import "time"

func Poll(interval int) {
	ticker := time.NewTicker(interval * time.Second)
	for {
		select {
		case <-ticker.C:
		}
	}
}
