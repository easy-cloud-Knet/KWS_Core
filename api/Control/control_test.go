package control

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
	getDomainFn   func(uuid string) (*domCon.Domain, error)
	sleepDomainFn func(domain *domCon.Domain, logger *zap.Logger) error
	removeDomainFn func(domain *domCon.Domain, uuid string, logger *zap.Logger) error
}

func (m *mockDomainController) GetDomain(uuid string) (*domCon.Domain, error) {
	return m.getDomainFn(uuid)
}
func (m *mockDomainController) SleepDomain(domain *domCon.Domain, logger *zap.Logger) error {
	return m.sleepDomainFn(domain, logger)
}
func (m *mockDomainController) RemoveDomain(domain *domCon.Domain, uuid string, logger *zap.Logger) error {
	return m.removeDomainFn(domain, uuid, logger)
}

func newTestHandler(dc DomainController) *Handler {
	return &Handler{
		DomainControl: dc,
		Logger:        zap.NewNop(),
	}
}

// ForceShutDownVM tests

func TestForceShutDownVM_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodPost, "/forceShutDownUUID", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.ForceShutDownVM(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestForceShutDownVM_GetDomainError(t *testing.T) {
	dc := &mockDomainController{
		getDomainFn: func(uuid string) (*domCon.Domain, error) {
			return nil, virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("domain not found"))
		},
	}
	h := newTestHandler(dc)
	r := testutil.MakeRequest(t, DomainControlRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.ForceShutDownVM(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// DeleteVM tests

func TestDeleteVM_BadRequest(t *testing.T) {
	h := newTestHandler(&mockDomainController{})
	r := httptest.NewRequest(http.MethodPost, "/DeleteVM", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.DeleteVM(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestDeleteVM_GetDomainError(t *testing.T) {
	dc := &mockDomainController{
		getDomainFn: func(uuid string) (*domCon.Domain, error) {
			return nil, virerr.ErrorGen(virerr.DomainSearchError, fmt.Errorf("domain search failed"))
		},
	}
	h := newTestHandler(dc)
	r := testutil.MakeRequest(t, DomainControlRequest{UUID: "test-uuid"})
	w := httptest.NewRecorder()

	h.DeleteVM(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
