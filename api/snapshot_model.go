package api

// Snapshot API structures
type SnapshotRequest struct {
	UUID string `json:"UUID"`
	Name string `json:"Name,omitempty"`
}
