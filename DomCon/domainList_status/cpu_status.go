package domListStatus

import (
	"runtime"
	"sync/atomic"
)

// 인터페이스 구현체

func (vs *VCPUStatus) EmitStatus(dls *DomainListStatus) {
	vs.Total = int(atomic.LoadInt64(&dls.VCPUTotal))
	vs.Allocated = int(atomic.LoadInt64(&dls.VcpuAllocated))
	vs.Sleeping = int(atomic.LoadInt64(&dls.VcpuSleeping))

	vs.Idle = max(vs.Total-vs.Allocated, 0)
}

func (dls *DomainListStatus) Update() {
	dls.UpdateCPUTotal()
}

func (dls *DomainListStatus) UpdateCPUTotal() {
	atomic.StoreInt64(&dls.VCPUTotal, int64(runtime.NumCPU()))
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
