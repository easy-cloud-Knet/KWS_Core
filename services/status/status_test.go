package status

import (
	"fmt"
	"testing"

	"libvirt.org/go/libvirt"
)

type mockConnect struct {
	domains []libvirt.Domain
	err     error
}

func (m *mockConnect) ListAllDomains(_ libvirt.ConnectListAllDomainsFlags) ([]libvirt.Domain, error) {
	return m.domains, m.err
}

func TestListAllDomainStates_ConnectError(t *testing.T) {
	mock := &mockConnect{err: fmt.Errorf("connect error")}
	_, err := ListAllDomainStates(mock)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListAllDomainStates_Empty(t *testing.T) {
	mock := &mockConnect{domains: []libvirt.Domain{}}
	result, err := ListAllDomainStates(mock)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
}
