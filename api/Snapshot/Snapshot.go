package snapshot

import (
	"fmt"
	"net/http"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	httputil "github.com/easy-cloud-Knet/KWS_Core/pkg/httputil"
	externalsnapshot "github.com/easy-cloud-Knet/KWS_Core/services/snapshot/external_snap"
	internalsnapshot "github.com/easy-cloud-Knet/KWS_Core/services/snapshot/internal_snap"
	"go.uber.org/zap"
)

func (h *Handler) CreateSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := httputil.ResponseGen[string]("Create Snapshot")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		h.Logger.Error("snapshot create decode failed", zap.Error(err))
		return
	}

	if param.Name == "" {
		resp.ResponseWriteErr(w, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot name required")), http.StatusBadRequest)
		return
	}

	h.Logger.Info("snapshot create start", zap.String("uuid", param.UUID), zap.String("snapshot_name", param.Name))

	dom, err := h.DomainControl.GetDomain(param.UUID)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("snapshot create failed - domain not found", zap.String("uuid", param.UUID), zap.Error(err))
		return
	}

	snapName, err := internalsnapshot.CreateSnapshot(dom, param.Name, &internalsnapshot.SnapshotOptions{
		Description: param.Description,
		Quiesce:     param.Quiesce,
	})
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("snapshot create failed", zap.String("uuid", param.UUID), zap.String("snapshot_name", param.Name), zap.Error(err))
		return
	}

	h.Logger.Info("snapshot create success", zap.String("uuid", param.UUID), zap.String("snapshot_name", snapName))
	resp.ResponseWriteOK(w, &snapName)
}

func (h *Handler) CreateExternalSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &ExternalSnapshotRequest{}
	resp := httputil.ResponseGen[ExternalSnapshotResponse]("Create External Snapshot")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		h.Logger.Error("external snapshot create decode failed", zap.Error(err))
		return
	}

	if param.Name == "" {
		resp.ResponseWriteErr(w, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot name required")), http.StatusBadRequest)
		return
	}

	h.Logger.Info("external snapshot create start", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name))

	dom, err := h.DomainControl.GetDomain(param.UUID)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("external snapshot create failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	snapName, err := externalsnapshot.CreateExternalSnapshot(dom, param.Name, &externalsnapshot.ExternalSnapshotOptions{
		BaseDir:     param.BaseDir,
		Description: param.Description,
		Quiesce:     param.Quiesce,
		Live:        param.Live,
	})
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("external snapshot create failed", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name), zap.Error(err))
		return
	}

	h.Logger.Info("external snapshot create success", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", snapName))
	resp.ResponseWriteOK(w, &ExternalSnapshotResponse{
		UUID:     param.UUID,
		SnapName: snapName,
	})
}

func (h *Handler) ListExternalSnapshots(w http.ResponseWriter, r *http.Request) {
	param := &ExternalSnapshotRequest{}
	resp := httputil.ResponseGen[ExternalSnapshotListResponse]("List External Snapshots")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		h.Logger.Error("external snapshot list decode failed", zap.Error(err))
		return
	}

	h.Logger.Info("external snapshot list start", zap.String("domain_uuid", param.UUID))

	dom, err := h.DomainControl.GetDomain(param.UUID)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("external snapshot list failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	names, err := externalsnapshot.ListExternalSnapshots(dom)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("external snapshot list failed", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	h.Logger.Info("external snapshot list success", zap.String("domain_uuid", param.UUID), zap.Int("snapshot_count", len(names)))
	resp.ResponseWriteOK(w, &ExternalSnapshotListResponse{
		UUID:      param.UUID,
		SnapNames: names,
	})
}

func (h *Handler) RevertExternalSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &ExternalSnapshotRequest{}
	resp := httputil.ResponseGen[ExternalSnapshotResponse]("Revert External Snapshot")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		h.Logger.Error("external snapshot revert decode failed", zap.Error(err))
		return
	}

	if param.Name == "" {
		resp.ResponseWriteErr(w, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot name required")), http.StatusBadRequest)
		return
	}

	h.Logger.Info("external snapshot revert start", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name))

	dom, err := h.DomainControl.GetDomain(param.UUID)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("external snapshot revert failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	if err := externalsnapshot.RevertExternalSnapshot(dom, param.Name); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("external snapshot revert failed", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name), zap.Error(err))
		return
	}

	h.Logger.Info("external snapshot revert success", zap.String("domain_uuid", param.UUID), zap.String("snapshot_name", param.Name))
	resp.ResponseWriteOK(w, &ExternalSnapshotResponse{
		UUID:     param.UUID,
		SnapName: param.Name,
	})
}

