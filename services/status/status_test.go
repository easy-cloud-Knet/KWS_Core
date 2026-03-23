package status

import (
	"errors"
	"fmt"
	"testing"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"libvirt.org/go/libvirt"
)

// mockDomain implements Domain interface
type mockDomain struct {
	infoResult      *libvirt.DomainInfo
	infoErr         error
	stateResult     libvirt.DomainState
	stateErr        error
	uuidResult      []byte
	uuidErr         error
	guestInfoResult *libvirt.DomainGuestInfo
	guestInfoErr    error
}

func (m *mockDomain) GetInfo() (*libvirt.DomainInfo, error) {
	return m.infoResult, m.infoErr
}
func (m *mockDomain) GetState() (libvirt.DomainState, int, error) {
	return m.stateResult, 0, m.stateErr
}
func (m *mockDomain) GetUUID() ([]byte, error) {
	return m.uuidResult, m.uuidErr
}
func (m *mockDomain) GetGuestInfo(types libvirt.DomainGuestInfoTypes, flags uint32) (*libvirt.DomainGuestInfo, error) {
	return m.guestInfoResult, m.guestInfoErr
}

// mockConnect implements Connect interface
type mockConnect struct {
	domains []libvirt.Domain
	err     error
}

func (m *mockConnect) ListAllDomains(flags libvirt.ConnectListAllDomainsFlags) ([]libvirt.Domain, error) {
	return m.domains, m.err
}

// valid 16-byte UUID for testing
var testUUIDBytes = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

func TestDomainInfo_GetInfo_Success(t *testing.T) {
	mock := &mockDomain{
		infoResult: &libvirt.DomainInfo{
			State:     1,
			MaxMem:    2048,
			Memory:    1024,
			NrVirtCpu: 2,
			CpuTime:   100,
		},
	}
	di := &DomainInfo{}
	if err := di.GetInfo(mock); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if di.MaxMem != 2048 || di.NrVirtCpu != 2 {
		t.Errorf("fields not populated correctly: %+v", di)
	}
}

func TestDomainInfo_GetInfo_Error(t *testing.T) {
	mock := &mockDomain{infoErr: fmt.Errorf("libvirt error")}
	err := (&DomainInfo{}).GetInfo(mock)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.DomainStatusError) {
		t.Errorf("expected DomainStatusError, got %v", err)
	}
}

func TestDomainState_GetInfo_Success(t *testing.T) {
	mock := &mockDomain{
		stateResult:     libvirt.DomainState(1),
		uuidResult:      testUUIDBytes,
		guestInfoResult: &libvirt.DomainGuestInfo{},
	}
	ds := &DomainState{}
	if err := ds.GetInfo(mock); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if ds.UUID == "" {
		t.Error("UUID not populated")
	}
}

func TestDomainState_GetInfo_StateError(t *testing.T) {
	mock := &mockDomain{stateErr: fmt.Errorf("state error")}
	err := (&DomainState{}).GetInfo(mock)
	if !errors.Is(err, virerr.DomainStatusError) {
		t.Errorf("expected DomainStatusError, got %v", err)
	}
}

func TestDomainState_GetInfo_UUIDError(t *testing.T) {
	mock := &mockDomain{
		stateResult: libvirt.DomainState(1),
		uuidErr:     fmt.Errorf("uuid error"),
	}
	err := (&DomainState{}).GetInfo(mock)
	if !errors.Is(err, virerr.InvalidUUID) {
		t.Errorf("expected InvalidUUID, got %v", err)
	}
}

func TestDomainState_GetInfo_GuestInfoError(t *testing.T) {
	mock := &mockDomain{
		stateResult:  libvirt.DomainState(1),
		uuidResult:   testUUIDBytes,
		guestInfoErr: fmt.Errorf("guest info error"),
	}
	err := (&DomainState{}).GetInfo(mock)
	if !errors.Is(err, virerr.DomainStatusError) {
		t.Errorf("expected DomainStatusError, got %v", err)
	}
}

func TestAllInstInfo_GetAllinstInfo_ListError(t *testing.T) {
	mock := &mockConnect{err: fmt.Errorf("libvirt error")}
	err := (&AllInstInfo{}).GetAllinstInfo(mock)
	if !errors.Is(err, virerr.HostStatusError) {
		t.Errorf("expected HostStatusError, got %v", err)
	}
}

func TestAllInstInfo_GetAllinstInfo_EmptyList(t *testing.T) {
	mock := &mockConnect{domains: []libvirt.Domain{}}
	aii := &AllInstInfo{}
	if err := aii.GetAllinstInfo(mock); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if aii.Totalmaxmem != 0 || aii.TotalVCpu != 0 {
		t.Errorf("expected zeros, got %+v", aii)
	}
}

func TestDataTypeRouter_ValidTypes(t *testing.T) {
	cases := []DomainDataType{DomState, BasicInfo, GuestInfoUser, GuestInfoOS, GuestInfoFS, GuestInfoDisk}
	for _, c := range cases {
		if _, err := DataTypeRouter(c); err != nil {
			t.Errorf("DataTypeRouter(%d) returned error: %v", c, err)
		}
	}
}

func TestDataTypeRouter_InvalidType(t *testing.T) {
	_, err := DataTypeRouter(DomainDataType(99))
	if !errors.Is(err, virerr.InvalidParameter) {
		t.Errorf("expected InvalidParameter, got %v", err)
	}
}
