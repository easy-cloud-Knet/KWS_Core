package conn

import (
	domCon "github.com/easy-cloud-Knet/KWS_Core.git/api/conn/DomCon"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)


type DomainController struct{
	DomainList *domCon.DomListControl
	Operator Operator
}

type Operator interface{
	Operation() (*libvirt.Domain,error)
}


func DomainControllerInjection(domainList *domCon.DomListControl,Operator Operator)(*DomainController){
	return &DomainController{
		DomainList: domainList,
		Operator: Operator,
	}
}

func (DC *DomainController)DomainAddWithOperation(logger *zap.Logger,uuid string) error{
	
		domain,err:=DC.Operator.Operation()
		if err!=nil{
			return err
		}
		dom:=domCon.NewDomainInstance(domain)
		DC.DomainList.AddNewDomain(dom, uuid)
		
		err=dom.Domain.Create()
		if err != nil {
			return err
		}
		return nil
	
}
func (DC *DomainController)Operate() error{
	_,err:=DC.Operator.Operation()
	if err!=nil{
		return err
	}
	return nil
}


func (DC *DomainController)DomainDeleteWithOperation(logger *zap.Logger,uuid string) error{
	
	dom,err:=DC.Operator.Operation()
	if err!=nil{
		return err
	}
	DC.DomainList.DeleteDomain(dom, uuid)
	return nil

}