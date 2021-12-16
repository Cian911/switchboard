package utils

import "testing"

const (
	filePath    = "/home/test/file.mp4"
	dirPath     = "/home/test"
	diffDirPath = "/home/test.movies/"
)

func TestUtils(t *testing.T) {
	t.Run("It tests ExtractFileExt()", func(t *testing.T) {
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
	})

	t.Run("It tests ValidatePath()", func(t *testing.T) {
		tests := []struct {
			input          string
			expectedOutput bool
		}{
			{"", false},
			{"/abba/asdas/asda", false},
			{t.TempDir(), true},
		}

		for _, tt := range tests {
			got := ValidatePath(tt.input)

			if got != tt.expectedOutput {
				t.Errorf("Failed extracting file extention: got=%t, want=%t", got, tt.expectedOutput)
			}
		}
	})

	t.Run("It tests ValidateFileExt()", func(t *testing.T) {
		tests := []struct {
			input          string
			expectedOutput bool
		}{
			{"", false},
			{"file", false},
			{"file...txt", false},
			{".txt", true},
		}

		for _, tt := range tests {
			got := ValidateFileExt(tt.input)

			if got != tt.expectedOutput {
				t.Errorf("Failed extracting file extention: got=%t, want=%t", got, tt.expectedOutput)
			}
		}
	})
}
