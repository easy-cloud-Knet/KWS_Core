package termination

import (
	"errors"
	"fmt"
	"testing"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

type mockDomain struct {
	active         bool
	destroyErr     error
	undefineErr    error
	destroyCalled  bool
	undefineCalled bool
}

func (m *mockDomain) IsActive() (bool, error) { return m.active, nil }
func (m *mockDomain) Destroy() error          { m.destroyCalled = true; return m.destroyErr }
func (m *mockDomain) Undefine() error         { m.undefineCalled = true; return m.undefineErr }

func TestTerminateDomain_Running(t *testing.T) {
	mock := &mockDomain{active: true}
	err := DomainTerminatorFactory(mock).TerminateDomain()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if !mock.destroyCalled {
		t.Error("Destroy not called on running domain")
	}
}

func TestTerminateDomain_Stopped(t *testing.T) {
	mock := &mockDomain{active: false}
	err := DomainTerminatorFactory(mock).TerminateDomain()
	if err == nil {
		t.Error("expected error for stopped domain, got nil")
	}
	if !errors.Is(err, virerr.DomainShutdownError) {
		t.Errorf("expected DomainShutdownError, got %v", err)
	}
	if mock.destroyCalled {
		t.Error("Destroy should not be called on stopped domain")
	}
}

func TestTerminateDomain_DestroyError(t *testing.T) {
	mock := &mockDomain{active: true, destroyErr: fmt.Errorf("libvirt error")}
	err := DomainTerminatorFactory(mock).TerminateDomain()
	if err == nil {
		t.Error("expected error from Destroy, got nil")
	}
	if !errors.Is(err, virerr.DomainShutdownError) {
		t.Errorf("expected DomainShutdownError, got %v", err)
	}
}

func TestDeleteDomain_SoftDelete_Running(t *testing.T) {
	mock := &mockDomain{active: true}
	err := DomainDeleterFactory(mock, SoftDelete, "test-uuid").DeleteDomain()
	if err == nil {
		t.Error("expected error for SoftDelete on running domain, got nil")
	}
	if !errors.Is(err, virerr.DeletionDomainError) {
		t.Errorf("expected DeletionDomainError, got %v", err)
	}
	if mock.destroyCalled || mock.undefineCalled {
		t.Error("Destroy/Undefine should not be called")
	}
}

func TestDeleteDomain_HardDelete_Running(t *testing.T) {
	mock := &mockDomain{active: true}
	err := DomainDeleterFactory(mock, HardDelete, "test-uuid").DeleteDomain()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if !mock.destroyCalled {
		t.Error("Destroy not called on HardDelete")
	}
	if !mock.undefineCalled {
		t.Error("Undefine not called after HardDelete")
	}
}

func TestDeleteDomain_SoftDelete_Stopped(t *testing.T) {
	mock := &mockDomain{active: false}
	err := DomainDeleterFactory(mock, SoftDelete, "test-uuid").DeleteDomain()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if mock.destroyCalled {
		t.Error("Destroy should not be called on stopped domain")
	}
	if !mock.undefineCalled {
		t.Error("Undefine not called")
	}
}
