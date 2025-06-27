package api

import (
	"fmt"
	"net/http"

	virerr "github.com/easy-cloud-Knet/KWS_Core.git/error"
	"github.com/easy-cloud-Knet/KWS_Core.git/vm/parsor"
	"github.com/easy-cloud-Knet/KWS_Core.git/vm/service/creation"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)


func (i *InstHandler) CreateVMFromBase(w http.ResponseWriter, r *http.Request) {

	resp := ResponseGen[libvirt.DomainInfo]("CreateVm")
	param := &parsor.VM_Init_Info{}
	domCon,_:= i.domainConGetter()
	if domCon==nil{
		fmt.Println("emrpy domcon")
	}
	

	if err := HttpDecoder(r, param); err != nil {
		resp.ResponseWriteErr(w, err, http.StatusBadRequest)
		i.Logger.Error("error occured while decoding user's parameter of requested creation")
		return
	}
	i.Logger.Info("Handling Create VM", zap.String("uuid", param.UUID))

	domainExisting,_:=domCon.GetDomain(param.UUID, i.LibvirtInst)
	if (domainExisting!=nil){
		fmt.Println(domainExisting)
		resp.ResponseWriteErr(w, nil, http.StatusBadRequest)
		i.Logger.Error("error handling creating vm, domain already exists", zap.String("uuid",param.UUID))
		return
	}

	DomConf := creation.LocalConfFactory(param, i.Logger)
	DomCreator:=creation.LocalCreatorFactory(DomConf, i.LibvirtInst,i.Logger)
	
	newDomain,err :=DomCreator.CreateVM()
	if err!=nil&& newDomain==nil{
		newErr:=virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf(" %w error while creating new domain, from CreateVM",err))
		i.Logger.Error("error from createvm" , zap.Error(newErr))
		resp.ResponseWriteErr(w, err, http.StatusInternalServerError)
		return		
	}

	domCon.AddNewDomain(newDomain,param.UUID)
	
	resp.ResponseWriteOK(w, nil)
	
}
