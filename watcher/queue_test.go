package watcher

import (
	"reflect"
	"testing"
	"time"

	"github.com/cian911/switchboard/event"
)

var (
	gFile = "sample.txt"
	gPath = "/var/sample.txt"
	gExt  = ".txt"
)

func TestQueue(t *testing.T) {
	t.Run("It adds one event to the queue", func(t *testing.T) {
		q := setupQueue()
		ev := testEvent(gFile, gPath, gExt)

		q.Add(*ev)

		if q.Size() != 1 {
			t.Errorf("Could size did not increase as expected. want=%d, got=%d", 1, q.Size())
		}
	})

	t.Run("It updates the event in the queue", func(t *testing.T) {
		q := setupQueue()
		ev := testEvent(gFile, gPath, gExt)

		q.Add(*ev)
		q.Add(*ev)
		q.Add(*ev)

		if q.Size() != 1 {
			// Queue size should not increase
			t.Errorf("Could size did not increase as expected. want=%d, got=%d", 1, q.Size())
		}
	})

	t.Run("It gets an item from the queue", func(t *testing.T) {
		q := setupQueue()
		ev := testEvent(gFile, gPath, gExt)

		hash := Hash(*ev)
		q.Add(*ev)
		e := q.Retrieve(hash)

		if !reflect.DeepEqual(ev, &e) {
			t.Errorf("Events are not the same. want=%v, got=%v", ev, e)
		}
	})

	t.Run("It removes an item from the queue", func(t *testing.T) {
		q := setupQueue()
		ev := testEvent(gFile, gPath, gExt)

		hash := Hash(*ev)
		q.Add(*ev)
		q.Remove(hash)

		if q.Size() != 0 {
			t.Errorf("Could size did not increase as expected. want=%d, got=%d", 0, q.Size())
		}
	})

	t.Run("It returns a unique hash for a given event", func(t *testing.T) {
		ev1 := testEvent(gFile, gPath, gExt)
		ev2 := testEvent("sample2.txt", "/var/sample2.txt", ".txt")

		h1 := Hash(*ev1)
		h2 := Hash(*ev2)

		if h1 == h2 {
			t.Errorf("Hashes are the same when they shouldn't be. want=%s, got=%s", h1, h2)
		}
	})
}

func setupQueue() *Q {
	return NewQueue()
}

func testEvent(file, path, ext string) *event.Event {
	return &event.Event{
		File:      file,
		Path:      path,
		Ext:       ext,
		Timestamp: time.Now(),
	}

}
