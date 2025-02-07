package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)

func (i *InstHandler) ReturnDomainByStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var param ReturnDomainFromStatus
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		CommonErrorHelper(w, err, http.StatusBadRequest, "error Decoding parameters ")
		return
	}

	DataHandle, _ := conn.DataTypeRouter(param.DataType)
	DomainSeeker := &conn.DomainSeekingByStatus{
		LibvirtInst: i.LibvirtInst,
		Status:      param.Status,
		DomList:     make([]*conn.Domain, 5),
	}
	doms := conn.DomainDetailFactory(DataHandle, DomainSeeker)
	// Domain Detail로 채우는 객체 생성
	err := doms.DomainSeeker.SetDomain()
	if err != nil {
		CommonErrorHelper(w, err, http.StatusInternalServerError, "error while setting Domain")
		return
	}
	//domain freeing need to be added
	list, err := doms.DomainSeeker.ReturnDomain()
	if err != nil {
		CommonErrorHelper(w, err, http.StatusInternalServerError, "error while retreving Domain")
		return
	}
	for i := range list {
		if list[i] == nil {
			CommonErrorHelper(w, err, http.StatusInternalServerError, "invalid Domain Object contained ")
			return
		}
		DataHandle.GetInfo(list[i])
		doms.DataHandle = append(doms.DataHandle, DataHandle)
	}

	data, err := json.Marshal(doms.DataHandle)
	if err != nil {
		CommonErrorHelper(w, err, http.StatusInternalServerError, "failed marshalling data ")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (i *InstHandler) ReturnStatusUUID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var param ReturnDomainFromUUID
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		CommonErrorHelper(w, err, http.StatusBadRequest, "Invalid Parameter")
		return
	}
	DomainSeeker := &conn.DomainSeekingByUUID{
		LibvirtInst: i.LibvirtInst,
		UUID:        string(param.UUID),
		Domain:      make([]*conn.Domain, 1),
	}
	outputStruct, _ := conn.DataTypeRouter(param.DataType)

	DomainDetail := conn.DomainDetailFactory(outputStruct, DomainSeeker)

	err := DomainDetail.DomainSeeker.SetDomain()
	if err != nil {
		fmt.Println("doamin error %w", err)
		CommonErrorHelper(w, err, http.StatusInternalServerError, "error occured while returning status")
		return
	}

	dom, err := DomainDetail.DomainSeeker.ReturnDomain()
	if err != nil {
		CommonErrorHelper(w, err, http.StatusBadRequest, "error occured while returning status")
		return
	}
	outputStruct.GetInfo(dom[0])

	DomainDetail.DataHandle = append(DomainDetail.DataHandle, outputStruct)
	data, err := json.Marshal(DomainDetail.DataHandle)
	if err != nil {
		CommonErrorHelper(w, err, http.StatusInternalServerError, "failed to marshal JSON")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (i *InstHandler) ReturnStatusHost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error read body", http.StatusBadRequest)
		return
	}
	var param ReturnHostFromStatus

	if err := json.Unmarshal(body, &param); err != nil {
		http.Error(w, "error decoding parameters", http.StatusBadRequest)
		return
	}

	dataHandle, err := conn.HostDataTypeRouter(param.HostDataType)
	if err != nil {
		http.Error(w, "fail host data type", http.StatusInternalServerError)
		return
	}

	host := conn.HostDetailFactory(dataHandle)

	data, err := json.Marshal(host.HostDataHandle)
	if err != nil {
		CommonErrorHelper(w, err, http.StatusInternalServerError, "failed marshalling data")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
