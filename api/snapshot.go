package api

import (
	"fmt"
	"net/http"

	snapshotpkg "github.com/easy-cloud-Knet/KWS_Core/vm/service/snapshot"
	"go.uber.org/zap"
)

// Snapshot API structures
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

// CreateSnapshot creates a snapshot for the specified domain UUID
func (i *InstHandler) CreateSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := ResponseGen[string]("Create Snapshot")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot create decode failed", zap.Error(err))
		return
	}

	if i.LibvirtInst == nil || i.DomainControl == nil {
		resp.ResponseWriteErr(w, fmt.Errorf("libvirt not initialized"), http.StatusInternalServerError)
		i.Logger.Error("libvirt not initialized")
		return
	}

	name := param.Name
	if name == "" {
		name = param.UUID + "-snap"
	}

	i.Logger.Info("snapshot create start", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", name))

	dom, err := i.DomainControl.GetDomain(param.UUID, i.LibvirtInst)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot create failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	snapName, err := snapshotpkg.CreateSnapshot(dom, name, &snapshotpkg.SnapshotOptions{
		Description: param.Description,
		Quiesce:     param.Quiesce,
	})
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot create failed", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", name), zap.Error(err))
		return
	}

	i.Logger.Info("snapshot create success", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", snapName))
	resp.ResponseWriteOK(w, &snapName)
}

// CreateExternalSnapshot creates an external snapshot for the specified domain UUID
func (i *InstHandler) CreateExternalSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &ExternalSnapshotRequest{}
	resp := ResponseGen[string]("Create External Snapshot")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("external snapshot create decode failed", zap.Error(err))
		return
	}

	if i.LibvirtInst == nil || i.DomainControl == nil {
		resp.ResponseWriteErr(w, fmt.Errorf("libvirt not initialized"), http.StatusInternalServerError)
		i.Logger.Error("libvirt not initialized")
		return
	}

	name := param.Name
	if name == "" {
		name = param.UUID + "-extsnap"
	}

	i.Logger.Info("external snapshot create start", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", name))

	dom, err := i.DomainControl.GetDomain(param.UUID, i.LibvirtInst)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("external snapshot create failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	snapName, err := snapshotpkg.CreateExternalSnapshot(dom, name, &snapshotpkg.ExternalSnapshotOptions{
		BaseDir:     param.BaseDir,
		Description: param.Description,
		Quiesce:     param.Quiesce,
		Live:        param.Live,
	})
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("external snapshot create failed", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", name), zap.Error(err))
		return
	}

	i.Logger.Info("external snapshot create success", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", snapName))
	resp.ResponseWriteOK(w, &snapName)
}

// ListExternalSnapshots returns external snapshot names for the specified domain UUID
func (i *InstHandler) ListExternalSnapshots(w http.ResponseWriter, r *http.Request) {
	param := &ExternalSnapshotRequest{}
	resp := ResponseGen[[]string]("List External Snapshots")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("external snapshot list decode failed", zap.Error(err))
		return
	}

	if i.LibvirtInst == nil || i.DomainControl == nil {
		resp.ResponseWriteErr(w, fmt.Errorf("libvirt not initialized"), http.StatusInternalServerError)
		i.Logger.Error("libvirt not initialized")
		return
	}

	i.Logger.Info("external snapshot list start", zap.String("domain_uuid", param.UUID))

	dom, err := i.DomainControl.GetDomain(param.UUID, i.LibvirtInst)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("external snapshot list failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	names, err := snapshotpkg.ListExternalSnapshots(dom)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("external snapshot list failed", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	i.Logger.Info("external snapshot list success", zap.String("domain_uuid", param.UUID), zap.Int("snapshot_count", len(names)))
	resp.ResponseWriteOK(w, &names)
}

