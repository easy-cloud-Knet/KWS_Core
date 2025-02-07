package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
	"libvirt.org/go/libvirt"
)

func (i *InstHandler) CreateVM(w http.ResponseWriter, r *http.Request) {
	param:=&parsor.VM_Init_Info{}
	resp:=ResponseGen("CreateVm")

	if err:=HttpDecoder(w,r,param); err!=nil{
		http.Error(w, "error decoding parameters", http.StatusBadRequest)
		return
	}
	//생성 방법에 따라 다른 Generator 선언 필요
	DomainFromLocal := &conn.DomainGeneratorLocal{
		DomainStatusManager: &conn.DomainStatusManager{
			UUID: param.UUID,
		},
		OS: param.OS,
	}
	
	err:= DomainFromLocal.CreateFolder()
	if err!=nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error creating ISO File",err), http.StatusInternalServerError)
		log.Printf("Error creating directory  %v", err)
		return 
	}

	DomainFromLocal.DataParsor.YamlParsor=&parsor.User_data_yaml{}
	err=DomainFromLocal.CloudInitConf(param)
	if err!=nil{
		log.Printf("Error writing user File  %v", err)
		resp.ResponseWriteErr(w,fmt.Errorf("%w error creating ISO File",err), http.StatusInternalServerError)
		return 
	}
	DomainFromLocal.DataParsor.YamlParsor = &parsor.Meta_data_yaml{}

	err=DomainFromLocal.CloudInitConf(param)
	if err!=nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error creating ISO File",err), http.StatusInternalServerError)
		log.Printf("Error writing Meta data  %v", err)
		return 
	}
	parsedXML := parsor.XML_Parsor(param)

 
	if err:= DomainFromLocal.CreateDiskImage();err!=nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error creating ISO File",err), http.StatusInternalServerError)
		log.Printf("Error writing XML file  %v", err)
		return 
	}
	if err:= DomainFromLocal.CreateISOFile();err!=nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error creating ISO File",err), http.StatusInternalServerError)
		log.Printf("Error Creating ISO File  %v", err)
		return 
	}

	dom , err := i.CreateDomainWithXML(parsedXML)
	if err!= nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error booting vm",err), http.StatusInternalServerError)
		log.Printf("Error Creating VM with defined XML File  %v", err)
		return 
	}
	err = dom.Create()
	if err!= nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error booting vm",err), http.StatusInternalServerError)
		log.Printf("Error starting VM, check for Host's Ram Capacity  %v", err)
		return 
	}

	domainInfo,err:= dom.GetInfo()
	if err!=nil{
		
	}
	resp.ResponseWriteOK(w,domainInfo)
}

func (i *InstHandler) CreateDomainWithXML(config []byte) (*libvirt.Domain, error) {

	// DomainCreateXMLWithFiles를 호출하여 도메인을 생성합니다.
	domain, err := i.LibvirtInst.DomainDefineXML(string(config))
	if err != nil {
		return nil, fmt.Errorf("error generating XML File %v",err)
	} 
	//이전까지 생성 된 파일 삭제 해야됨.
  return domain ,err
}

