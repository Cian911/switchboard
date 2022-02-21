package utils

import (
	"io/ioutil"
	"os"
	"testing"
)

const (
	filePath    = "/home/test/file.mp4"
	dirPath     = "/home/test"
	diffDirPath = "/home/test.movies/"
)

func TestUtils(t *testing.T) {
	t.Run("It tests ExtractFileExt()", func(t *testing.T) {
		tempDir := setupTempDir("'Author- [0.5] Version - Title (azw3 epub mobi)'", t)
		defer os.RemoveAll(tempDir)

		tests := []struct {
			input          string
			expectedOutput string
		}{
			{"/home/test/movie.mp4", ".mp4"},
			{"/home/test/movie.maaaap4.aaa.mp4", ".mp4"},
			{"/home/test/", ""},
			{"/home/test", ""},
			{"/home/test/movie.mp4/", ""},
			{"/home/test/'movie.mp4'", ".mp4"},
			{"/home/test/movie.mp4.part", ".part"},
			{"/home/test/Some weird folder name ([0.5] epub) sample.mp4", ".mp4"},
			{tempDir, ""},
		}

		for _, tt := range tests {
			got := ExtractFileExt(tt.input)

			if got != tt.expectedOutput {
				t.Errorf("Failed extracting file extention: got=%s, want=%s", got, tt.expectedOutput)
			}
		}
	})

	t.Run("It tests ValidatePath()", func(t *testing.T) {
		tempDir := setupTempDir("'Author - [0.5] - Title (azw3 epub mobi)'", t)
		defer os.RemoveAll(tempDir)

		tests := []struct {
			input          string
			expectedOutput bool
		}{
			{"", false},
			{"/abba/asdas/asda", false},
			{t.TempDir(), true},
			{tempDir, true},
		}

		for _, tt := range tests {
			got := ValidatePath(tt.input)

			if got != tt.expectedOutput {
				t.Errorf("Failed extracting file extention: got=%t, want=%t: input: %s", got, tt.expectedOutput, tt.input)
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

	t.Run("It tests ScanFilesInDir()", func(t *testing.T) {
		tempDir := setupTempDir("'Author - [0.5] - Title (azw3 epub mobi)'", t)
		// TODO: Make this more testable..
		_ = setupTempFile("'Author - [0.5] - Title.azw3'", tempDir, t)
		_ = setupTempFile("'Author - [0.5] - Title.epub'", tempDir, t)

		tests := []struct {
			expectedOutput int
		}{
			{2},
		}

		for _, tt := range tests {
			got, err := ScanFilesInDir(tempDir)

			if err != nil {
				t.Fatalf("Could not scan files in dir: %v", err)
			}

			if len(got) != tt.expectedOutput {
				t.Errorf("Failed scanning files in dir: got=%v, want=%d", got, tt.expectedOutput)
			}
		}
	})

	t.Run("It tests IsDir()", func(t *testing.T) {
		tempDir := setupTempDir("'Author- [0.5] Version - Title (azw3 epub mobi)'", t)
		defer os.RemoveAll(tempDir)

		tests := []struct {
			input          string
			expectedOutput bool
		}{
			{"/home/test/movie.mp4", false},
			{"/home/test/movie.maaaap4.aaa.mp4", false},
			{"/home/test/", false},
			{"/home/test", false},
			{"/home/test/Some weird folder name ([0.5] epub) sample.mp4", false},
			{tempDir, true},
			{"/tmp", true},
		}

		for _, tt := range tests {
			got := IsDir(tt.input)

			if got != tt.expectedOutput {
				t.Errorf("%s should be a dir but returned failed. want=%t, got=%t", tt.input, tt.expectedOutput, got)
			}
		}
	})
}

func setupTempDir(name string, t *testing.T) string {
	tempDir, err := ioutil.TempDir("", name)
	if err != nil {
		t.Errorf("Could not create temp dir: %v", err)
	}

	return tempDir
}

func setupTempFile(name, dir string, t *testing.T) *os.File {
	file, err := ioutil.TempFile(dir, name)

	if err != nil {
		t.Fatalf("Unable to create temp file: %v", err)
	}

	return file
}
