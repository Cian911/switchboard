package watcher

import (
	"testing"
	"time"

	"github.com/cian911/switchboard/event"
)

func TestQueue(t *testing.T) {
	t.Run("It adds one event to the queue", func(t *testing.T) {
		q := setupQueue()
		ev := testEvent()

		q.Add(*ev)

		if q.Size() != 1 {
			t.Errorf("Could size did not increase as expected. want=%d, got=%d", 1, q.Size())
		}
	})

	t.Run("It updates the event in the queue", func(t *testing.T) {
		q := setupQueue()
		ev := testEvent()

		q.Add(*ev)
		q.Add(*ev)
		q.Add(*ev)

		if q.Size() != 1 {
			// Queue size should not increase
			t.Errorf("Could size did not increase as expected. want=%d, got=%d", 1, q.Size())
		}
	})
}

func setupQueue() *Q {
	return NewQueue()
}

func testEvent() *event.Event {
	return &event.Event{
		File:      "sample.txt",
		Path:      "/var/sample.txt",
		Ext:       ".txt",
		Timestamp: time.Now(),
	}
}
