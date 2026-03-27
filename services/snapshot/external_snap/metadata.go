package external

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core/internal/config"
)

const defaultSnapshotRoot = "/var/lib/kws"
const externalSnapshotMetadataFile = ".kws_external_snapshots.json"

type externalSnapshotMetadataStore struct {
	Version   int                             `json:"version"`
	Snapshots []externalSnapshotMetadataEntry `json:"snapshots"`
}

type externalSnapshotMetadataEntry struct {
	Name      string                     `json:"name"`
	CreatedAt string                     `json:"created_at"`
	Disks     []externalSnapshotDiskMeta `json:"disks"`
}

type externalSnapshotDiskMeta struct {
	TargetDev    string `json:"target_dev"`
	SnapshotFile string `json:"snapshot_file"`
	BackingFile  string `json:"backing_file"`
}

func resolveSnapshotRoot(opts *ExternalSnapshotOptions) (string, error) {
	if opts == nil || opts.BaseDir == "" {
		return config.StorageBase, nil
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

func resolveDomainUUID(domain *domCon.Domain) (string, error) {
	if domain == nil || domain.Domain == nil {
		return "", fmt.Errorf("nil domain")
	}

	uuid, err := domain.Domain.GetUUIDString()
	if err != nil {
		return "", fmt.Errorf("failed to get domain uuid: %w", err)
	}

	return uuid, nil
}

func resolveExternalSnapshotMetadataPath(domain *domCon.Domain) (string, error) {
	uuid, err := resolveDomainUUID(domain)
	if err != nil {
		return "", err
	}

	return filepath.Join(defaultSnapshotRoot, uuid, externalSnapshotMetadataFile), nil
}

func loadExternalSnapshotMetadata(domain *domCon.Domain) (*externalSnapshotMetadataStore, string, error) {
	metadataPath, err := resolveExternalSnapshotMetadataPath(domain)
	if err != nil {
		return nil, "", err
	}

	store := &externalSnapshotMetadataStore{Version: 1, Snapshots: []externalSnapshotMetadataEntry{}}
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return store, metadataPath, nil
		}
		return nil, "", fmt.Errorf("failed to read external snapshot metadata: %w", err)
	}

	if err := json.Unmarshal(data, store); err != nil {
		return nil, "", fmt.Errorf("failed to parse external snapshot metadata: %w", err)
	}
	if store.Version == 0 {
		store.Version = 1
	}

	return store, metadataPath, nil
}

func writeExternalSnapshotMetadata(metadataPath string, store *externalSnapshotMetadataStore) error {
	if store == nil {
		return fmt.Errorf("nil metadata store")
	}

	if err := os.MkdirAll(filepath.Dir(metadataPath), 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	payload, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal external snapshot metadata: %w", err)
	}

	tmpPath := metadataPath + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0644); err != nil {
		return fmt.Errorf("failed to write external snapshot metadata: %w", err)
	}

	if err := os.Rename(tmpPath, metadataPath); err != nil {
		return fmt.Errorf("failed to replace external snapshot metadata: %w", err)
	}

	return nil
}

func appendExternalSnapshotMetadata(domain *domCon.Domain, snapName string, diskMetas []externalSnapshotDiskMeta) error {
	if snapName == "" {
		return fmt.Errorf("snapshot name required")
	}
	if len(diskMetas) == 0 {
		return fmt.Errorf("no disk metadata to record")
	}

	store, metadataPath, err := loadExternalSnapshotMetadata(domain)
	if err != nil {
		return err
	}

	store.Snapshots = append(store.Snapshots, externalSnapshotMetadataEntry{
		Name:      snapName,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Disks:     diskMetas,
	})

	return writeExternalSnapshotMetadata(metadataPath, store)
}

func findDiskBackingFromMetadata(store *externalSnapshotMetadataStore, targetDev string, currentSource string) (string, bool) {
	if store == nil {
		return "", false
	}

	for i := len(store.Snapshots) - 1; i >= 0; i-- {
		for _, d := range store.Snapshots[i].Disks {
			if d.TargetDev == targetDev && d.SnapshotFile == currentSource {
				return d.BackingFile, true
			}
		}
	}

	return "", false
}

func pruneMergedDisksFromMetadata(domain *domCon.Domain, merged map[string]struct{}) error {
	if len(merged) == 0 {
		return nil
	}

	store, metadataPath, err := loadExternalSnapshotMetadata(domain)
	if err != nil {
		return err
	}

	prunedSnapshots := make([]externalSnapshotMetadataEntry, 0, len(store.Snapshots))
	for _, snap := range store.Snapshots {
		remainingDisks := make([]externalSnapshotDiskMeta, 0, len(snap.Disks))
		for _, d := range snap.Disks {
			key := d.TargetDev + "|" + d.SnapshotFile
			if _, ok := merged[key]; ok {
				continue
			}
			remainingDisks = append(remainingDisks, d)
		}
		if len(remainingDisks) == 0 {
			continue
		}
		snap.Disks = remainingDisks
		prunedSnapshots = append(prunedSnapshots, snap)
	}

	store.Snapshots = prunedSnapshots
	return writeExternalSnapshotMetadata(metadataPath, store)
}

func deleteSnapshotMetadataByName(domain *domCon.Domain, snapName string) error {
	if snapName == "" {
		return fmt.Errorf("snapshot name required")
	}

	store, metadataPath, err := loadExternalSnapshotMetadata(domain)
	if err != nil {
		return err
	}

	filtered := make([]externalSnapshotMetadataEntry, 0, len(store.Snapshots))
	for _, snap := range store.Snapshots {
		if snap.Name != snapName {
			filtered = append(filtered, snap)
		}
	}

	store.Snapshots = filtered
	return writeExternalSnapshotMetadata(metadataPath, store)
}
