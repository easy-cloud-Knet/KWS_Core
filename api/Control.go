package api

import (
	"fmt"
	"net/http"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	domainStatus "github.com/easy-cloud-Knet/KWS_Core/DomCon/domain_status"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"github.com/easy-cloud-Knet/KWS_Core/vm/service/termination"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func (i *InstHandler) ForceShutDownVM(w http.ResponseWriter, r *http.Request) {
	param := &DomainControlRequest{}
	resp := ResponseGen[any]("domain number of" + param.UUID + ", Force Shutdown VM")

	if err := HttpDecoder(r, param); err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error shutting down vm, from forceShutdown vm "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to decode forceShutdown request", zap.Error(ERR))
		return
	}
	dom, err := i.DomainControl.GetDomain(param.UUID)
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error shutting down vm, retreving Get domin error "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to get domain for forceShutdown", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
	}

	DomainTerminator := termination.DomainTerminatorFactory(dom.Domain)

	err = DomainTerminator.TerminateDomain()
	if err != nil {
		resp.ResponseWriteErr(w, virerr.ErrorJoin(err, fmt.Errorf("error shutting down vm, retreving Get domin error ")), http.StatusInternalServerError)
		return
	}

	stat, err := i.DomainControl.DomainListStatus.GetDomStatus(dom.Domain, []domainStatus.SourceType{domainStatus.CPU}, i.Logger)
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error getting domain status for forceShutdown"))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to get domain status for forceShutdown", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
	}
	i.Logger.Info("Domain status retrieved", zap.Any("status", stat))

	i.DomainControl.DomainListStatus.AddSleepingCPU(int(stat.(map[domainStatus.SourceType]int)[domainStatus.CPU]))

	resp.ResponseWriteOK(w, nil)
}

func (i *InstHandler) DeleteVM(w http.ResponseWriter, r *http.Request) {
	param := &DomainControlRequest{}
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

	domain, err := i.DomainControl.GetDomain(param.UUID)
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error deleting vm, retreving Get domin error "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to get domain for deleteVM", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
		//error handling
	}

	stat, err := i.DomainControl.DomainListStatus.GetDomStatus(domain.Domain, []domainStatus.SourceType{domainStatus.CPU}, i.Logger)
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error getting domain status for deleteVM"))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to get domain status for deleteVM", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
	}
	i.Logger.Info("Domain status retrieved", zap.Any("status", stat))

	DomainDeleter := termination.DomainDeleterFactory(domain.Domain, param.DeletionType, param.UUID)
	if err := DomainDeleter.DeleteDomain(); err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error deleting vm, retreving Get domin error "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to delete domain", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
	}
	i.DomainControl.DeleteDomain(domain.Domain, param.UUID, int(stat.(map[domainStatus.SourceType]int)[domainStatus.CPU]))

	resp.ResponseWriteOK(w, nil)
}
