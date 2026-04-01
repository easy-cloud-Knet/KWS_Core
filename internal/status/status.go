package status

import (
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	"go.uber.org/zap"
)

type StatusRetriever interface {
	RetrieveStatus(map[vmtypes.SourceType]int, *zap.Logger) (map[vmtypes.SourceType]int, error)
}

type statusDomain interface {
	libvirtDom
	xmlDom
}

func New(dom statusDomain, isActive bool) StatusRetriever {
	if isActive {
		return &LibvirtStatus{dom: dom}
	}
	return &XMLStatus{dom: dom}
}
