package status

import (
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type DataDog interface {
	RetrieveStatus(*libvirt.Domain, map[vmtypes.SourceType]int, *zap.Logger) (map[vmtypes.SourceType]int, error)
}

func New(isActive bool) DataDog {
	if isActive {
		return &LibvirtStatus{}
	}
	return &XMLStatus{}
}
