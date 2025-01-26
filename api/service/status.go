package service

import (
	"encoding/json"
	"fmt"
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

	DomainDetail:=conn.DomainDetailFactory(outputStruct,DomainSeeker)
	
	err:=DomainDetail.DomainSeeker.SetDomain()
	if err!=nil{
		fmt.Print("error occured while returning status ")
		http.Error(w, "there is no such VM with that UUID", 1)
	}
	
	dom,_:=DomainDetail.DomainSeeker.ReturnDomain()
	fmt.Println(dom)

	outputStruct.GetInfo(dom[0])

	DomainDetail.DataHandle=append(DomainDetail.DataHandle, outputStruct)
	data, _ := json.Marshal(DomainDetail.DataHandle)
	if err != nil {
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

