package api

import (
	"fmt"
	"net/http"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"github.com/easy-cloud-Knet/KWS_Core/vm/service/termination"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func (i *InstHandler) ForceShutDownVM(w http.ResponseWriter, r *http.Request) {
	param := &DeleteDomain{}
	resp := ResponseGen[any]("domain number of" + param.UUID + ", Force Shutdown VM")

	if err := HttpDecoder(r, param); err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error shutting down vm, from forceShutdown vm "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to decode forceShutdown request", zap.Error(ERR))
		return
	}
	dom, err := i.DomainControl.GetDomain(param.UUID, i.LibvirtInst)
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error shutting down vm, retreving Get domin error "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to get domain for forceShutdown", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
	}
	vcpu, err := dom.Domain.GetMaxVcpus()
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error shutting down vm, retreving Get domin error "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to get vcpu count for forceShutdown", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
	}
	i.DomainControl.DomainListStatus.AddSleepingCPU(int(vcpu))

	DomainTerminator, _ := termination.DomainTerminatorFactory(dom)

	_, err = DomainTerminator.TerminateDomain()
	if err != nil {
		resp.ResponseWriteErr(w, virerr.ErrorJoin(err, fmt.Errorf("error shutting down vm, retreving Get domin error ")), http.StatusInternalServerError)
		return
	}

	resp.ResponseWriteOK(w, nil)
}

func (i *InstHandler) DeleteVM(w http.ResponseWriter, r *http.Request) {
	param := &DeleteDomain{}
	resp := ResponseGen[libvirt.DomainInfo]("Deleting Vm")

	if err := HttpDecoder(r, param); err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error deleting vm, unparsing HTTP request "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to decode deleteVM request", zap.Error(ERR))
		return
	}
	if _, err := domCon.ReturnUUID(param.UUID); err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error deleting vm,	invalid UUID "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("invalid UUID for deleteVM", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
	}
	// uuid 가 적합한지 확인

	domain, err := i.DomainControl.GetDomain(param.UUID, i.LibvirtInst)
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error deleting vm, retreving Get domin error "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to get domain for deleteVM", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
		//error handling
	}

	vcpu, err := domain.Domain.GetMaxVcpus()
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error can't retreving vcpu count "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to get vcpu count for deleteVM", zap.String("uuid", param.UUID), zap.Error(ERR))

		vcpu = 2
		//return
		//일단 지금은 해당 경우에 vcpu 숫자를 2로 설정
	} // 삭제된 도메인에서는 vcpu count 를 가져올 수 없으므로 미리 가져옴 . 맘에 안듦. 나중에 수정할 예정
	// TODO: GETMAXVCPU는 꺼진 도메인에 대해 동작하지 않음. DATADOG와 같은 인터페이스를 활용해서 상관없이 삭제할 수 있도록

	DomainDeleter, _ := termination.DomainDeleterFactory(domain, param.DeletionType, param.UUID)
	domDeleted, err := DomainDeleter.DeleteDomain()
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error deleting vm, retreving Get domin error "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to delete domain", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
	}
	i.DomainControl.DeleteDomain(domDeleted, param.UUID, int(vcpu))

	resp.ResponseWriteOK(w, nil)
}
