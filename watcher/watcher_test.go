package watcher

import "testing"

const (
	path        = "/tmp"
	destination = "/test"
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
