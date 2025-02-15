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

	if err := DomainGenerator.DataParsor.Generate(); err!=nil{
		fmt.Println("do someting")
	}


	dom , err := i.CreateDomainWithXML(parsedXML)
	if err!= nil{
		resp.ResponseWriteErr(w,err, http.StatusInternalServerError)
		log.Printf("Error Creating VM with defined XML File  %v", err)
		return 
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

func (i *InstHandler) CreateDomainWithXML(config []byte) (*libvirt.Domain, error) {

	// DomainCreateXMLWithFiles를 호출하여 도메인을 생성합니다.
	domain, err := i.LibvirtInst.DomainDefineXML(string(config))
	if err != nil {
		return nil, virerr.ErrorGen(virerr.DomainGenerationError,err)
		// cpu나 ip 중복 등을 검사하는 코드를 삽입하고, 그에 맞는 에러 반환 필요
	} 
	//이전까지 생성 된 파일 삭제 해야됨.
  return domain ,err
}

