package external

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// SnapshotFilePath returns the path of the primary disk snapshot file.
// Mirrors the path convention in createExternalSnapshot:
// {storageBase}/{domainUUID}/snapshots/{snapName}/{disk}.qcow2
func SnapshotFilePath(storageBase, domainUUID, snapName, disk string) string {
	return filepath.Join(storageBase, domainUUID, "snapshots", snapName, disk+".qcow2")
}

// UploadToPresignedURL uploads the file at filePath via an S3-compatible presigned PUT URL.
func UploadToPresignedURL(ctx context.Context, filePath, presignedURL string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open snapshot file %s: %w", filePath, err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat snapshot file: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, presignedURL, f)
	if err != nil {
		return fmt.Errorf("failed to build upload request: %w", err)
	}
	req.ContentLength = info.Size()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload snapshot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("upload returned status %d", resp.StatusCode)
	}

	return nil
}
