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
	var param DomainSortingByStatus[PredefinedStructures]
	if err:= json.NewDecoder(r.Body).Decode(&param);err != nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	domainSorter:=&DomainSortingByStatus[PredefinedStructures]{
		Status: param.Status,
	}	
	domainSorter.TypeBase.Setter(param.TypeBase.DataType, i.LibvirtInst)
	domList, err:= domainSorter.returnDomain()
	if err!= nil{
		http.Error(w, "error while fetcing domain list", http.StatusBadRequest)
	}

	// Domlist,_:= i.ReturnDomainNameList(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	encoder := json.NewEncoder(w)
	encoder.Encode(&domList)

}



func (i *InstHandler)ReturnStatusUUID(w http.ResponseWriter, r * http.Request){
	var param DomainSortingByUUID[PredefinedStructures]
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	domainSorter:=&DomainSortingByUUID[PredefinedStructures]{
		UUID: param.UUID,
	}		
	Domain,err:=domainSorter.returnDomain()
	if err!=nil{
		fmt.Print("error occured while returning status ")
		http.Error(w, "there is no such VM with that UUID", 1)
	}
	encoder:=json.NewEncoder(w)
	encoder.Encode(&Domain)
}



func (DSU *DomainSortingByUUID[T])returnDomain()([]*libvirt.Domain,error){
	parsedUUID, err:= uuid.Parse(DSU.UUID)
	domainList:=make([]*libvirt.Domain,1)

	if err != nil {
        return domainList, fmt.Errorf("invalid uuid format: %w", err)
	}

	domain,err := DSU.TypeBase.LibvirtInst.LookupDomainByUUID(parsedUUID[:])
	if err != nil {
        return domainList, fmt.Errorf("invalid uuid format: %w", err)
	}
	domainList=append(domainList, domain)

	return domainList,nil
}

func (DSS *DomainSortingByStatus[T])returnDomain()([]*libvirt.Domain,error){
	
	doms, err := DSS.TypeBase.LibvirtInst.ListAllDomains(DSS.Status)
	if err != nil {
		fmt.Println("error while retrieving domain List with status")
	}
	Domains := make([]*libvirt.Domain,0,len(doms))
	
	for i:= range doms{
		Domains = append(Domains, &doms[i])
	}
	return Domains, nil
}
//******************** this to function allocate domain struct inside memory 
//all domain needed to be freed after certain operation done. *****************



func (DT *DataType[T]) ReturnBasicInfo(doms []*libvirt.Domain)([]*T, error){
	var domInfos []*T
	for i :=range doms{
		info, err:= doms[i].GetInfo()
		if err!=nil{
			log.Println(err)
		}
		DomInfo:= &DomainInfo{
			State :info.State,
			MaxMem :info.MaxMem,
			Memory : info.Memory,
			NrVirtCpu :info.NrVirtCpu,
			CpuTime :info.CpuTime,
		}
		//basic info can be added
		domInfos = append(domInfos, any(DomInfo).(*T))
	}
	return domInfos,nil
}
 

func (DSS *DomainSortingByStatus[T])ReturnDomainSpecification(doms []*libvirt.Domain)([]*T,error) {
	DomainInfo:=make([]*T,0,5)
	var err error
	switch(DSS.TypeBase.DataType){
		case PowerStaus:
			DomainInfo, err = DSS.TypeBase.PowerStatus(doms)
		case BasicInfo:
			DomainInfo, err =DSS.TypeBase.ReturnBasicInfo(doms)
		case GuestInfoUser:
		case GuestInfoOS:
		case GuestInfoFS:
		case GuestInfoDisk:
	}
	if err!= nil{
		fmt.Printf("error occured whild retreiving infos from Dom Spec, %v", err)
		return DomainInfo, err
	}
	return DomainInfo, nil
}

func (DSU *DomainSortingByUUID[T])ReturnDomainSpecification(doms []*libvirt.Domain)([]*T,error) {
	DomainInfo:=make([]*T,0,1)
	var err error
	switch(DSU.TypeBase.DataType){
		case PowerStaus:
			DomainInfo, err = DSU.TypeBase.PowerStatus(doms)
		case BasicInfo:
			DomainInfo, err =DSU.TypeBase.ReturnBasicInfo(doms)
		case GuestInfoUser:
		case GuestInfoOS:
		case GuestInfoFS:
		case GuestInfoDisk:
	}
	if err!= nil{
		fmt.Printf("error occured whild retreiving infos from Dom Spec, %v", err)
		return DomainInfo, err
	}
	return DomainInfo, nil
}


func (DT *DataType[T]) Setter(dataType DomainDataType,LibvirtInst *libvirt.Connect){
	DT.DataType=dataType 
	DT.LibvirtInst=LibvirtInst
}

func (DT *DataType[T]) PowerStatus(doms []*libvirt.Domain)([]*T, error){
	infos := make([]*T, 0,5)
	
	for i:= range doms{
		info, err:= doms[i].GetInfo()
		if err!=nil{
			log.Println(err)
		}
		infos=append(infos, T(info))
	}
	

	return infos,nil
}


 