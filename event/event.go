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

func (e *Event) Move() error {
	log.Printf("Moving e.Path: %s to %s/%s\n", e.Path, e.Destination, e.File)

	err := os.Rename(e.Path, fmt.Sprintf("%s/%s", e.Destination, e.File))
	return err
}

func (e *Event) IsValidEvent(ext string) bool {
	if ext == e.Ext && e.Operation == "CREATE" {
		return true
	}

	return false
}
