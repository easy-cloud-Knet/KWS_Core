package api

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	virerr "github.com/easy-cloud-Knet/KWS_Core.git/error"
	"github.com/easy-cloud-Knet/KWS_Core.git/vm/service/status"
	"go.uber.org/zap"
)

func (i *InstHandler) ReturnStatusUUID(w http.ResponseWriter, r *http.Request) {
	param := &ReturnDomainFromUUID{}
	resp := ResponseGen[status.DataTypeHandler]("domain Status UUID")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		return
	}
	i.Logger.Info("retreving domain status", zap.String("uuid", param.UUID))

	outputStruct, err := status.DataTypeRouter(param.DataType)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		return
		// wrong parameter error 반환
	}
	dom, err := i.DomainControl.GetDomain(param.UUID, i.LibvirtInst)
	if err != nil {
		resp.ResponseWriteErr(w, virerr.ErrorJoin(err, errors.New("error returning status from uuid")), http.StatusInternalServerError)
	}

	DomainDetail := status.DomainDetailFactory(outputStruct, dom)

	outputStruct.GetInfo(dom)
	DomainDetail.DataHandle = outputStruct
	fmt.Println(outputStruct)
	resp.ResponseWriteOK(w, &DomainDetail.DataHandle)

}

func (i *InstHandler) ReturnStatusHost(w http.ResponseWriter, r *http.Request) {
	param := &ReturnHostFromStatus{}
	resp := ResponseGen[status.HostDataTypeHandler]("Host Status Return")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		http.Error(w, "error decoding parameters", http.StatusBadRequest)
		return
	}

	dataHandle, err := status.HostDataTypeRouter(param.HostDataType)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Println("data sending", reflect.TypeOf(dataHandle))
	host, err := status.HostDetailFactory(dataHandle)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
	}

	resp.ResponseWriteOK(w, &host.HostDataHandle)
}

func (i *InstHandler) ReturnInstAllInfo(w http.ResponseWriter, r *http.Request) {
	param := &ReturnInstAllData{}
	resp := ResponseGen[status.InstDataTypeHandler]("Inst Hardware Return")

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		http.Error(w, "error decoding parameters", http.StatusBadRequest)
		return
	}

	dataHandle, err := status.InstDataTypeRouter(param.InstDataType)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Println("data sending", reflect.TypeOf(dataHandle))
	inst, err := status.InstDetailFactory(dataHandle, i.LibvirtInst)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
	}

	resp.ResponseWriteOK(w, &inst.AllInstDataHandle)
}

// host 상태 조회는 두가지 에러를 반환할 수 있음.
// 1. Routing 등에서 일어나는 에러, (host data 타입등이 잘못 입력 된 경우) , InvalidParameter
// 2. 정확히 파악할 수 없는 오류, 사용하는 호스트 반환 패키지에서 반환되는 오류, , HostStatusError
// 이 두가지는 내부 함수에서 파악하여 올리기 때문에, 추가 없이  ResponseWirteErr 호출해도 괜찮을 듯

func (i *InstHandler) ReturnAllUUIDs(w http.ResponseWriter, r *http.Request) {
	i.Logger.Info("ReturnAllUUIDs handler entered")

	resp := ResponseGen[UUIDListResponse]("Get All UUIDs")

	uuids := i.DomainControl.GetAllUUIDs()
	respData := UUIDListResponse{UUIDs: uuids}

	resp.ResponseWriteOK(w, &respData)
}
