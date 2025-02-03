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
	var param parsor.VM_Init_Info

	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}
	//생성 방법에 따라 다른 Generator 선언 필요
	DomainFromLocal := &conn.DomainGeneratorLocal{
		DomainStatusManager: &conn.DomainStatusManager{
			UUID: param.UUID,
		},
		OS: param.OS,
	}

	err := DomainFromLocal.CreateFolder()
	if err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		log.Printf("Error creating directory  %v", err)
	}

	DomainFromLocal.DataParsor.YamlParsor = &parsor.User_data_yaml{}
	err = DomainFromLocal.CloudInitConf(&param)
	if err != nil {
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}
	DomainFromLocal.DataParsor.YamlParsor = &parsor.Meta_data_yaml{}

	err = DomainFromLocal.CloudInitConf(&param)
	if err != nil {
		http.Error(w, "Failed to marshal meta data", http.StatusInternalServerError)
		return
	}
	parsedXML := parsor.XML_Parsor(&param)

	if err := DomainFromLocal.CreateDiskImage(); err != nil {
		http.Error(w, "Failed to create disk image", http.StatusInternalServerError)

	}
	if err := DomainFromLocal.CreateISOFile(); err != nil {
		http.Error(w, "Failed to create disk image", http.StatusInternalServerError)

	}

	dom, err := i.CreateDomainWithXML(parsedXML)
	if err != nil {
		http.Error(w, "faild creating vm", http.StatusConflict)
	}
	err = dom.Create()
	if err != nil {
		http.Error(w, "faild starting vm", http.StatusConflict)
	}

	domainInfo, _ := dom.GetInfo()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "VM with UUID %s created successfully., %v", param.UUID, domainInfo)
}

func (i *InstHandler) CreateDomainWithXML(config []byte) (*libvirt.Domain, error) {

	// DomainCreateXMLWithFiles를 호출하여 도메인을 생성합니다.
	domain, err := i.LibvirtInst.DomainDefineXML(string(config))
	if err != nil {
		log.Fatal(err)
	}
	return domain, err

}
