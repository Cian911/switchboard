package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ExtractFileExt returns the file extension of a file
func ExtractFileExt(path string) string {
	// If the path is a directory, returns empty string
	if ValidatePath(path) && IsDir(path) {
		return ""
	}

	return filepath.Ext(strings.Trim(path, "'"))
}

// Extract file path without the extension
func ExtractPathWithoutExt(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

// Compare two filepaths and return a bool
func CompareFilePaths(p1, p2 string) bool {
	if ExtractPathWithoutExt(p1) == p2 {
		return true
	}

	return false
}

// ValidatePath checks if a path exists
func ValidatePath(path string) bool {
	if path == "" {
		return false
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("EVENT: Path does not exist: %s", path)
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

// IsDir returns a boolean if the given path is a directory
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Printf("Could not find path: %v", err)
		return false
	}

	if fileInfo.IsDir() {
		return true
	}

	return false
}
