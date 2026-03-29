package external

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/easy-cloud-Knet/KWS_Core/internal/config"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func resolveSnapshotRoot(opts *ExternalSnapshotOptions) (string, error) {
	if opts == nil || opts.BaseDir == "" {
		return config.StorageBase, nil
	}

	clean := filepath.Clean(opts.BaseDir)
	if !filepath.IsAbs(clean) {
		return "", virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("base dir must be absolute"))
	}
	if strings.Contains(clean, "..") {
		return "", virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("invalid base dir"))
	}

	return clean, nil
}

func resolveDomainUUID(domain externalSnapshotDomain) (string, error) {
	if domain == nil {
		return "", virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	uuid, err := domain.UUIDString()
	if err != nil {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to get domain uuid: %w", err))
	}

	return uuid, nil
}
