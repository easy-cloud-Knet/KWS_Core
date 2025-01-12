package conn

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"libvirt.org/go/libvirt"
)
 
func (DI *DomainInfo) GetInfo(domain *Domain) error{
		info, err:= domain.Domain.GetInfo()
		if err!=nil{
			log.Println(err)
		}
		DI.State=info.State
		DI.MaxMem=info.MaxMem
		DI.Memory= info.Memory
		DI.NrVirtCpu=info.NrVirtCpu
		DI.CpuTime=info.CpuTime
		//basic info can be added
	return nil
}
func (DP *DomainState) GetInfo(domain *Domain)error{
	info,_,err:= domain.Domain.GetState()
	//searching for coresponding second parameter, "Reason"
	if err!=nil{
		log.Println(err)
		return err
	}
	uuid,_ := domain.Domain.GetUUID()
	DP.DomainState = info
	DP.UUID= string(uuid)

	userInfo,err:= domain.Domain.GetGuestInfo(libvirt.DOMAIN_GUEST_INFO_USERS,0)
	if err!= nil{
		log.Println(err)
		return err
	}
	DP.Users=userInfo.Users
	
	return nil
}

func Generator(types DomainDataType)(DataTypeHandler,error){
	switch(types){
	case PowerStaus:
		return &DomainState{},nil
	case BasicInfo:
		return &DomainInfo{},nil
	case GuestInfoUser:
		return &DomainInfo{},nil
	case GuestInfoOS:
		return &DomainInfo{},nil
	case GuestInfoFS:
		return &DomainInfo{},nil
	case GuestInfoDisk:
}
	return &DomainInfo{},fmt.Errorf("error")
}

func (i * InstHandler) ReturnDomainByStatus(w http.ResponseWriter,r * http.Request){
	var param ReturnDomainFromStatus
	if err:= json.NewDecoder(r.Body).Decode(&param);err != nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	domainSorter:=&DomainDetail{
		DataHandle: []DataTypeHandler{},
		DomainSeeker: &DomainSeekingByStatus{
			LibvirtInst: i.LibvirtInst,
			Status: param.Status,
			DomList: make([]*Domain,5),
		},
	}
	// Domain Detail로 채우는 객체 생성
	err:= domainSorter.DomainSeeker.SetDomain()
	if err!= nil{
		http.Error(w, "error while fetcing domain list", http.StatusBadRequest)
	}
	//domain freeing need to be added
	list,_:=domainSorter.DomainSeeker.returnDomain()
	for i := range list{
		outputStruct, _ := Generator(param.DataType)
		outputInfo:= outputStruct
		outputInfo.GetInfo(list[i])
		// outputStruct.GetInfo(list[i])
		domainSorter.DataHandle=append(domainSorter.DataHandle, outputInfo)
	}
	encoder := json.NewEncoder(w)
	
	encoder.Encode(domainSorter.DataHandle)
}

func (i *InstHandler)ReturnStatusUUID(w http.ResponseWriter, r * http.Request){
	var param ReturnDomainFromUUID
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	domainSorter:=&DomainDetail{
		DataHandle: make([]DataTypeHandler,0,1),
		DomainSeeker: &DomainSeekinggByUUID{
			LibvirtInst: i.LibvirtInst,
			UUID: param.UUID,
			Domain: make([]*Domain,1),
		},
	}		
	outputStruct, _ := Generator(param.DataType)

	err:=domainSorter.DomainSeeker.SetDomain()
	if err!=nil{
		fmt.Print("error occured while returning status ")
		http.Error(w, "there is no such VM with that UUID", 1)
	}
	encoder:=json.NewEncoder(w)
	dom,_:=domainSorter.DomainSeeker.returnDomain()
	outputStruct.GetInfo(dom[0])
	domainSorter.DataHandle=append(domainSorter.DataHandle, outputStruct)
	encoder.Encode(domainSorter.DataHandle)
}

func (DSU *DomainSeekinggByUUID)returnDomain()([]*Domain,error){
	return DSU.Domain,nil
}

