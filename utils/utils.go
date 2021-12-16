package utils

import (
	"os"
	"path/filepath"
)

func ExtractFileExt(path string) string {
	return filepath.Ext(path)
}

func ValidatePath(path string) bool {
	if path == "" {
		return false
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func ValidateFileExt(ext string) bool {
	if ext == "" || ext[0:1] != "." {
		return false
	}

	return true
}
