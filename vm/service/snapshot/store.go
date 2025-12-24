package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
)

// SnapshotMeta holds minimal metadata for a snapshot.
type SnapshotMeta struct {
	Name         string    `json:"name"`
	SnapshotUUID string    `json:"snapshot_uuid,omitempty"`
	DomainUUID   string    `json:"domain_uuid,omitempty"`
	LibvirtName  string    `json:"libvirt_name,omitempty"`
	XMLDesc      string    `json:"xml_desc,omitempty"`
	Path         string    `json:"path,omitempty"`
	SizeBytes    int64     `json:"size_bytes,omitempty"`
	External     bool      `json:"external,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	Notes        string    `json:"notes,omitempty"`
}

func snapshotsDirForDomain(domain *domCon.Domain) string {
	// attempt to read UUID from domain object
	if domain == nil || domain.Domain == nil {
		return ""
	}
	uuid, err := domain.Domain.GetUUIDString()
	if err != nil {
		return ""
	}
	base := "/var/lib/kws/snapshots"
	return filepath.Join(base, uuid)
}

func ensureSnapshotsDir(path string) error {
	if path == "" {
		return fmt.Errorf("invalid path")
	}
	return os.MkdirAll(path, 0755)
}

// SaveSnapshotMetaForDomain appends a snapshot metadata entry for the domain.
func SaveSnapshotMetaForDomain(domain *domCon.Domain, meta SnapshotMeta) error {
	dir := snapshotsDirForDomain(domain)
	if dir == "" {
		return fmt.Errorf("could not determine snapshots dir")
	}
	if err := ensureSnapshotsDir(dir); err != nil {
		return err
	}

	file := filepath.Join(dir, "snapshots.json")
	var list []SnapshotMeta
	if _, err := os.Stat(file); err == nil {
		b, err := os.ReadFile(file)
		if err == nil {
			_ = json.Unmarshal(b, &list)
		}
	}

	// fill defaults if missing
	if meta.DomainUUID == "" && domain != nil && domain.Domain != nil {
		if u, err := domain.Domain.GetUUIDString(); err == nil {
			meta.DomainUUID = u
		}
	}
	if meta.CreatedAt.IsZero() {
		meta.CreatedAt = time.Now()
	}

	list = append(list, meta)
	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, b, 0644)
}

// LoadSnapshotMetaForDomain reads stored metadata list (best-effort).
func LoadSnapshotMetaForDomain(domain *domCon.Domain) ([]SnapshotMeta, error) {
	dir := snapshotsDirForDomain(domain)
	if dir == "" {
		return nil, fmt.Errorf("could not determine snapshots dir")
	}
	file := filepath.Join(dir, "snapshots.json")
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var list []SnapshotMeta
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}
	return list, nil
}
