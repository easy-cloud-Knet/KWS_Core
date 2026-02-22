package domCon

import (
	"sync"

	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domainList_status"
	"libvirt.org/go/libvirt"
)

// DomainListControl 은 libvirt 도메인 리스트를 관리하는 구조체
// Libvirt.Domain 객체를 추상화 함.
// libvirt-go 호출 최소화, vcpu 상태 관리, 도메인 추가/삭제 관리 등을 담당

// 실제 생성/삭제는 service 레이어에서 담당, 도메인 리스트 관리에 집중

type DomListControl struct {
	DomainList       map[string]*Domain
	domainListMutex  sync.Mutex
	DomainListStatus *domStatus.DomainListStatus
}

// 각 도메인을 관리하는 인메모리 구조체
// mutex 를 통해 동시성 제어
// 메모리 누수방지 + libvirt 접근 최소화 위해 libvirt.Domain 포인터를 보유

type Domain struct {
	Domain      *libvirt.Domain
	domainMutex sync.Mutex
}

// 각 도메인의 상태변경시에 사용하는 구조체
// mutex 를 통해 동시성 제어

type DomainSeekingByUUID struct {
	LibvirtInst *libvirt.Connect
	UUID        string
}

// 도메인 탐색 인터페이스
// 인메모리 도메인 리스트에 없을 경우 libvirt 에서 도메인을 찾아 반환하는 역할
