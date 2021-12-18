package utils

import (
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
