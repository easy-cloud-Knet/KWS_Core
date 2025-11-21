package domStatus

import (
	"runtime"
	"sync/atomic"
)

func (dls *DomainListStatus) UpdateCPUTotal() {
	totalCPU := runtime.NumCPU()
	dls.VCPUTotal = int32(totalCPU)
}

func (dls *DomainListStatus) UpdateActiveCPU(vcpu int) error {
	atomic.AddInt32(&dls.VcpuAllocated, int32(vcpu))
	return nil
}

func (dls *DomainListStatus) UpdateSleepingCPU(vcpu int) error {
	atomic.AddInt32(&dls.vcpuSleeping, int32(vcpu))
	return nil
}

func (dls *DomainListStatus) GetVCPUStatus() VCPUStatus {
	status := VCPUStatus{
		Total: int(dls.VCPUTotal),
		Allocated: int(dls.VcpuAllocated),
		Sleeping: int(dls.vcpuSleeping),
		Idle: int(dls.VCPUTotal - dls.VcpuAllocated),
	}
	return status
}
