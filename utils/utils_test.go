package utils

import "testing"

const (
	filePath    = "/home/test/file.mp4"
	dirPath     = "/home/test"
	diffDirPath = "/home/test.movies/"
)

func TestUtils(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"/home/test/movie.mp4", ".mp4"},
		{"/home/test/movie.maaaap4.aaa.mp4", ".mp4"},
		{"/home/test/", ""},
		{"/home/test", ""},
		{"/home/test/movie.mp4/", ""},
	}

	for _, tt := range tests {
		got := ExtractFileExt(tt.input)

		if got != tt.expectedOutput {
			t.Errorf("Failed extracting file extention: got=%s, want=%s", got, tt.expectedOutput)
		}
	}
}
