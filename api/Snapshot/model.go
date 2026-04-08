package snapshot

type SnapshotRequest struct {
	UUID        string `json:"UUID"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
	Quiesce     bool   `json:"Quiesce,omitempty"`
}

type ExternalSnapshotRequest struct {
	UUID        string `json:"UUID"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
	BaseDir     string `json:"BaseDir,omitempty"`
	Quiesce     bool   `json:"Quiesce,omitempty"`
	Live        bool   `json:"Live,omitempty"`
}

type ExternalSnapshotResponse struct {
	UUID     string `json:"UUID"`
	SnapName string `json:"SnapName"`
}

type ExternalSnapshotListResponse struct {
	UUID      string   `json:"UUID"`
	SnapNames []string `json:"SnapNames"`
}

type ExternalSnapshotMergeRequest struct {
	UUID string `json:"UUID"`
	Disk string `json:"Disk,omitempty"`
}

type ExternalSnapshotMergeResponse struct {
	UUID        string   `json:"UUID"`
	MergedDisks []string `json:"MergedDisks"`
}
