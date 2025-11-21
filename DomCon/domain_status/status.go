package domStatus

// cpu 는 각각 4개의 상태를 가짐
// 할당되어 활동중임
// 할당되어 있으나 유휴상태임(도메인이 꺼져있는 상태)
// 할당되어 있지 않음
// 총 cpu 수

// 우선순위 = 할당 되어 있지 않은 cpu 부터
// 그 이후 할당되어 있으나 유휴 상태인 cpu


type DomainListStatus struct {
	VCPUTotal int32 // 호스트 전체 cpu 수
	VcpuAllocated int32 // 할당 된 vcpu 수
	vcpuSleeping int32 // 유휴 상태인 vcpu 수
	// vcpuIdle = 할당되어 있지 않은 vcpu 수
	//VcpuIdle = VcpuTotal-VcpuAllocated
}
// 현재는 cpu 상태만이 필요한 상황
// 만약 추후에 메모리상에서 도메인 상태를 관리할 경우
// 새로운 패키지에서 호출할예정
type VCPUStatus struct{
	Total int `json:"total"`
	Allocated int `json:"allocated"`
	Sleeping int `json:"sleeping"`
	Idle int `json:"idle"`
}