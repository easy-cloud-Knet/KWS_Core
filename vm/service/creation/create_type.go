package creation

import (
	"github.com/easy-cloud-Knet/KWS_Core/vm/parsor"
	userconfig "github.com/easy-cloud-Knet/KWS_Core/vm/parsor/cloud-init"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

// 프로젝트내의 모든 도메인 생성은 VMCreator 라는 인터페이스를 통해 진행됨.
// 스냅샷 기반 생성, 파일기반 생성 등등 아래 규격을 따를 것.
// VMCreator의 조건
// 1. *libvirt.Domain 을 반환할 것(생성에 실패할 경우, 무조건 nil을 반환해야 함.)
// 2. CreateVM 이후 생성이 성공적이였음을 확인할 경우, DomCon에 vm 을 추가할 것
// 3. /var/lib/kws/{uuid} 내부에 가상 하드디스크, vm config에 대한 내용응 가지고 있어야 함.


type VMCreator interface {
	CreateVM() (*libvirt.Domain, error)
}



type LocalCreator struct{
	DomainConfiger *localConfigurer
	libvirtInst *libvirt.Connect
	logger *zap.Logger
}
// LocalCreator가 json에서 읽어온 데이터를 통해 새로운 vm을 만들 때 사용됨.
// CreateVM()의 구현체 

type NewDomainFromSnapshot struct{

}

type localConfigurer struct {
	VMDescription  *parsor.VM_Init_Info
	YamlParsorUser *userconfig.User_data_yaml
	YamlParsorMeta *userconfig.Meta_data_yaml
	DeviceDefiner  *parsor.VM_CREATE_XML
}
