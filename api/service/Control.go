package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)


func (i * InstHandler) ReturnDomainByStatus(w http.ResponseWriter,r * http.Request){
	var param ReturnDomainFromStatus
	if err:= json.NewDecoder(r.Body).Decode(&param);err != nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	DataHandle,_:= conn.DataTypeRouter(param.DataType)
	DomainSeeker:= &conn.DomainSeekingByStatus{
		LibvirtInst: i.LibvirtInst,
		Status: param.Status,
		DomList: make([]*conn.Domain,5),
	}
	doms:=conn.DomainDetailFactory(DataHandle,DomainSeeker)
 	// Domain Detail로 채우는 객체 생성
	err:= doms.DomainSeeker.SetDomain()
	if err!= nil{
		http.Error(w, "error while fetcing domain list", http.StatusBadRequest)
	}
	//domain freeing need to be added
	list,_:=doms.DomainSeeker.ReturnDomain()
	for i := range list{
		DataHandle.GetInfo(list[i])
		doms.DataHandle=append(doms.DataHandle, DataHandle)
	}

	data, _ := json.Marshal(doms.DataHandle)
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
	DomainSeeker:= &conn.DomainSeekingByUUID{
		LibvirtInst: i.LibvirtInst,
		UUID: string(param.UUID),
		Domain: make([]*conn.Domain,1),
	}
	outputStruct, _ := conn.DataTypeRouter(param.DataType)
	domain:=conn.DomainDetailFactory(outputStruct,DomainSeeker)
	
	err:=domain.DomainSeeker.SetDomain()
	if err!=nil{
		fmt.Print("error occured while returning status ")
		http.Error(w, "there is no such VM with that UUID", 1)
	}

	dom,_:=domain.DomainSeeker.ReturnDomain()
	outputStruct.GetInfo(dom[0])

	domain.DataHandle=append(domain.DataHandle, outputStruct)
	data, _ := json.Marshal(domain.DataHandle)
	if err != nil {
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}



func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	var param DeleteDomain
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		http.Error(w, "error decoding body", 1)
	}
	DomainSeeker:=&conn.DomainSeekingByUUID{
		LibvirtInst: i.LibvirtInst,
		UUID: param.UUID,
		Domain: make([]*conn.Domain, 0,1),
	}
	DomainTerminator,_:= conn.DomainTerminatorFactory(DomainSeeker)
	domainInfo,err:=DomainTerminator.ShutDownDomain()
	if err!= nil{
		http.Error(w,"error while shutting down Domain", http.StatusBadRequest)
		fmt.Println(err)
	}
	//uuid unparsing 중 에러, destroyDom 에서 에러
	// 수신시 에러 발생 가능 ,추후 에러 핸들링 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":   fmt.Sprintf("VM with UUID %s Shutdown successfully.", param.UUID),
		"domainInfo": domainInfo,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func (i *InstHandler)DeleteVM(w http.ResponseWriter, r *http.Request){
	var param DeleteDomain
	if err:= json.NewDecoder(r.Body).Decode(&param);err != nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	DomainSeeker:= &conn.DomainSeekingByUUID{
		LibvirtInst: i.LibvirtInst,
		UUID: param.UUID,
		Domain: make([]*conn.Domain, 0,1),
	}
	DomainDeleter,_:=conn.DomainDeleterFactory(DomainSeeker, param.DeletionType, param.UUID)
	
	domainInfo, err:=DomainDeleter.DeleteDomain()
	if err!=nil{
		fmt.Println(err)
		http.Error(w,"error while destroying vm", http.StatusInternalServerError)
	}
		//uuid unparsing 중 에러, undefine,destroyDom 에서 에러, 켜져 있지만 softdelete가 
	// 수신시 에러 발생 가능 ,추후 에러 핸들링 

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"message":   fmt.Sprintf("VM with UUID %s Deletion successfully.", param.UUID),
		"domainInfo": domainInfo,
	}	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
	
}

