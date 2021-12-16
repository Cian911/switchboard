package event

import (
	"errors"
	"fmt"
	"io/ioutil"
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
		event.Move()

		// If the file does not exist, log an error
		if _, err := os.Stat(fmt.Sprintf("%s/%s", event.Destination, event.File)); errors.Is(err, os.ErrNotExist) {
			t.Fatalf("Failed to move from %s/%s to %s/%s: %v : %v", event.Path, event.File, event.Destination, event.File, err, *event)
		}
	})

	t.Run("It does not move file from one dir to another dir", func(t *testing.T) {
		event := eventSetup(t)
		event.Destination = "/abcdefg"
		err := event.Move()

		if err == nil {
			log.Fatal("event.Move() should have thrown error but didn't.")
		}
	})
}

func eventSetup(t *testing.T) *Event {
	path := t.TempDir()
	_, err := ioutil.TempFile(path, file)

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
