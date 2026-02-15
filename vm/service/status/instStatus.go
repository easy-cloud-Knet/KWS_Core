package status

import (
	"fmt"

	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"libvirt.org/go/libvirt"
)

func (AII *AllInstInfo) GetAllinstInfo(LibvirtInst *libvirt.Connect) error {

	domains, err := LibvirtInst.ListAllDomains(0) //alldomain
	if err != nil {
		return virerr.ErrorGen(virerr.HostStatusError, fmt.Errorf("failed to list all domains: %w", err))
	}

	var totalMaxMem uint64
	var totalvCPU uint
	for _, dom := range domains {
		data, err := dom.GetInfo()
		if err != nil {
			dom.Free()
			continue
		}
		totalMaxMem += data.MaxMem
		totalvCPU += data.NrVirtCpu
		dom.Free()
	}

	AII.Totalmaxmem = totalMaxMem
	AII.TotalVCpu = totalvCPU
	return nil
}

func InstDataTypeRouter(types InstDataType) (InstDataTypeHandler, error) {
	switch types {
	case Vcpu_MaxMem:
		return &AllInstInfo{}, nil
	}

	return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("unsupported type"))
}

func InstDetailFactory(handler InstDataTypeHandler, LibvirtInst *libvirt.Connect) (*InstDetail, error) {
	if err := handler.GetAllinstInfo(LibvirtInst); err != nil {
		return nil, err
	}
	return &InstDetail{
		AllInstDataHandle: handler,
	}, nil
}
