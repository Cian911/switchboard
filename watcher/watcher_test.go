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
		pw, pc := setup(path, destination, ext, "")

		pw.Register(&pc)

		if len(pw.(*PathWatcher).Consumers) != 1 {
			t.Fatalf("Consumer was not registered when it should have been. want=%d, got=%d", 1, len(pw.(*PathWatcher).Consumers))
		}
	})

	t.Run("It unregisters a consumer", func(t *testing.T) {
		pw, pc := setup(path, destination, ext, "")

		pw.Register(&pc)
		pw.Unregister(&pc)

		if len(pw.(*PathWatcher).Consumers) != 0 {
			t.Fatalf("Consumer was not unregistered when it should have been. want=%d, got=%d", 0, len(pw.(*PathWatcher).Consumers))
		}
	})

	t.Run("It processes a new dir event", func(t *testing.T) {

		ev := eventSetup(t)
		ev.Path = t.TempDir()
		ev.File = utils.ExtractFileExt(ev.Path)

		pw, pc := setup(ev.Path, ev.Destination, ev.Ext, "")
		pw.Register(&pc)
		pw.Unregister(&pc)

		for i := 1; i <= 3; i++ {
			file := createTempFile(ev.Path, ".txt", t)
			defer os.Remove(file.Name())
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
			t.Fatalf("want: %d != got: %d - debug - files: %v, event: %v", want, got, filesInDir, ev)
		}
	})

	t.Run("It pocesses a pattern matched event", func(t *testing.T) {
		ev := eventSetup(t)
		ev.File = utils.ExtractFileExt(ev.Path)

		pw, pc := setup(ev.Path, ev.Destination, ev.Ext, "[0-9]+.txt")
		pw.Register(&pc)
		pw.Unregister(&pc)

		for i := 1; i <= 3; i++ {
			// Create 3 temp files
			file := createTempFile(ev.Path, ".txt", t)
			defer os.Remove(file.Name())
		}

		pc.Receive(ev.Path, "CREATE")

		// Scan dest dir for how many files it contains
		// if want == got, all files have been moved successfully
		filesInDir, err := utils.ScanFilesInDir(ev.Destination)
		t.Logf("Files in Dir: %v", filesInDir)
		if err != nil {
			t.Fatalf("Could not scan all files in destination dir: %v", err)
		}

		want := 3
		got := len(filesInDir)

		if want != got {
			t.Fatalf("want: %d != got: %d - debug - files: %v, event: %v", want, got, filesInDir, ev)
		}
	})
}

func setup(p, d, e, rp string) (Producer, Consumer) {
	var pw Producer = &PathWatcher{
		Path: p,
	}

	regexPattern, _ := utils.ValidateRegexPattern(rp)

	var pc Consumer = &PathConsumer{
		Path:        p,
		Destination: d,
		Ext:         e,
		Pattern:     *regexPattern,
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
		Timestamp:   time.Now(),
	}
}

func createTempFile(path, ext string, t *testing.T) *os.File {
	file, err := ioutil.TempFile(path, fmt.Sprintf("*%s", ext))
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}

	return file
}
