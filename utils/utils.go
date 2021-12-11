package utils

import "path/filepath"

func ExtractFileExt(path string) string {
	return filepath.Ext(path)
}
