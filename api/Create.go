package api

import (
	"errors"
	"fmt"
	"net/http"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	httputil "github.com/easy-cloud-Knet/KWS_Core/pkg/httputil"
	"github.com/easy-cloud-Knet/KWS_Core/services/creation"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func (i *InstHandler) BootVM(w http.ResponseWriter, r *http.Request) {
	resp := httputil.ResponseGen[libvirt.DomainInfo]("BootVM")
	param := &DomainBootRequest{}

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		i.Logger.Error("error occured while decoding user's parameter of requested creation")
		return
	}
	i.Logger.Info("Handling Boot VM", zap.String("uuid", param.UUID))

	DomainExisting, domainErr := i.DomainControl.GetDomain(param.UUID)
	if domainErr != nil {
		i.Logger.Error("error handling booting vm, failed to get domain", zap.String("uuid", param.UUID), zap.Error(domainErr))
		resp.ResponseWriteErr(w, domainErr, http.StatusInternalServerError)
		return
	}
	if DomainExisting == nil {
		notFoundErr := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("domain %s not found while booting vm", param.UUID))
		i.Logger.Error("error handling booting vm, domain not found", zap.String("uuid", param.UUID), zap.Error(notFoundErr))
		resp.ResponseWriteErr(w, notFoundErr, http.StatusNotFound)
		return
	}

	err := DomainExisting.Domain.Create()
	if err != nil {
		newErr := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf(" %w error while booting domain, from BootVM", err))
		i.Logger.Error("error from booting vm", zap.Error(newErr))
		resp.ResponseWriteErr(w, newErr, http.StatusInternalServerError)
		return
	}

	vcpu, err := DomainExisting.Domain.GetMaxVcpus()
	if err != nil {
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}
	i.DomainControl.DomainListStatus.TakeSleepingCPU(int(vcpu))

	resp.ResponseWriteOK(w, nil)
	i.Logger.Info("Boot VM request handled successfully", zap.String("uuid", param.UUID))
}

func (i *InstHandler) CreateVMFromBase(w http.ResponseWriter, r *http.Request) {

	resp := httputil.ResponseGen[libvirt.DomainInfo]("CreateVm")
	param := &CreateVMRequest{}

	if err := httputil.HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		i.Logger.Error("error occured while decoding user's parameter of requested creation")
		return
	}
	i.Logger.Info("Handling Create VM", zap.String("uuid", param.UUID))

	dom, err := i.DomainControl.GetDomain(param.UUID)
	if err != nil {
		if errors.Is(err, virerr.DomainSearchError) {
			i.Logger.Error("error handling creating vm, failed to get existing domain", zap.String("uuid", param.UUID), zap.Error(err))
			resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
			return
		} else if errors.Is(err, virerr.NoSuchDomain) {
			i.Logger.Info("no existing domain found with the same uuid, proceeding to create new domain", zap.String("uuid", param.UUID))
		} else {
			i.Logger.Error("unexpected error while checking existing domain", zap.String("uuid", param.UUID), zap.Error(err))
			resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
			return
		}
	}
	if dom != nil {
		i.Logger.Error("existing domain found with the same uuid, You cannot create a new domain with the same uuid", zap.String("uuid", param.UUID))
		resp.ResponseWriteErr(w, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("domain with uuid %s already exists", param.UUID)), http.StatusConflict)
		return
	}

	DomConf := creation.LocalConfFactory(param.toVMInitInfo(), i.Logger)
	DomCreator := creation.LocalCreatorFactory(DomConf, i.LibvirtInst, i.Logger)

	newDomain, err := DomCreator.CreateVM()
	if err != nil && newDomain == nil {
		newErr := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf(" %w error while creating new domain, from CreateVM", err))
		i.Logger.Error("error from createvm", zap.Error(newErr))
		resp.ResponseWriteErr(w, newErr, http.StatusInternalServerError)
		return
	}

	err = i.DomainControl.AddNewDomain(newDomain, param.UUID)
	if err != nil {
		i.Logger.Error("error from createvm", zap.Error(err))
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}

	resp.ResponseWriteOK(w, nil)

}
