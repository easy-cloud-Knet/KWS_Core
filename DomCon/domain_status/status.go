package domainStatus

import (
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type SourceType string

const (
	CPU    SourceType = "cpu"
	Memory SourceType = "memory"
)

// DataDog interface 는 도메인의 상태를 가져오는 인터페이스로, XMLStatus 와 LibvirtStatus 가 구현한다.
// 도메인의 상태에 따라 상태를 가져오는 방식이 다르기 때문에, 인터페이스로 추상화하여 구현체에서 각각의 방식으로 상태를 가져오도록 한다.
// source.go 에 각각 의 상태를 가져오는 함수를 구현한다. (RetrieveCPU, RetrieveMemory, RetrieveHDD 등)
// RetrieveStatus 함수는 도메인의 상태를 가져오는 함수로, enum 상태를 읽고 요구하는 데이터를 가져옴.
// 각각 다르게 구현하면 좋겠지만, libvirt 나 xml 메모리 크기를 고려해서 한번에 가져오는 식으로 구현 했음.

type DataDog interface {
	// 아작 반환 타입이 정해져 있지 않기 때문에, interface{} 로 반환 타입을 설정. 필요에 따라 구체적인 타입으로 변환하여 사용.
	RetrieveStatus(*libvirt.Domain, []SourceType, *zap.Logger) (interface{}, error)
}

type XMLStatus struct{}

type LibvirtStatus struct{}
