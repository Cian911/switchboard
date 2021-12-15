package event

import "testing"

var (
	event = &Event{
		File:        "readme.txt",
		Path:        "/input",
		Destination: "/output",
		Ext:         ".txt",
		Operation:   "CREATE",
	}

	ext = ".txt"
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
}
