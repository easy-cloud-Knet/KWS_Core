package create

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
	getDomainFn       func(uuid string) (*domCon.Domain, error)
	addNewDomainFn    func(domain *domCon.Domain, uuid string) error
	bootSleepingCPUFn func(domain *domCon.Domain) error
}

func (m *mockDomainController) GetDomain(uuid string) (*domCon.Domain, error) {
	return m.getDomainFn(uuid)
}
func (m *mockDomainController) AddNewDomain(domain *domCon.Domain, uuid string) error {
	return m.addNewDomainFn(domain, uuid)
}
func (m *mockDomainController) BootSleepingCPU(domain *domCon.Domain) error {
	return m.bootSleepingCPUFn(domain)
}

func newTestHandler(dc DomainController) *Handler {
	return &Handler{
		DomainControl: dc,
		Logger:        zap.NewNop(),
	}
}

// BootVM tests

func TestBootVM_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.BootVM(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestBootVM_GetDomainError(t *testing.T) {
	dc := &mockDomainController{
		getDomainFn: func(uuid string) (*domCon.Domain, error) {
			return nil, virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("domain error"))
		},
	}
	h := newTestHandler(dc)
	r := testutil.MakeRequest(t, DomainBootRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.BootVM(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestBootVM_DomainNotFound(t *testing.T) {
	dc := &mockDomainController{
		getDomainFn: func(uuid string) (*domCon.Domain, error) {
			return nil, nil
		},
	}
	h := newTestHandler(dc)
	r := testutil.MakeRequest(t, DomainBootRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.BootVM(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

// CreateVMFromBase tests

func TestCreateVMFromBase_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.CreateVMFromBase(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateVMFromBase_GetDomainSearchError(t *testing.T) {
	dc := &mockDomainController{
		getDomainFn: func(uuid string) (*domCon.Domain, error) {
			return nil, virerr.ErrorGen(virerr.DomainSearchError, fmt.Errorf("search error"))
		},
	}
	h := newTestHandler(dc)
	r := testutil.MakeRequest(t, CreateVMRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.CreateVMFromBase(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestCreateVMFromBase_DomainAlreadyExists(t *testing.T) {
	existing := &domCon.Domain{}
	dc := &mockDomainController{
		getDomainFn: func(uuid string) (*domCon.Domain, error) {
			return existing, nil
		},
	}
	h := newTestHandler(dc)
	r := testutil.MakeRequest(t, CreateVMRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.CreateVMFromBase(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("expected %d, got %d", http.StatusConflict, w.Code)
	}
}
