package domStatus

import (
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

// cpu 는 각각 4개의 상태를 가짐
// 할당되어 활동중임
// 할당되어 있으나 유휴상태임(도메인이 꺼져있는 상태)
// 할당되어 있지 않음
// 총 cpu 수

// 우선순위 = 할당 되어 있지 않은 cpu 부터
// 그 이후 할당되어 있으나 유휴 상태인 cpu

type DataDog interface {
	Retreive(*libvirt.Domain,*DomainListStatus ,zap.Logger) ( error  )  
}

type XMLStatus struct{
}
// 꺼져있는 도메인의 cpu 수

type libvirtStatus struct{
}
// 할당되어 활동중인 cpu 수

type DomainListStatus struct {
	VCPUTotal int32 // 호스트 전체 cpu 수
	VcpuAllocated int32 // 할당 된 vcpu 수
	VcpuSleeping int32 // 유휴 상태인 vcpu 수
	// vcpuIdle = 할당되어 있지 않은 vcpu 수
	//VcpuIdle = VcpuTotal-VcpuAllocated
}

//////////////////////////////

 
type StatusEmitter interface{
	EmitStatus() ( error)
}

type VCPUStatus struct{
	Total int `json:"total"`
	Allocated int `json:"allocated"`
	Sleeping int `json:"sleeping"`
	Idle int `json:"idle"`
}

func (vs *VCPUStatus) EmitStatus(dls *DomainListStatus) ( error) {
	vs.Total = int(dls.VCPUTotal)
	vs.Allocated = int(dls.VcpuAllocated)
	vs.Sleeping = int(dls.VcpuSleeping)
	
	vs.Idle = vs.Total - vs.Allocated
	if vs.Idle < 0 {
		vs.Idle = 0
	}

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