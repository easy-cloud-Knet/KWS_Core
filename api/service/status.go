package service

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
	"go.uber.org/zap"
)




func (i *InstHandler) ReturnStatusUUID(w http.ResponseWriter, r *http.Request) {
	param:=&ReturnDomainFromUUID{}
	resp:=ResponseGen[conn.DataTypeHandler]("domain Status UUID")

	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,err, http.StatusBadRequest)
		return
	}
	i.Logger.Info("retreving domain status",
	 zap.String("uuid", param.UUID),
	zap.Int("method", int(param.DataType)))

	dom, err:= i.DomainControl.GetDomain(param.UUID, i.LibvirtInst, i.Logger)
	if err!=nil{
		detailErr := virerr.ErrorGen(virerr.DomainStatusError,fmt.Errorf("error getting domain while serving ReturnStatusUUID ,UUID of %s, %w",param.UUID, err))
		resp.ResponseWriteErr(w,detailErr, http.StatusInternalServerError)
		return
	}

	outputStruct, err := conn.DataTypeRouter(param.DataType)
	if err!=nil{
		detailErr := virerr.ErrorJoin(err,fmt.Errorf("error domain type routing while serving ReturnStatusUUID ,UUID of %s, %w",param.UUID, err))
		resp.ResponseWriteErr(w, detailErr, http.StatusBadRequest)
		return
	}

	DomainDetail := conn.DomainDetailFactory(outputStruct, dom)


	outputStruct.GetInfo(dom)
	i.Logger.Sugar().Info("retreving domain info successfully", outputStruct)
	DomainDetail.DataHandle = outputStruct
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
	fmt.Println("data sending", reflect.TypeOf(dataHandle))
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
