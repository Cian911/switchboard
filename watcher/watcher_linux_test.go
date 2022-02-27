//go:build linux
// +build linux

package watcher

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/cian911/switchboard/event"
)

func TestObserve(t *testing.T) {
	return true
}

func testEventSetup(op string, t *testing.T) *event.Event {
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
		Operation:   op,
		Timestamp:   time.Now(),
	}
}
