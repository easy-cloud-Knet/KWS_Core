package safepath

import (
	"path/filepath"
	"strings"
)

func GetSafeFilePath(baseDir, uuid string) (string, bool) {
	if !filepath.IsAbs(baseDir) {
		return "", false
	}

	fullPath := filepath.Join(baseDir, uuid)
	cleanedPath := filepath.Clean(fullPath)
	if !strings.HasPrefix(cleanedPath, baseDir) {
		return "", false
	}

	return cleanedPath, true
}
