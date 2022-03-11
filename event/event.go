package event

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cian911/switchboard/utils"
)

var validOperations = map[string]bool{
	"CREATE": true,
	"WRITE":  true,
}

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
	// Timestamp in unix time epoch
	Timestamp time.Time
}

// New creates and returns a new event struct
func New(file, path, dest, ext string) *Event {
	return &Event{
		File:        file,
		Path:        path,
		Destination: dest,
		Ext:         ext,
		Timestamp:   time.Now(),
	}
}

// Move moves the file to the destination
func (e *Event) Move(path, file string) error {
	log.Printf("Moving e.Path: %s to %s/%s\n", path, e.Destination, e.File)

	sourcePath := filepath.Join(path, file)
	destPath := filepath.Join(e.Destination, e.File)

	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

// IsValidEvent checks if the event operation and file extension is valid
func (e *Event) IsValidEvent(ext string) bool {
	if ext == e.Ext && validOperations[e.Operation] {
		return true
	}

	return false
}

// IsNewDirEvent returns a bool if the given path is a directory or not
func (e *Event) IsNewDirEvent() bool {
	if e.Ext == "" && utils.ValidatePath(e.Path) && utils.IsDir(e.Path) {
		return true
	}

	return false
}
