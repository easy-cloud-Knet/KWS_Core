package creation

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core.git/api/conn/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
	userconfig "github.com/easy-cloud-Knet/KWS_Core.git/api/parsor/cloud-init"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type DomainCreator struct{
	DomainConfiger *NewDomainFromBase
	libvirtInst *libvirt.Connect
	Domain *domCon.Domain
	logger *zap.Logger
}
// 아직은 확장성만 고려, snapshot 이나 firecracker 등 다른 방식을 보고
// DomainConfiger 나 다른 필드를 추상화 할 수 있을 지 생각

type NewDomainFromSnapshot struct{

}

type NewDomainFromBase struct {
	VMDescription  *parsor.VM_Init_Info
	YamlParsorUser *userconfig.User_data_yaml
	YamlParsorMeta *userconfig.Meta_data_yaml
	DeviceDefiner  *parsor.VM_CREATE_XML
}
// service에서 도메인 생성요청이 들어오면 사용되는 구조체 
// 현재는 로컬 파일을(cloud-init) 을 생성해서 만드는 거 밖에 없음,
//firecracker를 사용한 vm생성이나, snapshot 기반으로 만들거나 할때도 사용할 수 있게 하면 좋을 듯. 
