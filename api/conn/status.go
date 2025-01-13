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
	uuidBytes,_ := domain.Domain.GetUUID()
	uuidParsed,_:=uuid.FromBytes(uuidBytes)
	fmt.Println(uuidParsed.String())
	DP.DomainState = info
	DP.UUID= string(uuidParsed.String())

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
	
	data, _ := json.Marshal(domainSorter.DataHandle)
	if err != nil {
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (i *InstHandler)ReturnStatusUUID(w http.ResponseWriter, r * http.Request){
	var param ReturnDomainFromUUID
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	domainSorter:=&DomainDetail{
		DataHandle: make([]DataTypeHandler,0,1),
		DomainSeeker: &DomainSeekingByUUID{
			LibvirtInst: i.LibvirtInst,
			UUID: string(param.UUID),
			Domain: make([]*Domain,1),
		},
	}		
	outputStruct, _ := Generator(param.DataType)

	err:=domainSorter.DomainSeeker.SetDomain()
	if err!=nil{
		fmt.Print("error occured while returning status ")
		http.Error(w, "there is no such VM with that UUID", 1)
	}
	dom,_:=domainSorter.DomainSeeker.returnDomain()
	outputStruct.GetInfo(dom[0])
	fmt.Println(domainSorter.DataHandle)
	domainSorter.DataHandle=append(domainSorter.DataHandle, outputStruct)
	data, _ := json.Marshal(domainSorter.DataHandle)
	if err != nil {
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}



func (DSU *DomainSeekingByUUID)returnDomain()([]*Domain,error){
	return DSU.Domain,nil
}

func (DSS *DomainSeekingByStatus)returnDomain()([]*Domain,error){
	return DSS.DomList,nil
}


func (DSU *DomainSeekingByUUID)ReturnUUID()(uuid.UUID,error){
	UUID, err := uuid.Parse(DSU.UUID)
	if err!=nil{
		return uuid.UUID{}, err
	}
	return UUID,nil
}

func (DSU *DomainSeekingByUUID)SetDomain()(error){
	parsedUUID, err:= uuid.Parse(DSU.UUID)
	if err != nil {
        return  fmt.Errorf("invalid uuid format: %w", err)
	}
	domain,err := DSU.LibvirtInst.LookupDomainByUUID(parsedUUID[:])
	if err != nil {
        return  fmt.Errorf("invalid uuid format: %w", err)
	}
	Dom:=make([]*Domain,1)
	Dom[0]=&Domain{
		Domain:domain,
		DomainMutex: sync.Mutex{},
	}
	fmt.Println(Dom)
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
	fmt.Println(Domains)

	DSS.DomList=Domains
	return  nil
}
//******************** this to function allocate domain struct inside memory 
//all domain needed to be freed after certain operation done. *****************