func (DSS *DomainSeekingByStatus)returnDomain()([]*Domain,error){
	return DSS.DomList,nil
}

func (DSU *DomainSeekinggByUUID)SetDomain()(error){
	parsedUUID, err:= uuid.Parse(DSU.UUID)
	if err != nil {
        return  fmt.Errorf("invalid uuid format: %w", err)
	}

	domain,err := DSU.LibvirtInst.LookupDomainByUUID(parsedUUID[:])
	if err != nil {
        return  fmt.Errorf("invalid uuid format: %w", err)
	}
	Dom:=make([]*Domain,0,1)
	Dom[0]=&Domain{
		Domain:domain,
		DomainMutex: sync.Mutex{},
	}
	DSU.Domain=Dom
	return nil
}

func (DSS *DomainSeekingByStatus)SetDomain()(error){
	doms, err := DSS.LibvirtInst.ListAllDomains(DSS.Status)
	if err != nil {
		fmt.Println("error while retrieving domain List with status")
        return  fmt.Errorf("invalid uuid format: %w", err)
	}
	Domains := make([]*Domain,0,len(doms))
	
	for i:= range doms{
		Domains = append(Domains, &Domain{Domain:&doms[i], DomainMutex: sync.Mutex{}})
	}
	DSS.DomList=Domains
	return  nil
}
//******************** this to function allocate domain struct inside memory 
//all domain needed to be freed after certain operation done. *****************



// func (DT *DomainDetail) ReturnBasicInfo(doms []*libvirt.Domain)([]*T, error){
// 	var domInfos []*T
// 	for i :=range doms{
// 		info, err:= doms[i].GetInfo()
// 		if err!=nil{
// 			log.Println(err)
// 		}
// 		DomInfo:= &DomainInfo{
// 			State :info.State,
// 			MaxMem :info.MaxMem,
// 			Memory : info.Memory,
// 			NrVirtCpu :info.NrVirtCpu,
// 			CpuTime :info.CpuTime,
// 		}
// 		//basic info can be added
// 		domInfos = append(domInfos, any(DomInfo).(*T))
// 	}
// 	return domInfos,nil
// }
 

// func (DSS *DomainDetail)ReturnDomainSpecification(doms []*libvirt.Domain)([]*T,error) {
// 	DomainInfo:=make([]*T,0,5)
// 	var err error

// 	if err!= nil{
// 		fmt.Printf("error occured whild retreiving infos from Dom Spec, %v", err)
// 		return DomainInfo, err
// 	}
// 	return DomainInfo, nil
// }

// func (DSU *DomainDetail)ReturnDomainSpecification(doms []*libvirt.Domain)([]*T,error) {
// 	DomainInfo:=make([]*T,0,1)
// 	var err error
// 	switch(DSU.TypeBase.DataType){
// 		case PowerStaus:
// 			DomainInfo, err = DSU.TypeBase.PowerStatus(doms)
// 		case BasicInfo:
// 			DomainInfo, err =DSU.TypeBase.ReturnBasicInfo(doms)
// 		case GuestInfoUser:
// 		case GuestInfoOS:
// 		case GuestInfoFS:
// 		case GuestInfoDisk:
// 	}
// 	if err!= nil{
// 		fmt.Printf("error occured whild retreiving infos from Dom Spec, %v", err)
// 		return DomainInfo, err
// 	}
// 	return DomainInfo, nil
// }


// func (DT *DomainDetail) Setter(dataType DomainDataType,LibvirtInst *libvirt.Connect){
	
// }

// func (DT *DomainDetail) PowerStatus(doms []*libvirt.Domain)([]*T, error){
// 	infos := make([]*T, 0,5)
	
// 	for i:= range doms{
// 		info, err:= doms[i].GetInfo()
// 		if err!=nil{
// 			log.Println(err)
// 		}
// 		infos=append(infos, T(info))
// 	}
	
// 	return infos,nil
// }


 