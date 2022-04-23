package watcher

import (
	"testing"
	"time"

	"github.com/cian911/switchboard/event"
)

var (
	pollInterval = 1
)

func TestPoller(t *testing.T) {
	t.Run("It successfully notifies of a new event", func(t *testing.T) {
		pw := TestPathWatcher()
		pw.Poll(pollInterval)

		ev := eventSetup(t)
		pw.Queue.Add(*ev)

		if pw.Queue.Size() != 1 {
			t.Errorf("Queue size did not increase. want=%d, got=%d", 1, pw.Queue.Size())
		}
		<-time.After(3 * time.Second)

		if pw.Queue.Size() != 0 {
			t.Errorf("Queue size did not decrease. want=%d, got=%d", 0, pw.Queue.Size())
		}
	})

	t.Run("It successfully notifies on a new dir event", func(t *testing.T) {
		pw := TestPathWatcher()
		pw.Poll(pollInterval)

		ev := setupNewDirEvent(t)
		pw.Queue.Add(*ev)

		if pw.Queue.Size() != 1 {
			t.Errorf("Queue size did not increase. want=%d, got=%d", 1, pw.Queue.Size())
		}
		<-time.After(3 * time.Second)

		if pw.Queue.Size() != 0 {
			t.Errorf("Queue size did not decrease. want=%d, got=%d", 0, pw.Queue.Size())
		}
	})
}

func setupNewDirEvent(t *testing.T) *event.Event {
	path := t.TempDir()

	return &event.Event{
		Path:        path,
		Destination: t.TempDir(),
		Ext:         ".txt",
		Operation:   "CREATE",
		Timestamp:   time.Now(),
	}
}
