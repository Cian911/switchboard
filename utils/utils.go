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

func ScanFilesInDir(path string) (map[string]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fileList := map[string]string{}
	/* for _, file := range files { */
	/*      */
	/* } */

	return fileList, nil
}
