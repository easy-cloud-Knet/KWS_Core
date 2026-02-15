package domCon

import (
	"sync"

	domStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domain_status"
	"libvirt.org/go/libvirt"
)

// libvirt 를 통해 도메인을 찾는 행위를 최소한 하기위해 관리하는 리스트(찾을때 걸리는 시간, 사용 후 도메인을 해제하는 과정 최소화)
// libvirt 내에서 domain은 *libvirt.domain으로 관리 됨
// DomainList 에서 uuid형태로 각각의 도메인을 관리






type DomListControl struct {
	DomainList map[string]*Domain
	domainListMutex sync.Mutex 
	DomainListStatus * domStatus.DomainListStatus
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


type DomainSeeker interface {
	ReturnDomain() (*Domain, error)
}
