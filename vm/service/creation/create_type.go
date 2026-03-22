package creation

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core/vm/parsor"
	userconfig "github.com/easy-cloud-Knet/KWS_Core/vm/parsor/cloud-init"
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/vm/types"
	"go.uber.org/zap"
)

// VMCreator is implemented by all VM creation strategies.
type VMCreator interface {
	CreateVM() (*domCon.Domain, error)
}

type LocalCreator struct {
	DomainConfiger Configurer
	libvirtInst    LibvirtConnect
	logger         *zap.Logger
}

// LocalCreator가 json에서 읽어온 데이터를 통해 새로운 vm을 만들 때 사용됨.
// CreateVM()의 구현체

type NewDomainFromSnapshot struct {
}

type localConfigurer struct {
	VMDescription  *vmtypes.VM_Init_Info
	YamlParsorUser *userconfig.User_data_yaml
	YamlParsorMeta *userconfig.Meta_data_yaml
	DeviceDefiner  *parsor.VM_CREATE_XML
}
