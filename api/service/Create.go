package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
	"libvirt.org/go/libvirt"
)

func (i *InstHandler) CreateVM(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var param parsor.VM_Init_Info

	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		CommonErrorHelper(w,err,http.StatusInternalServerError, "error while Decoding reqeust")
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
		CommonErrorHelper(w,err, http.StatusInternalServerError,"Error creating directory")
		log.Printf("Error creating directory  %v", err)
		return 
	}

	DomainFromLocal.DataParsor.YamlParsor=&parsor.User_data_yaml{}
	err=DomainFromLocal.CloudInitConf(&param)
	if err!=nil{
		log.Printf("Error writing user File  %v", err)
		CommonErrorHelper(w,err, http.StatusInternalServerError,"Error writing XML file ")
		return 
	}
	DomainFromLocal.DataParsor.YamlParsor = &parsor.Meta_data_yaml{}

	err=DomainFromLocal.CloudInitConf(&param)
	if err!=nil{
		CommonErrorHelper(w,err, http.StatusInternalServerError,"Error writing MetaData ")
		log.Printf("Error writing Meta data  %v", err)
		return 
	}
	parsedXML := parsor.XML_Parsor(&param)

 
	if err:= DomainFromLocal.CreateDiskImage();err!=nil{
		CommonErrorHelper(w,err, http.StatusInternalServerError,"Error Creating DiskImage")
		log.Printf("Error writing XML file  %v", err)
		return 
	}
	if err:= DomainFromLocal.CreateISOFile();err!=nil{
		CommonErrorHelper(w,err, http.StatusInternalServerError,"Error Creating ISO File ")

		log.Printf("Error Creating ISO File  %v", err)
		return 
	}

	dom , err := i.CreateDomainWithXML(parsedXML)
	if err!= nil{
		CommonErrorHelper(w,err, http.StatusInternalServerError,"Error Creating VM with defined XML File ")
		log.Printf("Error Creating VM with defined XML File  %v", err)
		return 
	}
	err = dom.Create()
	if err!= nil{
		CommonErrorHelper(w,err, http.StatusInternalServerError,"Error starting VM, check for Host Capacity")
		log.Printf("Error starting VM, check for Host's Ram Capacity  %v", err)
		return 
	}

	domainInfo,_:= dom.GetInfo()
	data, err:= json.Marshal(domainInfo)
	if err!=nil{
		CommonErrorHelper(w,err, http.StatusInternalServerError,"Error Mashaling Data")
		return
	}
	w.WriteHeader(http.StatusOK)

	w.Write(data)

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


func DomainCleanUp(){

}
