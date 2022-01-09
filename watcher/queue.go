package watcher

import (
	"crypto/md5"
	"fmt"

	"github.com/cian911/switchboard/event"
)

// Q holds the Queue
type Q struct {
	Queue map[string]event.Event
}

// NewQueue create a new Q object
func NewQueue() *Q {
	return &Q{
		Queue: make(map[string]event.Event),
	}
}

// Add adds to the queue
func (q *Q) Add(ev event.Event) {
	q.Queue[Hash(ev)] = ev
}

// Retrieve get an item from the queue given a valid hash
func (q *Q) Retrieve(hash string) event.Event {
	return q.Queue[hash]
}

// Remove removes an item from the queue
func (q *Q) Remove(hash string) {
	delete(q.Queue, hash)
}

// Size returns the size of the queue
func (q *Q) Size() int {
	return len(q.Queue)
}

// Empty returns a bool indicating if the queue is empty or not
func (q *Q) Empty() bool {
	return len(q.Queue) == 0
}

// Hash returns a md5 hash composed of an event File, Path, and Ext
func Hash(ev event.Event) string {
	data := []byte(fmt.Sprintf("%s%s%s", ev.File, ev.Path, ev.Ext))
	return fmt.Sprintf("%x", md5.Sum(data))
}
