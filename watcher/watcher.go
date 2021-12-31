package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cian911/switchboard/event"
	"github.com/cian911/switchboard/utils"
	"github.com/fsnotify/fsnotify"
)

// Producer interface for the watcher
// Must implement Register(), Unregister(), and Observe(), and notify()
type Producer interface {
	// Register a consumer to the producer
	Register(consumer *Consumer)
	// Unregister a consumer from the producer
	Unregister(consumer *Consumer)
	// Notify consumers of an event
	notify(path, event string)
	// Observe the producer
	Observe()
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
	log.Printf("Event Received: %s, Path: %s\n", ev, path)

	// TODO: Move IsNewDirEvent to utils and call func on event struct
	// TODO: If is a dir event, there should not be a file ext
	e := &event.Event{
		File:        filepath.Base(path),
		Path:        path,
		Destination: pc.Destination,
		Ext:         utils.ExtractFileExt(path),
		Operation:   ev,
	}

	if e.IsNewDirEvent() {
		log.Println("Event is a new dir")

		// Recursively scan dir for items with our ext
		// Then add all recursive dirs as paths
		pc.ProcessDirEvent(e)
	} else if e.IsValidEvent(pc.Ext) {
		log.Println("Event is valid")
		pc.Process(e)
	}
}

// Process takes an event and moves it to the destination
func (pc *PathConsumer) Process(e *event.Event) {
	err := e.Move(e.Path, "")
	if err != nil {
		log.Fatalf("Unable to move file from { %s } to { %s }: %v", e.Path, e.Destination, err)
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
func (pw *PathWatcher) Observe() {
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
				}

				pw.notify(event.Name, event.Op.String())
			case err := <-watcher.Errors:
				log.Printf("Watcher encountered an error when observing %s: %v", pw.Path, err)
			}
		}
	}()

	<-done
}

// Notify consumers of an event
func (pw *PathWatcher) notify(path, event string) {
	for _, cons := range pw.Consumers {
		(*cons).Receive(path, event)
	}
}
