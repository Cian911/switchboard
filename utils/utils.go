package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// ExtractFileExt returns the file extension of a file
func ExtractFileExt(path string) string {
	return filepath.Ext(path)
}

// ValidatePath checks if a path exists
func ValidatePath(path string) bool {
	if path == "" {
		return false
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// ValidateFileExt checks if a file extension is valid
func ValidateFileExt(ext string) bool {
	if ext == "" || ext[0:1] != "." {
		return false
	}

	return true
}

// ScanFilesInDir scans and returns a map of all files in a given directory.
// The returned map is a key value of the filename, and if a directory.
// im_a_dir -> true
// sample.txt -> false
func ScanFilesInDir(path string) (map[string]bool, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// "sample.txt" -> false
	// "folder1" -> true
	fileList := make(map[string]bool)
	for _, file := range files {
		fileList[file.Name()] = file.IsDir()
	}

	return fileList, nil
}
