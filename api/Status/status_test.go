package status

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domainList_status"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	svcstatus "github.com/easy-cloud-Knet/KWS_Core/services/status"
	testutil "github.com/easy-cloud-Knet/KWS_Core/test"
	"go.uber.org/zap"
)

// ─── mock ───────────────────────────────────────────────────────────────────

type mockDomainController struct {
	getDomainFn          func(uuid string) (*domCon.Domain, error)
	getAllUUIDsFn         func() []string
	getDomainListStatusFn func() *domStatus.DomainListStatus
}

func (m *mockDomainController) GetDomain(uuid string) (*domCon.Domain, error) {
	return m.getDomainFn(uuid)
}

func (m *mockDomainController) GetAllUUIDs() []string {
	return m.getAllUUIDsFn()
}

func (m *mockDomainController) GetDomainListStatus() *domStatus.DomainListStatus {
	return m.getDomainListStatusFn()
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
			return nil, virerr.ErrorGen(virerr.DomainSearchError, nil)
		},
	}
}

// ─── ReturnStatusUUID ────────────────────────────────────────────────────────

func TestReturnStatusUUID_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodGet, "/getStatusUUID", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.ReturnStatusUUID(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestReturnStatusUUID_InvalidDataType(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	// DomainDataType 99 는 DataTypeRouter 에서 에러 반환
	r := testutil.MakeRequest(t, DomainStatusRequest{UUID: "test-uuid", DataType: svcstatus.DomainDataType(99)})
	w := httptest.NewRecorder()

	h.ReturnStatusUUID(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestReturnStatusUUID_GetDomainError(t *testing.T) {
	h := newTestHandler(domainErrMock())
	r := testutil.MakeRequest(t, DomainStatusRequest{UUID: "test-uuid", DataType: svcstatus.DomState})
	w := httptest.NewRecorder()

	h.ReturnStatusUUID(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ─── ReturnStatusHost ────────────────────────────────────────────────────────

func TestReturnStatusHost_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodGet, "/getStatusHost", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.ReturnStatusHost(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestReturnStatusHost_InvalidHostDataType(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	// HostDataType 99 는 HostDataTypeRouter 에서 에러 반환
	r := testutil.MakeRequest(t, HostStatusRequest{HostDataType: svcstatus.HostDataType(99)})
	w := httptest.NewRecorder()

	h.ReturnStatusHost(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ─── ReturnInstAllInfo ───────────────────────────────────────────────────────

func TestReturnInstAllInfo_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodGet, "/getInstAllInfo", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.ReturnInstAllInfo(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestReturnInstAllInfo_InvalidInstDataType(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	// InstDataType 99 는 InstDataTypeRouter 에서 에러 반환
	r := testutil.MakeRequest(t, InstInfoRequest{InstDataType: svcstatus.InstDataType(99)})
	w := httptest.NewRecorder()

	h.ReturnInstAllInfo(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// ─── ReturnAllUUIDs ──────────────────────────────────────────────────────────

func TestReturnAllUUIDs_EmptyList(t *testing.T) {
	dc := &mockDomainController{
		getAllUUIDsFn: func() []string { return []string{} },
	}
	h := newTestHandler(dc)
	r := httptest.NewRequest(http.MethodGet, "/getAllUUIDs", nil)
	w := httptest.NewRecorder()

	h.ReturnAllUUIDs(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestReturnAllUUIDs_WithUUIDs(t *testing.T) {
	dc := &mockDomainController{
		getAllUUIDsFn: func() []string {
			return []string{"uuid-1", "uuid-2", "uuid-3"}
		},
	}
	h := newTestHandler(dc)
	r := httptest.NewRequest(http.MethodGet, "/getAllUUIDs", nil)
	w := httptest.NewRecorder()

	h.ReturnAllUUIDs(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}
}
