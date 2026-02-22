package domListStatus

// cpu 는 각각 4개의 상태를 가짐
// 할당되어 활동중임
// 할당되어 있으나 유휴상태임(도메인이 꺼져있는 상태)
// 할당되어 있지 않음
// 총 cpu 수

type DomainListStatus struct {
	VCPUTotal     int64 // 호스트 전체 cpu 수
	VcpuAllocated int64 // 할당 된 vcpu 수
	VcpuSleeping  int64 // 유휴 상태인 vcpu 수
	// vcpuIdle = 할당되어 있지 않은 vcpu 수
	//VcpuIdle = VcpuTotal-VcpuAllocated
}

//////////////////////////////

type VCPUStatus struct {
	Total     int `json:"total"`
	Allocated int `json:"allocated"`
	Sleeping  int `json:"sleeping"`
	Idle      int `json:"idle"`
}
