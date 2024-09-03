package conn

import (
	"fmt"
	"net/http"

	"libvirt.org/go/libvirt"
)



type DomainList struct{
	RequestType string `json:"requestType"` 
	// libvirt.ConnectListAllDomainsFlags
}

func (i * InstHandler) ReturnStatus(w http.ResponseWriter,r * http.Request){
	fmt.Println("getStatus request income")

	i.ReturnDomainNameList(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
}
