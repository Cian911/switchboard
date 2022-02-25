//go:build linux
// +build linux

package watcher

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cian911/switchboard/event"
	"github.com/cian911/switchboard/utils"
	"github.com/fsnotify/fsnotify"
)

// Monitor for IN_CLOSE_WRITE events on these file exts
// A create event should immediatly follow
var specialWatchedFilExts = map[string]int{
	".part": 1,
}

// Producer interface for the watcher
// Must implement Register(), Unregister(), and Observe(), and notify()
type Producer interface {
	// Register a consumer to the producer
	Register(consumer *Consumer)
	// Unregister a consumer from the producer
	Unregister(consumer *Consumer)
	// Notify consumers of an event
	Notify(path, event string)
	// Observe the producer
	Observe(pollInterval int)
}

// Consumer interface
// Must implement Receive(), and Process() methods
type Consumer interface {
	// Receive an event from the producer
	Receive(path, event string)
	// Process an event
	Process(e *event.Event)
	// Process a dir event
	ProcessDirEvent(e *event.Event)
}

// PathWatcher is a producer that watches a path for events
type PathWatcher struct {
	// List of consumers
	Consumers []*Consumer
	// Queue
	Queue *Q
	// Watcher instance
	Watcher fsnotify.Watcher
	// Path to watch
	Path string
}

// PathConsumer is a consumer that consumes events from a path
// and moves them to a destination
type PathConsumer struct {
	// Path to watch
	Path string
	// Destination to move files to
	Destination string
	// File extenstion
	Ext string
}

// Receive takes a path and an event operation, determines its validity
// and passes it to be processed it if valid
func (pc *PathConsumer) Receive(path, ev string) {
	e := &event.Event{
		File:        filepath.Base(path),
		Path:        path,
		Destination: pc.Destination,
		Ext:         utils.ExtractFileExt(path),
		Timestamp:   time.Now(),
		Operation:   ev,
	}

	log.Printf("EVENT_EXT: %s - %s - %s", e.Ext, pc.Ext, pc.Path)

	if !e.IsNewDirEvent() && e.Ext != pc.Ext && filepath.Dir(path) != pc.Path {
		log.Printf("Not processing event - %v - %v", e, pc)
		// Do not process event for consumers not watching file
		return
	}

	if e.IsNewDirEvent() {
		pc.ProcessDirEvent(e)
	} else if e.IsValidEvent(pc.Ext) {
		pc.Process(e)
	}
}

// Process takes an event and moves it to the destination
func (pc *PathConsumer) Process(e *event.Event) {
	err := e.Move(e.Path, "")
	if err != nil {
		log.Printf("Unable to move file from { %s } to { %s }: %v\n", e.Path, e.Destination, err)
	} else {
		log.Println("Event has been processed.")
	}
}

// ProcessDirEvent takes an event and scans files ext
func (pc *PathConsumer) ProcessDirEvent(e *event.Event) {
	files, err := utils.ScanFilesInDir(e.Path)

	if err != nil {
		log.Fatalf("Unable to scan files in dir event: error: %v, path: %s", err, e.Path)
	}

	for file := range files {
		if utils.ExtractFileExt(file) == pc.Ext {
			ev := event.New(file, e.Path, e.Destination, pc.Ext)
			err = ev.Move(ev.Path, "/"+file)

			if err != nil {
				log.Printf("Unable to move file: %s from path: %s to dest: %s: %v", file, ev.Path, ev.Destination, err)
			}
		}
	}
}

// AddPath adds a path to the watcher
func (pw *PathWatcher) AddPath(path string) {
	pw.Watcher.Add(path)
}

// Register a consumer to the producer
func (pw *PathWatcher) Register(consumer *Consumer) {
	pw.Consumers = append(pw.Consumers, consumer)
}

// Unregister a consumer from the producer
func (pw *PathWatcher) Unregister(consumer *Consumer) {
	for i, cons := range pw.Consumers {
		if cons == consumer {
			pw.Consumers[i] = pw.Consumers[len(pw.Consumers)-1]
			pw.Consumers = pw.Consumers[:len(pw.Consumers)-1]
		}
	}
}

// Observe the producer
func (pw *PathWatcher) Observe(pollInterval int) {
	eventQueue := NewQueue()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Could not create new watcher: %v", err)
	}

	defer watcher.Close()

	// fsnotify doesnt support recursive folders, so we can here
	if err := filepath.Walk(pw.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Error walking path structure. Please ensure to use absolute path: %v", err)
		}

		if info.Mode().IsDir() {
			watcher.Add(path)
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not parse recursive path: %v", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op.String() == "CREATE" && utils.IsDir(event.Name) {
					watcher.Add(event.Name)
				} else if event.Op.String() == "CLOSEWRITE" {
					// Most files being written to should have a IN_CLOSE_WRITE event for the file name itself.
					// One exception is .part files, wherein the close_write event occurs on the .part file, and then merge the .part files into the expected file name.
					// The problem with this is that there is no close_write event for this, only a new CREATE event.
					// Need to come up with a solution to account for this.

					// Check if the event is in the specialExt list
					// If it is, then throw it on a queue, as we'd expect to see a create event.

					ev := newEvent(event.Name, event.Op.String())
					log.Printf("IN_CLOSE_WRITE EVENT: %v\n", ev)

					if specialWatchedFilExts[ev.Ext] == 1 {
						log.Printf("Ext added to queue: %s\n\n", ev.Ext)
						eventQueue.Add(*ev)
					} else {
						pw.Notify(ev.Path, ev.Operation)
					}
				} else if event.Op.String() == "CREATE" {
					// Check against the specialExt queue for files matching
					// this new CREATE event.
					// Notify subscribers of the event if valid.
					createEvent := newEvent(event.Name, event.Op.String())
					log.Printf("CREATE EVENT: %v\n", createEvent)
					log.Printf("QUEU_SIZE: %d\n", eventQueue.Size())

					for hsh, ev := range eventQueue.Queue {
						log.Printf("EVENT: %v, HSH: %s - CREATE_EVENT: %v\n\n", createEvent, hsh, ev)
						log.Printf("Paths: %s = %s", utils.ExtractPathWithoutExt(ev.Path), utils.ExtractPathWithoutExt(createEvent.Path))
						if utils.CompareFilePaths(ev.Path, createEvent.Path) {
							pw.Notify(createEvent.Path, createEvent.Operation)
							eventQueue.Remove(hsh)
						}
					}
				}
			case err := <-watcher.Errors:
				log.Printf("Watcher encountered an error when observing %s: %v", pw.Path, err)
			}
		}
	}()

	<-done
}

// Notify consumers of an event
func (pw *PathWatcher) Notify(path, event string) {
	for _, cons := range pw.Consumers {
		(*cons).Receive(path, event)
	}
}

func newEvent(path, ev string) *event.Event {
	return &event.Event{
		File:      filepath.Base(path),
		Path:      path,
		Ext:       utils.ExtractFileExt(path),
		Timestamp: time.Now(),
		Operation: ev,
	}
}
