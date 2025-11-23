package domStatus

import (
	"fmt"

	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)



func NewDataDog(state libvirt.ConnectListAllDomainsFlags) DataDog {
	switch state {
		case libvirt.CONNECT_LIST_DOMAINS_ACTIVE:
			fmt.Println("returning active")
			return &libvirtStatus{}
		case libvirt.CONNECT_LIST_DOMAINS_INACTIVE:
			fmt.Println("returning inactive")
			return &XMLStatus{}
		default:
			return nil
	}
}


func (ds *XMLStatus) Retreive(dom *libvirt.Domain,DLS *DomainListStatus, logger zap.Logger) ( error) {
	domcnf, err := XMLUnparse(dom)
	if err != nil {
		logger.Error("failed to unparse domain XML", zap.Error(err))
		return nil
	}
		DLS.AddAllocatedCPU(int(domcnf.VCPU.Value))
		DLS.AddSleepingCPU(int(domcnf.VCPU.Value))
		return nil

}	
	
func (ls *libvirtStatus) Retreive(dom *libvirt.Domain, DLS *DomainListStatus, logger zap.Logger) (error) {
	cpuCount, err := dom.GetMaxVcpus()
	if err != nil {
		logger.Error("failed to get live vcpu count", zap.Error(err))
		return nil
	}

	fmt.Printf("%+v", cpuCount)
	DLS.AddAllocatedCPU(int(cpuCount))
	return nil
	
}