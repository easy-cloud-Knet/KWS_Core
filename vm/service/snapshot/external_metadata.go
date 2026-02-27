package snapshot

import (
	"fmt"
	"path/filepath"
	"strings"
)

// TODO: 실제 사용 시 아래 두 nolint 태그 제거할 것
const defaultSnapshotRoot = "/var/lib/kws" //nolint:unused

type ExternalSnapshotOptions struct {
	BaseDir     string
	Description string
	Quiesce     bool
	Live        bool
}

func resolveSnapshotRoot(opts *ExternalSnapshotOptions) (string, error) { //nolint:unused
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
