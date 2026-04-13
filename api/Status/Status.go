package status

import (
	"errors"
	"net/http"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	httputil "github.com/easy-cloud-Knet/KWS_Core/pkg/httputil"
	svcstatus "github.com/easy-cloud-Knet/KWS_Core/services/status"
	"go.uber.org/zap"
)

func (h *Handler) ReturnStatusUUID(w http.ResponseWriter, r *http.Request) {
	param := &DomainStatusRequest{}
	resp := httputil.ResponseGen[svcstatus.DataTypeHandler]("domain Status UUID")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		return
	}
	h.Logger.Info("retreving domain status", zap.String("uuid", param.UUID))

	outputStruct, err := svcstatus.DataTypeRouter(param.DataType)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		return
	}
	dom, err := h.DomainControl.GetDomain(param.UUID)
	if err != nil {
		resp.ResponseWriteErr(w, virerr.ErrorJoin(err, errors.New("error returning status from uuid")), http.StatusInternalServerError)
		return
	}

	DomainDetail := svcstatus.DomainDetailFactory(outputStruct, dom.Domain)

	if err := outputStruct.GetInfo(dom.Domain); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}
	DomainDetail.DataHandle = outputStruct
	resp.ResponseWriteOK(w, &DomainDetail.DataHandle)
}

func (h *Handler) ReturnStatusHost(w http.ResponseWriter, r *http.Request) {
	param := &HostStatusRequest{}
	resp := httputil.ResponseGen[svcstatus.HostDataTypeHandler]("Host Status Return")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}

	dataHandle, err := svcstatus.HostDataTypeRouter(param.HostDataType)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}

	host, err := svcstatus.HostInfoHandler(dataHandle, h.DomainControl.GetDomainListStatus())
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}

	resp.ResponseWriteOK(w, &host.HostDataHandle)
}

func (h *Handler) ReturnInstAllInfo(w http.ResponseWriter, r *http.Request) {
	param := &InstInfoRequest{}
	resp := httputil.ResponseGen[svcstatus.InstDataTypeHandler]("Inst Hardware Return")

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}

	dataHandle, err := svcstatus.InstDataTypeRouter(param.InstDataType)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}
	inst, err := svcstatus.InstDetailFactory(dataHandle, h.LibvirtInst)
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}

	resp.ResponseWriteOK(w, &inst.AllInstDataHandle)
}

func (h *Handler) ReturnAllUUIDs(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("ReturnAllUUIDs handler entered")

	resp := httputil.ResponseGen[UUIDListResponse]("Get All UUIDs")

	uuids := h.DomainControl.GetAllUUIDs()
	respData := UUIDListResponse{UUIDs: uuids}

	resp.ResponseWriteOK(w, &respData)
}

func (h *Handler) getAllDomainStates() ([]DomainStateResponse, error) {
	domains, err := h.LibvirtInst.ListAllDomains(0)
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, d := range domains {
			d.Free()
		}
	}()

	var result []DomainStateResponse
	for _, domain := range domains {
		uuid, err := domain.GetUUIDString()
		if err != nil {
			continue
		}
		state, _, err := domain.GetState()
		if err != nil {
			continue
		}
		result = append(result, DomainStateResponse{
			UUID:        uuid,
			DomainState: state,
		})
	}
	return result, nil
}

func (h *Handler) ReturnAllDomainStates(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("ReturnAllDomainStates handler entered")

	resp := httputil.ResponseGen[[]DomainStateResponse]("Get All Domain States")

	states, err := h.getAllDomainStates()
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}

	resp.ResponseWriteOK(w, &states)
}
