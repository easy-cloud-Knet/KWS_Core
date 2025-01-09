package conn

import (
	"encoding/json"
	"net/http"
	"fmt"

	_ "libvirt.org/go/libvirt"
	"github.com/google/uuid"
)


func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	var param SpecifiyUUID
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		http.Error(w, "error decoding body", 1)
	}

	parsedUUID, err := uuid.Parse(param.UUID)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	domain,err := i.LibvirtInst.LookupDomainByUUID(parsedUUID[:])

	if err!=nil{
		fmt.Println(err)	
		//error handler needed
		return
	}
	err=domain.Destroy()
	if err!=nil{
		fmt.Println(err)	
		//error handler needed
		return
	}
	domainInfo,err:= domain.GetInfo()
	if err!=nil{
		fmt.Println(err)	
		//error handler needed
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(&domainInfo)
	
}
