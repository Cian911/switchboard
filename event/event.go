package event

import (
	"fmt"
	"log"
	"os"
)

type Event struct {
	File        string
	Path        string
	Destination string
	Ext         string
	Operation   string
}

func (e *Event) Move() {
	log.Printf("Moving e.Path: %s to %s/%s\n", e.Path, e.Destination, e.File)
	err := os.Rename(e.Path, fmt.Sprintf("%s/%s", e.Destination, e.File))
	if err != nil {
		log.Fatalf("Unable to move file from { %s } to { %s }: %v", e.Path, e.Destination, err)
	}
}

func (e *Event) IsValidEvent(ext string) bool {
	if ext == e.Ext && e.Operation == "CREATE" {
		return true
	}

	return false
}
