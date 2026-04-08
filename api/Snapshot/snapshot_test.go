package snapshot

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	testutil "github.com/easy-cloud-Knet/KWS_Core/test"
	"go.uber.org/zap"
)

type mockDomainController struct {
	getDomainFn func(uuid string) (*domCon.Domain, error)
}

func (m *mockDomainController) GetDomain(uuid string) (*domCon.Domain, error) {
	return m.getDomainFn(uuid)
}

func newTestHandler(dc DomainController) *Handler {
	return &Handler{
		DomainControl: dc,
		Logger:        zap.NewNop(),
	}
}

func domainErrMock() *mockDomainController {
	return &mockDomainController{
		getDomainFn: func(uuid string) (*domCon.Domain, error) {
			return nil, virerr.ErrorGen(virerr.DomainSearchError, fmt.Errorf("domain not found"))
		},
	}
}

// CreateSnapshot

func TestCreateSnapshot_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodPost, "/CreateSnapshot", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.CreateSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateSnapshot_EmptyName(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := testutil.MakeRequest(t, SnapshotRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.CreateSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateSnapshot_GetDomainError(t *testing.T) {
	h := newTestHandler(domainErrMock())
	r := testutil.MakeRequest(t, SnapshotRequest{UUID: "test-uuid", Name: "snap1"})
	w := httptest.NewRecorder()

	h.CreateSnapshot(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// CreateExternalSnapshot

func TestCreateExternalSnapshot_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodPost, "/CreateExternalSnapshot", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.CreateExternalSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateExternalSnapshot_EmptyName(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := testutil.MakeRequest(t, ExternalSnapshotRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.CreateExternalSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateExternalSnapshot_GetDomainError(t *testing.T) {
	h := newTestHandler(domainErrMock())
	r := testutil.MakeRequest(t, ExternalSnapshotRequest{UUID: "test-uuid", Name: "snap1"})
	w := httptest.NewRecorder()

	h.CreateExternalSnapshot(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ListSnapshots

func TestListSnapshots_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodGet, "/ListSnapshots", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.ListSnapshots(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestListSnapshots_GetDomainError(t *testing.T) {
	h := newTestHandler(domainErrMock())
	r := testutil.MakeRequest(t, SnapshotRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.ListSnapshots(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ListExternalSnapshots

func TestListExternalSnapshots_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodGet, "/ListExternalSnapshots", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.ListExternalSnapshots(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestListExternalSnapshots_GetDomainError(t *testing.T) {
	h := newTestHandler(domainErrMock())
	r := testutil.MakeRequest(t, ExternalSnapshotRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.ListExternalSnapshots(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// RevertSnapshot

func TestRevertSnapshot_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodPost, "/RevertSnapshot", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.RevertSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRevertSnapshot_EmptyName(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := testutil.MakeRequest(t, SnapshotRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.RevertSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRevertSnapshot_GetDomainError(t *testing.T) {
	h := newTestHandler(domainErrMock())
	r := testutil.MakeRequest(t, SnapshotRequest{UUID: "test-uuid", Name: "snap1"})
	w := httptest.NewRecorder()

	h.RevertSnapshot(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// RevertExternalSnapshot

func TestRevertExternalSnapshot_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodPost, "/RevertExternalSnapshot", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.RevertExternalSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRevertExternalSnapshot_EmptyName(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := testutil.MakeRequest(t, ExternalSnapshotRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.RevertExternalSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRevertExternalSnapshot_GetDomainError(t *testing.T) {
	h := newTestHandler(domainErrMock())
	r := testutil.MakeRequest(t, ExternalSnapshotRequest{UUID: "test-uuid", Name: "snap1"})
	w := httptest.NewRecorder()

	h.RevertExternalSnapshot(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// MergeExternalSnapshot

func TestMergeExternalSnapshot_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodPost, "/MergeExternalSnapshot", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.MergeExternalSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMergeExternalSnapshot_GetDomainError(t *testing.T) {
	h := newTestHandler(domainErrMock())
	r := testutil.MakeRequest(t, ExternalSnapshotMergeRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.MergeExternalSnapshot(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// DeleteSnapshot

func TestDeleteSnapshot_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodPost, "/DeleteSnapshot", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.DeleteSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteSnapshot_EmptyName(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := testutil.MakeRequest(t, SnapshotRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.DeleteSnapshot(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteSnapshot_GetDomainError(t *testing.T) {
	h := newTestHandler(domainErrMock())
	r := testutil.MakeRequest(t, SnapshotRequest{UUID: "test-uuid", Name: "snap1"})
	w := httptest.NewRecorder()

	h.DeleteSnapshot(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
