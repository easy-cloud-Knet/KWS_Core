package api

import (
	"fmt"
	"net/http"

	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"github.com/easy-cloud-Knet/KWS_Core/vm/parsor"
	"github.com/easy-cloud-Knet/KWS_Core/vm/service/creation"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

func (i *InstHandler) BootVM(w http.ResponseWriter, r *http.Request) {
	resp := ResponseGen[libvirt.DomainInfo]("BootVM")
	param := &StartDomain{}

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		i.Logger.Error("error occured while decoding user's parameter of requested creation")
		return
	}
	i.Logger.Info("Handling Boot VM", zap.String("uuid", param.UUID))

	domCon, _ := i.domainConGetter()

	DomainExisting, _ := domCon.GetDomain(param.UUID, i.LibvirtInst)
	if DomainExisting == nil {
		resp.ResponseWriteErr(w, nil, http.StatusBadRequest)
		i.Logger.Error("error handling booting vm, domain not found", zap.String("uuid", param.UUID))
		return
	}

	err := DomainExisting.Domain.Create()
	if err != nil {
		newErr := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf(" %w error while booting domain, from BootVM", err))
		i.Logger.Error("error from booting vm", zap.Error(newErr))
		resp.ResponseWriteErr(w, newErr, http.StatusInternalServerError)
		return
	}

	resp.ResponseWriteOK(w, nil)
	i.Logger.Info("Boot VM request handled successfully", zap.String("uuid", param.UUID))
}

func (i *InstHandler) CreateVMFromBase(w http.ResponseWriter, r *http.Request) {

	resp := ResponseGen[libvirt.DomainInfo]("CreateVm")
	param := &parsor.VM_Init_Info{}

	domCon, _ := i.domainConGetter()
	if domCon == nil {
		fmt.Println("emrpy domcon")
	}

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		i.Logger.Error("error occured while decoding user's parameter of requested creation")
		return
	}
	i.Logger.Info("Handling Create VM", zap.String("uuid", param.UUID))

	// 같은 UUID를 가진 Instance가 있는지 확인
	domainExisting, _ := domCon.GetDomain(param.UUID, i.LibvirtInst)
	if domainExisting != nil {
		fmt.Println(domainExisting)
		resp.ResponseWriteErr(w, nil, http.StatusBadRequest)
		i.Logger.Error("error handling creating vm, domain already exists", zap.String("uuid", param.UUID))
		return
	}

	// Instance 생성에 필요한 객체 생성
	DomConf := creation.LocalConfFactory(param, i.Logger)
	DomCreator := creation.LocalCreatorFactory(DomConf, i.LibvirtInst, i.Logger)

	// Instance 생성 시도
	newDomain, err := DomCreator.CreateVM()
	if err != nil && newDomain == nil {
		newErr := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf(" %w error while creating new domain, from CreateVM", err))
		i.Logger.Error("error from createvm", zap.Error(newErr))
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return
	}

	// 도메인 컨트롤러에 만들어진 도메인 등록
	domCon.AddNewDomain(newDomain, param.UUID)

	resp.ResponseWriteOK(w, nil)

}
