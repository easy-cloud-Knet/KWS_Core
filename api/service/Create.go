package service

import (
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn/creation"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)


func (i *InstHandler) CreateVMFromBase(w http.ResponseWriter, r *http.Request) {

	resp := ResponseGen[libvirt.DomainInfo]("CreateVm")
	param := &parsor.VM_Init_Info{}

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		i.Logger.Error("error occured while decoding user's parameter of requested creation")
		return
	}
	i.Logger.Info("Handling Create VM", zap.String("uuid", param.UUID))

	DomConf := creation.LocalDomainerFactory(param, i.Logger)
	DomCreator:=creation.DomainCreatorFactory(DomConf, i.LibvirtInst,i.Logger)
	domController := conn.DomainControllerInjection(i.DomainControl,DomCreator)
	
	err:=domController.DomainAddWithOperation(i.Logger,param.UUID)
	if err != nil {
		i.Logger.Error("domain Generating failed, while encoding user's config", zap.String("uuid", param.UUID))
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}
	
	resp.ResponseWriteOK(w, nil)
}
