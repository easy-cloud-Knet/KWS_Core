package domListStatus

import (
	"fmt"
	"testing"

	instatus "github.com/easy-cloud-Knet/KWS_Core/internal/status"
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	"go.uber.org/zap"
)

type mockDataDog struct {
	result map[vmtypes.SourceType]int
	err    error
}

func (m *mockDataDog) RetrieveStatus(_ instatus.Domain, _ map[vmtypes.SourceType]int, _ *zap.Logger) (map[vmtypes.SourceType]int, error) {
	return m.result, m.err
}

var nopLogger = zap.NewNop()

func TestUpdateFromDomain_ActiveIncreasesAllocated(t *testing.T) {
	dls := &DomainListStatus{}
	mock := &mockDataDog{result: map[vmtypes.SourceType]int{vmtypes.CPU: 4}}

	if err := dls.UpdateFromDomain(mock, nil, true, map[vmtypes.SourceType]int{vmtypes.CPU: 0}, nopLogger); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dls.VcpuAllocated != 4 {
		t.Errorf("expected VcpuAllocated=4, got %d", dls.VcpuAllocated)
	}
	if dls.VcpuSleeping != 0 {
		t.Errorf("expected VcpuSleeping=0 for active domain, got %d", dls.VcpuSleeping)
	}
}

func TestUpdateFromDomain_InactiveIncreasesBoth(t *testing.T) {
	dls := &DomainListStatus{}
	mock := &mockDataDog{result: map[vmtypes.SourceType]int{vmtypes.CPU: 4}}

	if err := dls.UpdateFromDomain(mock, nil, false, map[vmtypes.SourceType]int{vmtypes.CPU: 0}, nopLogger); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dls.VcpuAllocated != 4 {
		t.Errorf("expected VcpuAllocated=4, got %d", dls.VcpuAllocated)
	}
	if dls.VcpuSleeping != 4 {
		t.Errorf("expected VcpuSleeping=4 for inactive domain, got %d", dls.VcpuSleeping)
	}
}

func TestUpdateFromDomain_DataDogError(t *testing.T) {
	dls := &DomainListStatus{}
	mock := &mockDataDog{err: fmt.Errorf("retrieval error")}

	err := dls.UpdateFromDomain(mock, nil, true, map[vmtypes.SourceType]int{vmtypes.CPU: 0}, nopLogger)
	if err == nil {
		t.Error("expected error from DataDog, got nil")
	}
	if dls.VcpuAllocated != 0 {
		t.Errorf("expected no CPU change on error, got %d", dls.VcpuAllocated)
	}
}
