package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
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
		Domain: conn.Domain{},
		DataParsor: DomainParsor,
	}

	dom,err := DomainGenerator.DataParsor.Generate(i.LibvirtInst); 
	if err!=nil{
		fmt.Println("do someting")
	}

	err = dom.Create()
	if err!= nil{
		resp.ResponseWriteErr(w,err, http.StatusInternalServerError)
		log.Printf("Error starting VM, check for Host's Ram Capacity  %v", err)
		return 
	}
	

	domainInfo,err:= dom.GetInfo()
	if err!=nil{
		appendingErorr:=virerr.ErrorJoin(virerr.DomainStatusError, errors.New("retreving Domain Status Error in creating VM workload"))
		resp.ResponseWriteErr(w,appendingErorr, http.StatusInternalServerError)
		return 
	}
	resp.ResponseWriteOK(w,domainInfo)
}


