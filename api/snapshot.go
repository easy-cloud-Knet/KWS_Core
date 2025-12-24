package api

import (
	"fmt"
	"net/http"
)

// Snapshot API structures
type SnapshotRequest struct {
	UUID string `json:"UUID"`
	Name string `json:"Name,omitempty"`
}

// CreateSnapshot creates a snapshot for the specified domain UUID
func (i *InstHandler) CreateSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := ResponseGen[string]("Create Snapshot")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error(err.Error())
		return
	}

	if i.SnapshotManager == nil {
		resp.ResponseWriteErr(w, fmt.Errorf("snapshot manager not initialized"), http.StatusInternalServerError)
		i.Logger.Error("snapshot manager not initialized")
		return
	}

	name := param.Name
	if name == "" {
		name = param.UUID + "-snap"
	}

	snapName, err := i.SnapshotManager.Create(param.UUID, name)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error(err.Error())
		return
	}

	resp.ResponseWriteOK(w, &snapName)
}

// ListSnapshots returns all snapshot names for the specified domain UUID
func (i *InstHandler) ListSnapshots(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := ResponseGen[[]string]("List Snapshots")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error(err.Error())
		return
	}

	if i.SnapshotManager == nil {
		resp.ResponseWriteErr(w, fmt.Errorf("snapshot manager not initialized"), http.StatusInternalServerError)
		i.Logger.Error("snapshot manager not initialized")
		return
	}

	names, err := i.SnapshotManager.List(param.UUID)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error(err.Error())
		return
	}

	resp.ResponseWriteOK(w, &names)
}

// RevertSnapshot reverts the domain to a named snapshot
func (i *InstHandler) RevertSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := ResponseGen[any]("Revert Snapshot")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error(err.Error())
		return
	}

	if i.SnapshotManager == nil {
		resp.ResponseWriteErr(w, fmt.Errorf("snapshot manager not initialized"), http.StatusInternalServerError)
		i.Logger.Error("snapshot manager not initialized")
		return
	}

	if param.Name == "" {
		resp.ResponseWriteErr(w, fmt.Errorf("snapshot name required"), http.StatusBadRequest)
		return
	}

	if err := i.SnapshotManager.Revert(param.UUID, param.Name); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error(err.Error())
		return
	}

	resp.ResponseWriteOK(w, nil)
}
