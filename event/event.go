package event

import (
	"fmt"
	"log"
	"os"

	"github.com/cian911/switchboard/utils"
)

// Event is a struct that holds the information for a file event
type Event struct {
	// File is the name of the file
	File string
	// Path is the path to the file
	Path string
	// Destination is the path to the destination
	Destination string
	// Ext is the file extension
	Ext string
	// Operation is the operation that was performed
	Operation string
	// IsDir is the new create vent a directory
	IsDir bool
}

func New(file, path, dest, ext string) *Event {
	return &Event{
		File:        file,
		Path:        path,
		Destination: dest,
		Ext:         ext,
	}
}

// Move moves the file to the destination
func (e *Event) Move(path string) error {
	log.Printf("Moving e.Path: %s to %s/%s\n", path, e.Destination, e.File)

	err := os.Rename(path, fmt.Sprintf("%s/%s", e.Destination, e.File))
	return err
}

// IsValidEvent checks if the event operation and file extension is valid
func (e *Event) IsValidEvent(ext string) bool {
	if ext == e.Ext && e.Operation == "CREATE" {
		return true
	}

	return false
}

func (e *Event) IsNewDirEvent() bool {
	if e.Ext == "" && utils.ValidatePath(e.Path) {
		return true
	}

	return false
}
