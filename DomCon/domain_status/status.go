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
	deadCPU int
}
// 꺼져있는 도메인의 cpu 수

type libvirtStatus struct{
	liveCPU int
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
	EmitStatus() (VCPUStatus, error)
}

type VCPUStatus struct{
	Total int `json:"total"`
	Allocated int `json:"allocated"`
	Sleeping int `json:"sleeping"`
	Idle int `json:"idle"`
}

func (vs *DomainListStatus) EmitStatus() (VCPUStatus, error) {
	total := int(vs.VCPUTotal)
	allocated := int(vs.VcpuAllocated)
	sleeping := int(vs.VcpuSleeping)

	idle := total - allocated
	if idle < 0 {
		idle = 0
	}

	return VCPUStatus{
		Total:     total,
		Allocated: allocated,
		Sleeping:  sleeping,
		Idle:      idle,
	}, nil
}