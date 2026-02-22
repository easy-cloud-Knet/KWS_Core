package domListStatus

import (
	"runtime"
	"sync/atomic"
)

type VCPUStatus struct {
	Total     int `json:"total"`
	Allocated int `json:"allocated"`
	Sleeping  int `json:"sleeping"`
	Idle      int `json:"idle"`
}

// 인터페이스 구현체

func (vs *VCPUStatus) EmitStatus(dls *DomainListStatus) error {
	vs.Total = int(dls.VCPUTotal)
	vs.Allocated = int(dls.VcpuAllocated)
	vs.Sleeping = int(dls.VcpuSleeping)

	vs.Idle = vs.Total - vs.Allocated
	if vs.Idle < 0 {
		vs.Idle = 0
	}

	return nil
}

func (dls *DomainListStatus) Update() {
	dls.UpdateCPUTotal()
}

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
