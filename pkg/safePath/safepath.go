package safepath

import (
	"fmt"
	"path/filepath"
	"strings"
)

func GetSafeFilePath(baseDir, uuid string) (string, error) {
	if !filepath.IsAbs(baseDir) {
		return "", fmt.Errorf("base directory is not absolute: %s", baseDir)
	}

	fullPath := filepath.Join(baseDir, uuid)
	cleanedPath := filepath.Clean(fullPath)
	if !strings.HasPrefix(cleanedPath, baseDir) {
		return "", fmt.Errorf("path traversal detected for uuid: %s", uuid)
	}

	return cleanedPath, nil
}
