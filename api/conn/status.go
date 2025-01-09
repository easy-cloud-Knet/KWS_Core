package conn

import (
	"encoding/json"
	"fmt"
	"net/http"

	"libvirt.org/go/libvirt"
)



func (i * InstHandler) ReturnStatus(w http.ResponseWriter,r * http.Request){
	fmt.Println("getStatus request income")

	Domlist,_:= i.ReturnDomainNameList(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	encoder := json.NewEncoder(w)
	encoder.Encode(&Domlist)

}

func (i *InstHandler)ReturnStatusUUID(w http.ResponseWriter, r * http.Request){
	var param SpecifiyUUID
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}
	// i.InstHandler
}