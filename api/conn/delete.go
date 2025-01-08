package conn

import (
	"encoding/json"
	"net/http"

	"libvirt.org/go/libvirt"
)


func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	var param SpecifiyUUID
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		http.Error(w, "error decoding body", 1)
	}
	domain,err := i.LibvirtInst.LookupDomainByUUID([]byte(param.UUID))
	if err!=nil{
		//error handler needed
		return
	}
	err=domain.Destroy()
	if err!=nil{
		//error handler needed
		return
	}
	Domlist,_:= i.ReturnDomainNameList(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	encoder := json.NewEncoder(w)
	encoder.Encode(&Domlist)
	// i.LibvirtInst.LookupDomainByUUID()
	
}