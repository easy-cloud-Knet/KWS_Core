package domStatus

import (
	"runtime"
	"sync/atomic"
)

func (dls *DomainListStatus) UpdateCPUTotal() {
	totalCPU := runtime.NumCPU()
	dls.VCPUTotal = int64(totalCPU)
}

func (dls *DomainListStatus) AddAllocatedCPU(vcpu int) error {
	atomic.AddInt64(&dls.VcpuAllocated, int64(vcpu))
	return nil
}

func (dls *DomainListStatus) AddSleepingCPU(vcpu int) error {
	atomic.AddInt64(&dls.VcpuSleeping, int64(vcpu))
	return nil
}
func (dls *DomainListStatus) TakeAllocatedCPU(vcpu int) error {
	
	atomic.AddInt64(&dls.VcpuAllocated, -int64(vcpu))
	return nil
}

func (dls *DomainListStatus) TakeSleepingCPU(vcpu int) error {
	
	atomic.AddInt64(&dls.VcpuSleeping, -int64(vcpu))
	return nil
}


