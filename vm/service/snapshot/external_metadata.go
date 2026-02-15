package snapshot

import (
	"fmt"
	"path/filepath"
	"strings"
)

const defaultSnapshotRoot = "/var/lib/kws"

type ExternalSnapshotOptions struct {
	BaseDir     string
	Description string
	Quiesce     bool
	Live        bool
}

func resolveSnapshotRoot(opts *ExternalSnapshotOptions) (string, error) {
	if opts == nil || opts.BaseDir == "" {
		return defaultSnapshotRoot, nil
	}

	clean := filepath.Clean(opts.BaseDir)
	if !filepath.IsAbs(clean) {
		return "", fmt.Errorf("base dir must be absolute")
	}
	if strings.Contains(clean, "..") {
		return "", fmt.Errorf("invalid base dir")
	}

	return clean, nil
}
