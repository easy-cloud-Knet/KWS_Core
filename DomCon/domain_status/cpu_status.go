package domStatus

import (
	"runtime"
	"sync/atomic"
)

func (dls *DomainListStatus) UpdateCPUTotal() {
	totalCPU := runtime.NumCPU()
	dls.VCPUTotal = int32(totalCPU)
}

func (dls *DomainListStatus) AddAllocatedCPU(vcpu int) error {
	atomic.AddInt32(&dls.VcpuAllocated, int32(vcpu))
	return nil
}

func (dls *DomainListStatus) AddSleepingCPU(vcpu int) error {
	atomic.AddInt32(&dls.VcpuSleeping, int32(vcpu))
	return nil
}
func (dls *DomainListStatus) TakeAllocatedCPU(vcpu int) error {
	num := int(dls.VcpuAllocated) 
	
	atomic.SwapInt32(&dls.VcpuAllocated, int32(num-vcpu))
	return nil
}

func (dls *DomainListStatus) TakeSleepingCPU(vcpu int) error {
	num := int(dls.VcpuSleeping) 
	
	atomic.SwapInt32(&dls.VcpuSleeping, int32(num-vcpu))
	return nil
}

func (dls *DomainListStatus) GetVCPUStatus() VCPUStatus {
	status := VCPUStatus{
		Total: int(dls.VCPUTotal),
		Allocated: int(dls.VcpuAllocated),
		Sleeping: int(dls.VcpuSleeping),
		Idle: int(dls.VCPUTotal - dls.VcpuAllocated),
	}
	return status
}
