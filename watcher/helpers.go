package watcher

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/cian911/switchboard/event"
	"github.com/cian911/switchboard/utils"
)

var (
	HelperPath        string
	HelperDestination string
	HelperFile        string
	HelperExt         string
	HelperOperation   string
	HelperPattern     string
)

func TestEventSetup(t *testing.T) *event.Event {
	path := t.TempDir()
	_, err := ioutil.TempFile(path, HelperFile)

	if err != nil {
		t.Fatalf("Unable to create temp file: %v", err)
	}

	return &event.Event{
		File:        HelperFile,
		Path:        path,
		Destination: t.TempDir(),
		Ext:         HelperExt,
		Operation:   HelperOperation,
		Timestamp:   time.Now(),
	}
}

// SimulateMultipleEvents takes a list of operations as args
// ["CREATE", "WRITE", "IN_CLOSE_WRITE"]
func TestSimulateMultipleEvents(operationList []string, t *testing.T) []event.Event {
	eventList := []event.Event{}

	for _, op := range operationList {
		HelperOperation = op
		eventList = append(eventList, *TestEventSetup(t))
	}

	return eventList
}

func TestProducerConsumer() (Producer, Consumer) {
	var pw Producer = &PathWatcher{
		Path: HelperPath,
	}

	pattern, _ := utils.ValidateRegexPattern(HelperPattern)

	var pc Consumer = &PathConsumer{
		Path:        HelperPath,
		Destination: HelperDestination,
		Ext:         HelperExt,
		Pattern:     *pattern,
	}

	return pw, pc
}

func TestPathWatcher() *PathWatcher {
	return &PathWatcher{
		Queue: NewQueue(),
	}
}
