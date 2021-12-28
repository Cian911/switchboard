package watcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/cian911/switchboard/event"
	"github.com/cian911/switchboard/utils"
)

const (
	path        = "/tmp"
	destination = "/test"
)

var (
	e = &event.Event{
		File:        "readme.txt",
		Path:        "/input",
		Destination: "/output",
		Ext:         ".txt",
		Operation:   "CREATE",
	}

	ext  = ".txt"
	file = "sample.txt"
)

func TestWatcher(t *testing.T) {
	t.Run("It registers a consumer", func(t *testing.T) {
		pw, pc := setup()

		pw.Register(&pc)

		if len(pw.(*PathWatcher).Consumers) != 1 {
			t.Fatalf("Consumer was not registered when it should have been. want=%d, got=%d", 1, len(pw.(*PathWatcher).Consumers))
		}
	})

	t.Run("It unregisters a consumer", func(t *testing.T) {
		pw, pc := setup()

		pw.Register(&pc)
		pw.Unregister(&pc)

		if len(pw.(*PathWatcher).Consumers) != 0 {
			t.Fatalf("Consumer was not unregistered when it should have been. want=%d, got=%d", 0, len(pw.(*PathWatcher).Consumers))
		}
	})

	t.Run("It processes a new dir event", func(t *testing.T) {
		pw, pc := setup()

		pw.Register(&pc)
		pw.Unregister(&pc)

		ev := eventSetup(t)
		ev.Path = t.TempDir()
		ev.File = utils.ExtractFileExt(ev.Path)

		for i := 1; i <= 3; i++ {
			createTempFile(ev.Path, ".txt", t)
		}

		pc.Receive(ev.Path, "CREATE")

		// Scan dest dir for how many files it contains
		// if want == got, all files have been moved successfully
		filesInDir, err := utils.ScanFilesInDir(ev.Destination)
		if err != nil {
			t.Fatalf("Could not scan all files in destination dir: %v", err)
		}

		want := 3
		got := len(filesInDir)

		if want != got {
			t.Fatalf("want: %d != got: %d", want, got)
		}
	})
}

func setup() (Producer, Consumer) {
	var pw Producer = &PathWatcher{
		Path: path,
	}

	var pc Consumer = &PathConsumer{
		Path:        path,
		Destination: destination,
	}

	return pw, pc
}

func eventSetup(t *testing.T) *event.Event {
	path := t.TempDir()
	_, err := ioutil.TempFile(path, file)

	if err != nil {
		t.Fatalf("Unable to create temp file: %v", err)
	}

	return &event.Event{
		File:        file,
		Path:        path,
		Destination: t.TempDir(),
		Ext:         ext,
		Operation:   "CREATE",
	}
}

func createTempFile(path, ext string, t *testing.T) {
	data := []byte("hello\nworld\n")
	fileName := fmt.Sprintf("%d.%s", time.Now().Unix(), ext)
	err := os.WriteFile(fmt.Sprintf("%s/%s", path, fileName), data, 0644)

	if err != nil {
		t.Fatalf("Could not create test file: %v", err)
	}
}
