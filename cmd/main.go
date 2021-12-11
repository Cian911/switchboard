package main

import (
	"github.com/cian911/switchboard/watcher"
)

func main() {
	path := "/Users/cian.gallagher/input"
	destination := "/Users/cian.gallagher/output"

	var pw watcher.Producer = &watcher.PathWatcher{
		Path: path,
	}

	var pc watcher.Consumer = &watcher.PathConsumer{
		Path:        path,
		Destination: destination,
		Ext:         ".txt",
	}

	pw.Register(&pc)
	pw.Observe()
}
