package service

import (
	"fmt"
	"log"
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
		return
	}

	DomainParsor:= parsor.ParsorFactoryFromRequest(param)

	DomainGenerator := &conn.DomainGenerator{
		DataParsor: DomainParsor,
	}

	dom,err := DomainGenerator.DataParsor.Generate(i.LibvirtInst); 
	if err!=nil{
		resp.ResponseWriteErr(w,err, http.StatusInternalServerError)

	}
	fmt.Println("domain specification", dom)
	err = dom.Create()
	if err!= nil{
		resp.ResponseWriteErr(w,err, http.StatusInternalServerError)
		log.Printf("Error starting VM, check for Host's Ram Capacity  %v", err)
		return 
	}
	newDomain := conn.DomGen(dom)

	i.DomainControl.AddNewDomain(newDomain,param.UUID)

	// domainInfo,err:= dom.GetInfo()
	// if err!=nil{
	// 	appendingErorr:=virerr.ErrorJoin(virerr.DomainStatusError, errors.New("retreving Domain Status Error in creating VM workload"))
	// 	resp.ResponseWriteErr(w,appendingErorr, http.StatusInternalServerError)
	// 	return 
	// }
	resp.ResponseWriteOK(w,nil)
}


