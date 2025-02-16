package service

import (
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
	"libvirt.org/go/libvirt"
)


func (i *InstHandler) CreateVMLocal(w http.ResponseWriter, r *http.Request) {

	resp:=ResponseGen[libvirt.DomainInfo]("CreateVm")
	param:=&parsor.VM_Init_Info{}
	
	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,err, http.StatusBadRequest)
		i.Logger.Errorf("error occured while decoding user's parameter of requested creation")
		return
	}
	i.Logger.Infoln("Handling Create VM of uuid of %s", param.UUID)

	DomainParsor:= parsor.ParsorFactoryFromRequest(param, i.Logger)

	DomainGenerator := &conn.DomainGenerator{
		DataParsor: DomainParsor,
	}

	dom,err := DomainGenerator.DataParsor.Generate(i.LibvirtInst, i.Logger); 
	if err!=nil{
		resp.ResponseWriteErr(w,err, http.StatusInternalServerError)
		return
	}

	err = dom.Create()
	if err!= nil{
		resp.ResponseWriteErr(w,err, http.StatusInternalServerError)
		return 
	}
	newDomain := conn.NewDomainInstance(dom)

	i.DomainControl.AddNewDomain(newDomain,param.UUID)

	resp.ResponseWriteOK(w,nil)
}


