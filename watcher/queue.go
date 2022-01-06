package watcher

import (
	"crypto/md5"
	"fmt"

	"github.com/cian911/switchboard/event"
)

type Queue struct {
	queue map[string]event.Event
}

func New() *Queue {
	return &Queue{
		queue: make(map[string]event.Event),
	}
}

func (q *Queue) Add(hash string, ev event.Event) {
	q.queue[hash] = ev
}

func (q *Queue) Retrieve(hash string) event.Event {
	return q.queue[hash]
}

func (q *Queue) Remove(hash string) {
	delete(q.queue, hash)
}

func generateHash(ev event.Event) string {
	data := []byte(fmt.Sprintf("%s%s%s%s", ev.File, ev.Path, ev.Destination, ev.Ext))
	return fmt.Sprintf("%x", md5.Sum(data))
}
