package watcher

import (
	"log"
	"time"
)

// Poll polls the queue for valid events given an interval (in seconds)
func (pw *PathWatcher) Poll(interval int) {
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		for {
			select {
			case <-ticker.C:
				log.Printf("Polling... - Queue Size: %d\n", pw.Queue.Size())

				for hsh, ev := range pw.Queue.Queue {
					timeDiff := ev.Timestamp.Sub(time.Now())
					if timeDiff < (time.Duration(-interval) * time.Second) {
						pw.Notify(ev.Path, ev.Operation)
						pw.Queue.Remove(hsh)
					}
				}
			}
		}
	}()
}
