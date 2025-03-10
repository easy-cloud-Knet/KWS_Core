package termination

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core.git/api/conn/DomCon"
	"libvirt.org/go/libvirt"
)


func DomainTerminatorFactory(Domain *domCon.Domain) (*DomainTerminator, error) {
	return &DomainTerminator{
		domain: Domain,
	}, nil
}

func (DD *DomainTerminator) Operation()(*libvirt.Domain,error){
	dom:= DD.domain

	isRunning, _ := dom.Domain.IsActive()
	if !isRunning {
		return  nil,fmt.Errorf("requested Domain to shutdown is already Dead ")
	}

	if err := dom.Domain.Destroy(); err != nil {
		fmt.Println("error occured while deleting Domain")
		return nil,fmt.Errorf("internal Error in Libvirt occured while shutting down domain")
	}

	return dom.Domain,nil
}


