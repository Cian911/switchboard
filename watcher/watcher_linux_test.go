//go:build linux
// +build linux

package watcher

import (
	"os"
	"testing"
	"time"

	"github.com/cian911/switchboard/utils"
)

func TestObserve(t *testing.T) {
	t.Run("CLOSE_WRITE", func(t *testing.T) {
		HelperPath = t.TempDir()
		HelperDestination = t.TempDir()
		HelperExt = ".txt"

		pw, pc := TestProducerConsumer()
		pw.Register(&pc)

		go pw.Observe(1)
		<-time.After(1 * time.Second)
		// Fire event
		h, err := os.Create(HelperPath + "/sample2.txt")
		if err != nil {
			t.Fatalf("Failed to create file in testdir: %v", err)
		}
		h.Close()
		os.OpenFile(HelperPath+"/sample.txt", 0, os.FileMode(int(0777)))
		<-time.After(3 * time.Second)

		files, _ := utils.ScanFilesInDir(HelperDestination)

		if len(files) != 1 {
			t.Errorf("CLOSE_WRITE event was not processed - want: %d, got: %d", 1, len(files))
		}
	})

	t.Run("CREATE", func(t *testing.T) {
		HelperPath = t.TempDir()
		HelperDestination = t.TempDir()
		HelperExt = ".txt"

		pw, pc := TestProducerConsumer()
		pw.Register(&pc)

		go pw.Observe(1)
		<-time.After(3 * time.Second)
		// Fire event
		os.Create(HelperPath + "/sample2.txt")
		<-time.After(3 * time.Second)

		files, _ := utils.ScanFilesInDir(HelperDestination)

		if len(files) != 1 {
			t.Errorf("CREATE event was not processed - want: %d, got: %d", 1, len(files))
		}
	})

	t.Run("WRITE", func(t *testing.T) {
	})
}
