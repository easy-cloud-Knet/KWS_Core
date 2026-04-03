package api

import (
	"fmt"
	"net/http"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	instatus "github.com/easy-cloud-Knet/KWS_Core/internal/status"
	httputil "github.com/easy-cloud-Knet/KWS_Core/pkg/httputil"
	"github.com/easy-cloud-Knet/KWS_Core/services/termination"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

// DI pattern
// 실행 시점에 삽입 (test, libvirt.Connect, afadsfadf)

func (i *InstHandler) ForceShutDownVM(w http.ResponseWriter, r *http.Request) {
	param := &DomainControlRequest{}
	resp := httputil.ResponseGen[any]("domain number of" + param.UUID + ", Force Shutdown VM")

	if err := httputil.HttpDecoder(r, param); err != nil {
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

	sources := map[instatus.SourceType]int{instatus.CPU: 0}
	stat, err := i.DomainControl.DomainListStatus.GetDomStatus(dom.Domain, sources, i.Logger)
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error getting domain status for forceShutdown"))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to get domain status for forceShutdown", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
	}
	i.Logger.Info("Domain status retrieved", zap.Any("status", stat))

	i.DomainControl.DomainListStatus.AddSleepingCPU(stat[instatus.CPU])

	resp.ResponseWriteOK(w, nil)
}

func (i *InstHandler) DeleteVM(w http.ResponseWriter, r *http.Request) {
	param := &DomainControlRequest{}
	resp := httputil.ResponseGen[libvirt.DomainInfo]("Deleting Vm")

	if err := httputil.HttpDecoder(r, param); err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error deleting vm, unparsing HTTP request "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to decode deleteVM request", zap.Error(ERR))
		return
	}
	domain, err := i.DomainControl.GetDomain(param.UUID)
	if err != nil {
		ERR := virerr.ErrorJoin(err, fmt.Errorf("error deleting vm, retreving Get domin error "))
		resp.ResponseWriteErr(w, ERR, http.StatusInternalServerError)
		i.Logger.Error("failed to get domain for deleteVM", zap.String("uuid", param.UUID), zap.Error(ERR))
		return
	}

	sources := map[instatus.SourceType]int{instatus.CPU: 0}
	stat, err := i.DomainControl.DomainListStatus.GetDomStatus(domain.Domain, sources, i.Logger)
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
	i.DomainControl.DeleteDomain(domain.Domain, param.UUID, stat[instatus.CPU])

	resp.ResponseWriteOK(w, nil)
}