func (h *Handler) MergeExternalSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &ExternalSnapshotMergeRequest{}
	resp := httputil.ResponseGen[ExternalSnapshotMergeResponse]("Merge External Snapshot")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		h.Logger.Error("external snapshot merge decode failed", zap.Error(err))
		return
	}

	h.Logger.Info("external snapshot merge start", zap.String("domain_uuid", param.UUID), zap.String("disk", param.Disk))

	dom, err := h.DomainControl.GetDomain(param.UUID)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("external snapshot merge failed - domain not found", zap.String("domain_uuid", param.UUID), zap.Error(err))
		return
	}

	mergedDisks, err := externalsnapshot.MergeExternalSnapshot(dom, param.Disk)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("external snapshot merge failed", zap.String("domain_uuid", param.UUID), zap.String("disk", param.Disk), zap.Error(err))
		return
	}

	h.Logger.Info("external snapshot merge success", zap.String("domain_uuid", param.UUID), zap.String("disk", param.Disk), zap.Int("merged_disk_count", len(mergedDisks)))
	resp.ResponseWriteOK(w, &ExternalSnapshotMergeResponse{
		UUID:        param.UUID,
		MergedDisks: mergedDisks,
	})
}

func (h *Handler) ListSnapshots(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := httputil.ResponseGen[[]string]("List Snapshots")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		h.Logger.Error("snapshot list decode failed", zap.Error(err))
		return
	}

	h.Logger.Info("snapshot list start", zap.String("uuid", param.UUID))

	dom, err := h.DomainControl.GetDomain(param.UUID)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("snapshot list failed - domain not found", zap.String("uuid", param.UUID), zap.Error(err))
		return
	}

	names, err := internalsnapshot.ListSnapshots(dom)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("snapshot list failed", zap.String("uuid", param.UUID), zap.Error(err))
		return
	}

	h.Logger.Info("snapshot list success", zap.String("uuid", param.UUID), zap.Int("snapshot_count", len(names)))
	resp.ResponseWriteOK(w, &names)
}

func (h *Handler) RevertSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := httputil.ResponseGen[any]("Revert Snapshot")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		h.Logger.Error("snapshot revert decode failed", zap.Error(err))
		return
	}

	if param.Name == "" {
		resp.ResponseWriteErr(w, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot name required")), http.StatusBadRequest)
		return
	}

	h.Logger.Info("snapshot revert start", zap.String("uuid", param.UUID), zap.String("snapshot_name", param.Name))

	dom, err := h.DomainControl.GetDomain(param.UUID)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("snapshot revert failed - domain not found", zap.String("uuid", param.UUID), zap.Error(err))
		return
	}

	if err := internalsnapshot.RevertToSnapshot(dom, param.Name); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("snapshot revert failed", zap.String("uuid", param.UUID), zap.String("snapshot_name", param.Name), zap.Error(err))
		return
	}

	h.Logger.Info("snapshot revert success", zap.String("uuid", param.UUID), zap.String("snapshot_name", param.Name))
	resp.ResponseWriteOK(w, nil)
}

func (h *Handler) DeleteSnapshot(w http.ResponseWriter, r *http.Request) {
	param := &SnapshotRequest{}
	resp := httputil.ResponseGen[any]("Delete Snapshot")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		h.Logger.Error("snapshot delete decode failed", zap.Error(err))
		return
	}

	if param.Name == "" {
		resp.ResponseWriteErr(w, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot name required")), http.StatusBadRequest)
		return
	}

	h.Logger.Info("snapshot delete start", zap.String("uuid", param.UUID), zap.String("snapshot_name", param.Name))

	dom, err := h.DomainControl.GetDomain(param.UUID)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("snapshot delete failed - domain not found", zap.String("uuid", param.UUID), zap.Error(err))
		return
	}

	if err := internalsnapshot.DeleteSnapshot(dom, param.Name); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		h.Logger.Error("snapshot delete failed", zap.String("uuid", param.UUID), zap.String("snapshot_name", param.Name), zap.Error(err))
		return
	}

	h.Logger.Info("snapshot delete success", zap.String("uuid", param.UUID), zap.String("snapshot_name", param.Name))
	resp.ResponseWriteOK(w, nil)
}
