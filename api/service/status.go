package service

import (
	"errors"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)


func (i *InstHandler) ReturnDomainByStatus(w http.ResponseWriter, r *http.Request) {
	param:=&ReturnDomainFromStatus{}
	resp:=ResponseGen[[]conn.DataTypeHandler]("domain status with Status")
	
	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,err, http.StatusBadRequest)
		//StatusInternalError 말고 다른 에러 반환 고민
		return
	}

	DomainSeeker:=conn.DomSeekStatusFactory(i.LibvirtInst, param.Status)
	DataHandle, err := conn.DataTypeRouter(param.DataType)
	if err!=nil{
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		return
		// wrong parameter error 반환
	}
	doms := conn.DomainDetailFactory(DataHandle, DomainSeeker)
	// Domain Detail로 채우는 객체 생성
 
	list, err := doms.DomainSeeker.ReturnDomain()
	if err != nil {
		appendingErorr:=conn.ErrorJoin(err, errors.New("retreving Domain from Return Domain Status"))
		resp.ResponseWriteErr(w, appendingErorr, http.StatusInternalServerError)
		return 
	}
	for i := range list {
		if list[i] == nil {
			appendingError:= conn.ErrorGen(conn.DomainSearchError, errors.New("internal error dereferncing domain pointer in Sercing Status"))
			resp.ResponseWriteErr(w,appendingError, http.StatusInternalServerError)
			return
		}
		err:= DataHandle.GetInfo(list[i])
		if err!=nil{
			appendingErorr:=conn.ErrorJoin(conn.DomainStatusError, errors.New("retreving Domain Status Error"))
			resp.ResponseWriteErr(w,appendingErorr, http.StatusInternalServerError)
		}
		doms.DataHandle = append(doms.DataHandle, DataHandle)
	}
	resp.ResponseWriteOK(w, &doms.DataHandle)
	//slice를 포인터로 넘기는 문제가 있음, 타입 어설션 등 고민,, 
}
// InvalidUUID, UUID 변환 중 오류가 발생할 시
// NoSuchDomain, 도메인 검색 중, 해당 도메인을 찾을 수 없을 시
// DomainStatusError, 알수 없는 에러, Libvirt 내부 오류 발생 시
// Wrong paramter, 도메인 상태 플래그가 범위를 벗어났을 시

func (i *InstHandler) ReturnStatusUUID(w http.ResponseWriter, r *http.Request) {
	param:=&ReturnDomainFromUUID{}
	resp:=ResponseGen[[]conn.DataTypeHandler]("domain Status UUID")

	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,err, http.StatusBadRequest)
		return
	}

	DomainSeeker:= conn.DomSeekUUIDFactory(i.LibvirtInst, param.UUID)
	outputStruct, err := conn.DataTypeRouter(param.DataType)
	if err!=nil{
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		return
		// wrong parameter error 반환
	}

	DomainDetail := conn.DomainDetailFactory(outputStruct, DomainSeeker)

	dom, err := DomainDetail.DomainSeeker.ReturnDomain()
	if err != nil {
		appendingErorr:=conn.ErrorJoin(err, errors.New("retreving Domain from Return Domain Status funcion error"))
		resp.ResponseWriteErr(w, appendingErorr, http.StatusInternalServerError)
		return 
	}
	outputStruct.GetInfo(dom[0])
	DomainDetail.DataHandle = append(DomainDetail.DataHandle, outputStruct)

	resp.ResponseWriteOK(w, &DomainDetail.DataHandle)

}

func (i *InstHandler) ReturnStatusHost(w http.ResponseWriter, r *http.Request) {
	param:=&ReturnHostFromStatus{}
	resp:=ResponseGen[conn.HostDataTypeHandler]("Host Status Return")

	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,err,http.StatusInternalServerError)
		http.Error(w, "error decoding parameters", http.StatusBadRequest)
		return
	}
 
	dataHandle, err := conn.HostDataTypeRouter(param.HostDataType)
	if err != nil {
		resp.ResponseWriteErr(w,err,http.StatusInternalServerError)
		return
	}
	
	host,err := conn.HostDetailFactory(dataHandle)
	if err!= nil{
		resp.ResponseWriteErr(w,err,http.StatusInternalServerError)
	 }
	
	resp.ResponseWriteOK(w, &host.HostDataHandle)
}
// host 상태 조회는 두가지 에러를 반환할 수 있음.
// 1. Routing 등에서 일어나는 에러, (host data 타입등이 잘못 입력 된 경우) , InvalidParameter 
// 2. 정확히 파악할 수 없는 오류, 사용하는 호스트 반환 패키지에서 반환되는 오류, , HostStatusError 
// 이 두가지는 내부 함수에서 파악하여 올리기 때문에, 추가 없이  ResponseWirteErr 호출해도 괜찮을 듯
