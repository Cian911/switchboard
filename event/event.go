package event

import (
	"fmt"
	"log"
	"os"
)

// Event is a struct that holds the information for a file event
type Event struct {
	File        string
	Path        string
	Destination string
	Ext         string
	Operation   string
}

// Move moves the file to the destination
func (e *Event) Move() error {
	log.Printf("Moving e.Path: %s to %s/%s\n", e.Path, e.Destination, e.File)

	err := os.Rename(e.Path, fmt.Sprintf("%s/%s", e.Destination, e.File))
	return err
}

// IsValidEvent checks if the event operation and file extension is valid
func (e *Event) IsValidEvent(ext string) bool {
	if ext == e.Ext && e.Operation == "CREATE" {
		return true
	}

	return false
}
