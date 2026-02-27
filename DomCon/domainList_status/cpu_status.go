package domListStatus

import (
	"runtime"
	"sync/atomic"
)

// 인터페이스 구현체

func (vs *VCPUStatus) EmitStatus(dls *DomainListStatus) {
	vs.Total = int(dls.VCPUTotal)
	vs.Allocated = int(dls.VcpuAllocated)
	vs.Sleeping = int(dls.VcpuSleeping)

	vs.Idle = vs.Total - vs.Allocated
	if vs.Idle < 0 {
		vs.Idle = 0
	}
}

func (dls *DomainListStatus) Update() {
	dls.UpdateCPUTotal()
}

func (dls *DomainListStatus) UpdateCPUTotal() {
	totalCPU := runtime.NumCPU()
	dls.VCPUTotal = int64(totalCPU)
}

func (dls *DomainListStatus) AddAllocatedCPU(vcpu int) {
	atomic.AddInt64(&dls.VcpuAllocated, int64(vcpu))
}

func (dls *DomainListStatus) AddSleepingCPU(vcpu int) {
	atomic.AddInt64(&dls.VcpuSleeping, int64(vcpu))
}

func (dls *DomainListStatus) TakeAllocatedCPU(vcpu int) {
	atomic.AddInt64(&dls.VcpuAllocated, -int64(vcpu))
}

func (dls *DomainListStatus) TakeSleepingCPU(vcpu int) {
	atomic.AddInt64(&dls.VcpuSleeping, -int64(vcpu))
}
