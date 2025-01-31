package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)



func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var param DeleteDomain
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		CommonErrorHelper(w,err,http.StatusBadRequest, "error Decoding parameters ")

		return
	}
	DomainSeeker:=&conn.DomainSeekingByUUID{
		LibvirtInst: i.LibvirtInst,
		UUID: param.UUID,
		Domain: make([]*conn.Domain, 0,1),
	}
	DomainTerminator,_:= conn.DomainTerminatorFactory(DomainSeeker)

	domainInfo,err:=DomainTerminator.ShutDownDomain()
	if err!= nil{
		CommonErrorHelper(w,err,http.StatusInternalServerError, "error while shutting down domain")
		return
	}
	//uuid unparsing 중 에러, destroyDom 에서 에러
	// 수신시 에러 발생 가능 ,추후 에러 핸들링 
	response := map[string]interface{}{
		"message":   fmt.Sprintf("VM with UUID %s Shutdown successfully.", param.UUID),
		"domainInfo": domainInfo,
	}	
	resp, err:=json.Marshal(response)
	if err!= nil{
		CommonErrorHelper(w,err,http.StatusInternalServerError, "error while Marshaling response ")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}




func (i *InstHandler)DeleteVM(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var param DeleteDomain
	if err:= json.NewDecoder(r.Body).Decode(&param);err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w,"Error while decoding deletion parameter, %v", err)
		return
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w,"Error while deleting Domain %v", err)
		return
	}
		//uuid unparsing 중 에러, undefine,destroyDom 에서 에러, 켜져 있지만 softdelete가 
	// 수신시 에러 발생 가능 ,추후 에러 핸들링 

	response := map[string]interface{}{
		"message":   fmt.Sprintf("VM with UUID %s Deletion successfully.", param.UUID),
		"domainInfo": domainInfo,
	}	
	
	resp, err:=json.Marshal(response)
	if err!= nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w,"error while Marshaling response %v",err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

