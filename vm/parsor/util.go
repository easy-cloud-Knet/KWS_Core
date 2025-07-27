package parsor

import (
	"path/filepath"
	"strings"
)


func UUIDValidator(uuid string) bool {
    if len(uuid) != 36 {
        return false
    }

    // A quick check for path traversal characters.
    // This is a good first line of defense.
    if strings.ContainsAny(uuid, "/\\.") {
        return false
    }

    // Your original UUID format validation loop
    for i, c := range uuid {
        if (i == 8 || i == 13 || i == 18 || i == 23) && c != '-' {
            return false
        }
        if (i != 8 && i != 13 && i != 18 && i != 23) && (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
            return false
        }
    }
    return true
}

func GetSafeFilePath(baseDir, uuid string) (string, bool) {
    if !UUIDValidator(uuid) {
        return "", false
    }

    fullPath := filepath.Join(baseDir, uuid)
    cleanedPath := filepath.Clean(fullPath)
    if !strings.HasPrefix(cleanedPath, baseDir) {
        return "", false
    }

    return cleanedPath, true
}
