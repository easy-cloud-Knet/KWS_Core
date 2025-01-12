package conn

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	_ "libvirt.org/go/libvirt"
)


func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	var param DomainSeekinggByUUID
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
