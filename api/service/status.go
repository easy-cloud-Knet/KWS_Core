package service

import (
	"fmt"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)


func (i *InstHandler) ReturnDomainByStatus(w http.ResponseWriter, r *http.Request) {
	param:=&ReturnDomainFromStatus{}
	resp:=ResponseGen[[]conn.DataTypeHandler]("domain status with Status")
	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error booting vm",err), http.StatusInternalServerError)
		return
	}
	
	DomainSeeker:=conn.DomSeekStatusFactory(i.LibvirtInst, param.Status)
	DataHandle, _ := conn.DataTypeRouter(param.DataType)
	doms := conn.DomainDetailFactory(DataHandle, DomainSeeker)
	// Domain Detail로 채우는 객체 생성
 
	list, err := doms.DomainSeeker.ReturnDomain()
	if err != nil {
		resp.ResponseWriteErr(w,fmt.Errorf("%w error booting vm",err), http.StatusInternalServerError)
		return
	}
	for i := range list {
		if list[i] == nil {
			resp.ResponseWriteErr(w,fmt.Errorf("%w error booting vm",err), http.StatusInternalServerError)
			return
		}
		DataHandle.GetInfo(list[i])
		doms.DataHandle = append(doms.DataHandle, DataHandle)
	}
	resp.ResponseWriteOK(w, &doms.DataHandle)
	//slice를 포인터로 넘기는 문제가 있음, 타입 어설션 등 고민,, 
}

func (i *InstHandler) ReturnStatusUUID(w http.ResponseWriter, r *http.Request) {
	param:=&ReturnDomainFromUUID{}
	resp:=ResponseGen[[]conn.DataTypeHandler]("domain Status UUID")

	if err:=HttpDecoder(r,param); err!=nil{
		http.Error(w, "error decoding parameters", http.StatusBadRequest)
		return
	}

	DomainSeeker:= conn.DomSeekUUIDFactory(i.LibvirtInst, param.UUID)
	outputStruct, _ := conn.DataTypeRouter(param.DataType)

	DomainDetail := conn.DomainDetailFactory(outputStruct, DomainSeeker)

	dom, err := DomainDetail.DomainSeeker.ReturnDomain()
	if err != nil {
		resp.ResponseWriteErr(w,fmt.Errorf("%w error booting vm",err), http.StatusInternalServerError)
		return
	}
	fmt.Println(dom)
	outputStruct.GetInfo(dom[0])
	DomainDetail.DataHandle = append(DomainDetail.DataHandle, outputStruct)

	resp.ResponseWriteOK(w, &DomainDetail.DataHandle)

}

func (i *InstHandler) ReturnStatusHost(w http.ResponseWriter, r *http.Request) {
	param:=&ReturnHostFromStatus{}
	resp:=ResponseGen[conn.HostDataTypeHandler]("Host Status Return")

	if err:=HttpDecoder(r,param); err!=nil{
		http.Error(w, "error decoding parameters", http.StatusBadRequest)
		return
	}
 
	dataHandle, err := conn.HostDataTypeRouter(param.HostDataType)
	if err != nil {
		http.Error(w, "fail host data type", http.StatusInternalServerError)
		return
	}

	host := conn.HostDetailFactory(dataHandle)
	
	resp.ResponseWriteOK(w, &host.HostDataHandle)

}
