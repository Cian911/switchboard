package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Producer interface {
	Register(consumer *Consumer)
	Unregister(consumer *Consumer)
	notify(path, event string)
	Observe()
}

type Consumer interface {
	Receive(path, event string)
	Process(path, destination string)
}

type PathWatcher struct {
	Consumers []*Consumer
	Watcher   fsnotify.Watcher
	Path      string
}

type PathConsumer struct {
	Path        string
	Destination string
}

func (pc *PathConsumer) Receive(path, event string) {
	log.Printf("CONSUMER EVENT: path: %s, event: %s", path, event)
}

func (pc *PathConsumer) Process(path, destination string) {
	log.Printf("CONSUMER PROCESSING EVENT: path: %s, event: %s", path, destination)
}

func (pw *PathWatcher) Register(consumer *Consumer) {
	pw.Consumers = append(pw.Consumers, consumer)
}

func (pw *PathWatcher) Unregister(consumer *Consumer) {
	for i, cons := range pw.Consumers {
		if cons == consumer {
			pw.Consumers[i] = pw.Consumers[len(pw.Consumers)-1]
			pw.Consumers = pw.Consumers[:len(pw.Consumers)-1]
		}
	}
}

func (pw *PathWatcher) notify(path, event string) {
	for _, cons := range pw.Consumers {
		(*cons).Receive(path, event)
	}
}

func (pw *PathWatcher) Observe() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Could not create new watcher: %v", err)
	}

	defer watcher.Close()

	// fsnotify doesnt support recursive folders, so we can here
	if err := filepath.Walk(pw.Path, func(path string, info os.FileInfo, err error) error {
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
				pw.notify(event.Name, event.Op.String())
			case err := <-watcher.Errors:
				log.Printf("Watcher encountered an error when observing %s: %v", pw.Path, err)
			}
		}
	}()

	<-done
}
