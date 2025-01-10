package conn

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"libvirt.org/go/libvirt"
)
 





func (i * InstHandler) ReturnDomainByStatus(w http.ResponseWriter,r * http.Request){
	fmt.Println("getStatus request income")
	var param DomainSortingByStatus
	if err:= json.NewDecoder(r.Body).Decode(&param);err != nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	domainSorter:=&DomainSortingByStatus{
		Status: param.Status,
	}	
	domainSorter.TypeBase.ParamSetter(param.TypeBase, i.LibvirtInst)
	domList, err:= domainSorter.returnStatus()
	if err!= nil{
		http.Error(w, "error while fetcing domain list", http.StatusBadRequest)
	}

	// Domlist,_:= i.ReturnDomainNameList(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	encoder := json.NewEncoder(w)
	encoder.Encode(&domList)

}
func (DT *DataType) ParamSetter(dataType DomainDataType,LibvirtInst *libvirt.Connect){
	DT.DataType=dataType 
	DT.LibvirtInst=LibvirtInst
}
func (i *InstHandler)ReturnStatusUUID(w http.ResponseWriter, r * http.Request){
	var param DomainSortingByUUID
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	domainSorter:=&DomainSortingByUUID{
		UUID: param.UUID,
	}		
	Domain,err:=domainSorter.returnStatus()
	if err!=nil{
		fmt.Print("error occured while returning status ")
		http.Error(w, "there is no such VM with that UUID", 1)
	}
	encoder:=json.NewEncoder(w)
	encoder.Encode(&Domain)
}



func (DSU *DomainSortingByUUID)returnStatus()([]*libvirt.Domain,error){
	parsedUUID, err:= uuid.Parse(DSU.UUID)
	if err != nil {
        return make([]*libvirt.Domain, 1), fmt.Errorf("invalid uuid format: %w", err)
	}
	domain,err := DSU.LibvirtInst.LookupDomainByUUID(parsedUUID[:])
	if err != nil {
        return make([]*libvirt.Domain, 1), fmt.Errorf("invalid uuid format: %w", err)
	}
	domainList:=make([]*libvirt.Domain,1)
	domainList=append(domainList, domain)

	return domainList,nil
}

func (DSS *DomainSortingByStatus)returnStatus()([]*libvirt.Domain,error){
	
	doms, err := DSS.LibvirtInst.ListAllDomains(DSS.Status)
	if err != nil {
		panic(err)
	}
	Domains := make([]*libvirt.Domain,0,len(doms))
	
	for i:= range doms{
		Domains = append(Domains, &doms[i])
	}
	return Domains, nil
}



func (i *InstHandler)ReturnDomainNameList(flag libvirt.ConnectListAllDomainsFlags)([]*DomainInfo,error) {
	var Domains []*DomainInfo

	doms, err := i.LibvirtInst.ListAllDomains(flag)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {
		info, err := dom.GetInfo()
		if  err!=nil {
				log.Println(err)
		}
		uuid, err := dom.GetUUIDString()
		if err!= nil{
			log.Panicln(err)
		}
		DomInfo:= &DomainInfo{
				State :info.State,
				MaxMem :info.MaxMem,
				Memory : info.Memory,
				NrVirtCpu :info.NrVirtCpu,
				CpuTime :info.CpuTime,
				UUID : uuid,
		}
		
		Domains=append(Domains,DomInfo)
		dom.Free()
	}
	Use(Domains)
	return Domains,nil
}

// func (DSU *DomainSortingByUUID)InstGetter(CurrentInst *libvirt.Connect){
// 	DSU.LibvirtInst=CurrentInst	
// }
// func (DSS *DomainSortingByStatus)InstGetter(CurrentInst *libvirt.Connect){
// 	DSS.LibvirtInst=CurrentInst	
// }
// func (DSU *DomainSortingByUUID)InstSetter()(*libvirt.Connect,error){
// 	return DSU.LibvirtInst,nil
// }
// func (DSS *DomainSortingByStatus)InstSetter()(*libvirt.Connect,error){
// 	return DSS.LibvirtInst,nil
// }