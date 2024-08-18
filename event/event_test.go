package event

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
)

var (
	event = &Event{
		File:        "readme.txt",
		Path:        "/input",
		Destination: "/output",
		Ext:         ".txt",
		Operation:   "CREATE",
	}

	ext  = ".txt"
	file = "sample.txt"
)

func TestEvent(t *testing.T) {
	t.Run("It returns true when event is valid", func(t *testing.T) {
		want := true
		got := event.IsValidEvent(ext)

		if want != got {
			t.Errorf("event extension is not valid, when it should have been: want=%t, got=%t", want, got)
		}
	})

	t.Run("It returns false when event is valid", func(t *testing.T) {
		want := false
		got := event.IsValidEvent(".mp4")

		if want != got {
			t.Errorf("event extension is not valid, when it should have been: want=%t, got=%t", want, got)
		}
	})

	t.Run("It moves file from one dir to another dir", func(t *testing.T) {
		event := eventSetup(t)
		// event := helpers.TestEventSetup(t)
		event.Move(event.Path, "")

		// If the file does not exist, log an error
		if _, err := os.Stat(fmt.Sprintf("%s/%s", event.Destination, event.File)); errors.Is(err, os.ErrNotExist) {
			t.Fatalf("Failed to move from %s/%s to %s/%s: %v : %v", event.Path, event.File, event.Destination, event.File, err, *event)
		}

		// If the file still exists in the source directory, log an error
		if _, err := os.Stat(fmt.Sprintf("%s/%s", event.Path, event.File)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("Failed to delete file from %s/%s to %s/%s after Move: %v : %v", event.Path, event.File, event.Destination, event.File, err, *event)
		}
	})

	t.Run("It moves file from one dir to another dir with valid destPath", func(t *testing.T) {
		event := eventSetup(t)
		event.Move(fmt.Sprintf("%s/", event.Path), "")

		// If the file does not exist, log an error
		if _, err := os.Stat(fmt.Sprintf("%s/%s", event.Destination, event.File)); errors.Is(err, os.ErrNotExist) {
			t.Fatalf("Failed to move from %s/%s to %s/%s: %v : %v", event.Path, event.File, event.Destination, event.File, err, *event)
		}

		// If the file still exists in the source directory, log an error
		if _, err := os.Stat(fmt.Sprintf("%s/%s", event.Path, event.File)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("Failed to delete file from %s/%s to %s/%s after Move: %v : %v", event.Path, event.File, event.Destination, event.File, err, *event)
		}
	})

	t.Run("It does not move file from one dir to another dir", func(t *testing.T) {
		event := eventSetup(t)
		event.Destination = "/abcdefg"
		err := event.Move(event.Path, "")

		if err == nil {
			log.Fatal("event.Move() should have thrown error but didn't.")
		}
	})

	t.Run("It determines if the event is a new dir", func(t *testing.T) {
		event := eventSetup(t)
		event.File = "input"
		event.Ext = ""

		want := true
		got := event.IsNewDirEvent()

		if want != got {
			t.Errorf("Wanted new dir event but didn't get it: want=%t, got=%t", want, got)
		}
	})
}

func eventSetup(t *testing.T) *Event {
	path := t.TempDir()
	_, err := os.CreateTemp(path, file)
	if err != nil {
		t.Fatalf("Unable to create temp file: %v", err)
	}

	return &Event{
		File:        file,
		Path:        path,
		Destination: t.TempDir(),
		Ext:         ext,
		Operation:   "CREATE",
	}
}