// RevertExternalSnapshot reverts the domain to a named external snapshot
func (i *InstHandler) RevertExternalSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &ExternalSnapshotRequest{}
	resp := ResponseGen[any]("Revert External Snapshot")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("external snapshot revert decode failed", zap.Error(err))
		return
	}

	if i.LibvirtInst == nil || i.DomainControl == nil {
		resp.ResponseWriteErr(w, fmt.Errorf("libvirt not initialized"), http.StatusInternalServerError)
		i.Logger.Error("libvirt not initialized")
		return
	}

	if param.Name == "" {
		resp.ResponseWriteErr(w, fmt.Errorf("snapshot name required"), http.StatusBadRequest)
		return
	}

	i.Logger.Info("external snapshot revert start", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name))

	dom, err := i.DomainControl.GetDomain(param.UUID, i.LibvirtInst)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("external snapshot revert failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	if err := snapshotpkg.RevertExternalSnapshot(dom, param.Name); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("external snapshot revert failed", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name), zap.Error(err))
		return
	}

	i.Logger.Info("external snapshot revert success", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name))
	resp.ResponseWriteOK(w, nil)
}

// ListSnapshots returns all snapshot names for the specified domain UUID
func (i *InstHandler) ListSnapshots(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := ResponseGen[[]string]("List Snapshots")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot list decode failed", zap.Error(err))
		return
	}

	if i.LibvirtInst == nil || i.DomainControl == nil {
		resp.ResponseWriteErr(w, fmt.Errorf("libvirt not initialized"), http.StatusInternalServerError)
		i.Logger.Error("libvirt not initialized")
		return
	}

	i.Logger.Info("snapshot list start", zap.String("domain_uuid", param.UUID))

	dom, err := i.DomainControl.GetDomain(param.UUID, i.LibvirtInst)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot list failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	names, err := snapshotpkg.ListSnapshots(dom)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot list failed", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	i.Logger.Info("snapshot list success", zap.String("domain_uuid", param.UUID), zap.Int("snapshot_count", len(names)))
	resp.ResponseWriteOK(w, &names)
}

// RevertSnapshot reverts the domain to a named snapshot
func (i *InstHandler) RevertSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := ResponseGen[any]("Revert Snapshot")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot revert decode failed", zap.Error(err))
		return
	}

	if i.LibvirtInst == nil || i.DomainControl == nil {
		resp.ResponseWriteErr(w, fmt.Errorf("libvirt not initialized"), http.StatusInternalServerError)
		i.Logger.Error("libvirt not initialized")
		return
	}

	if param.Name == "" {
		resp.ResponseWriteErr(w, fmt.Errorf("snapshot name required"), http.StatusBadRequest)
		return
	}

	i.Logger.Info("snapshot revert start", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name))

	dom, err := i.DomainControl.GetDomain(param.UUID, i.LibvirtInst)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot revert failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	if err := snapshotpkg.RevertToSnapshot(dom, param.Name); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot revert failed", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name), zap.Error(err))
		return
	}

	i.Logger.Info("snapshot revert success", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name))
	resp.ResponseWriteOK(w, nil)
}

// DeleteSnapshot deletes a snapshot by name
func (i *InstHandler) DeleteSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := ResponseGen[any]("Delete Snapshot")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot delete decode failed", zap.Error(err))
		return
	}

	if i.LibvirtInst == nil || i.DomainControl == nil {
		resp.ResponseWriteErr(w, fmt.Errorf("libvirt not initialized"), http.StatusInternalServerError)
		i.Logger.Error("libvirt not initialized")
		return
	}

	if param.Name == "" {
		resp.ResponseWriteErr(w, fmt.Errorf("snapshot name required"), http.StatusBadRequest)
		return
	}

	i.Logger.Info("snapshot delete start", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name))

	dom, err := i.DomainControl.GetDomain(param.UUID, i.LibvirtInst)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot delete failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	if err := snapshotpkg.DeleteSnapshot(dom, param.Name); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		i.Logger.Error("snapshot delete failed", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name), zap.Error(err))
		return
	}

	i.Logger.Info("snapshot delete success", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name))
	resp.ResponseWriteOK(w, nil)
}
