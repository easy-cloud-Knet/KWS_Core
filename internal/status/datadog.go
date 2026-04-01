package status

import (
	"go.uber.org/zap"
)

type StatusRetriever interface {
	RetrieveStatus(map[SourceType]int, *zap.Logger) (map[SourceType]int, error)
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
