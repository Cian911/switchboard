package main

import (
	"github.com/cian911/switchboard/watcher"
)

func main() {
	path := "/Users/cian.gallagher/test_events"
	destination := "/tmp"

	var pw watcher.Producer = &watcher.PathWatcher{
		Path: path,
	}

	var pc watcher.Consumer = &watcher.PathConsumer{
		Path:        path,
		Destination: destination,
	}

	pw.Register(&pc)
	pw.Observe()
}
