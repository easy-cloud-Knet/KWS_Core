package service

import (
	"fmt"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)



func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	
	param:= &DeleteDomain{}
	resp:=ResponseGen[libvirt.DomainInfo]("domain number of"+param.UUID+", Force Shutdown VM")

	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,virerr.ErrorGen(virerr.DomainShutdownError,err), http.StatusInternalServerError)
		i.Logger.Info("error deparsing http request whilie serving ForceShutdown", zap.Error(err))
		return
	}
	i.Logger.Info("http reqeust for shutting down VM", zap.String("uuid",param.UUID))

	dom,err:= i.DomainControl.GetDomain(param.UUID,i.LibvirtInst)
	if err!= nil{
		resp.ResponseWriteErr(w,virerr.ErrorGen(virerr.DomainShutdownError,fmt.Errorf("error retreiving domiain, while serving ForceShutDown VM")), http.StatusInternalServerError)
		i.Logger.Error("error retreiving domiain, while serving ForceShutDown VM", zap.Error(err))
		return
	}

	DomainTerminator,_:= conn.DomainTerminatorFactory(dom)

	domainInfo,err:=DomainTerminator.ShutDownDomain()
	if err!= nil{
		resp.ResponseWriteErr(w,fmt.Errorf("error shutting down Domain, %w ",err), http.StatusInternalServerError)
		i.Logger.Error("error shutting down Domain", zap.Error(err))
		return
	}

 
	resp.ResponseWriteOK(w,domainInfo)
	}




func (i *InstHandler)DeleteVM(w http.ResponseWriter, r *http.Request){
	param:=&DeleteDomain{}
	resp:= ResponseGen[libvirt.DomainInfo]("Deleting Vm")

	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,virerr.ErrorGen(virerr.DomainShutdownError,err), http.StatusInternalServerError)
		i.Logger.Error("error deparsing http request whilie serving DeleteVM", zap.Error(err))
		return
	}
	if _, err := conn.ReturnUUID(param.UUID); err!=nil{
		detailErr := virerr.ErrorGen(virerr.DeletionDomainError,fmt.Errorf("error translating uuid while serving DeleteVM,UUID of %s, %w",param.UUID, err))
		resp.ResponseWriteErr(w,detailErr, http.StatusBadRequest)
		i.Logger.Info("error deparsing http request whilie serving DeleteVM", zap.Error(err))
		return
	}
	// uuid 가 적합한지 확인

	domain, err := i.DomainControl.GetDomain(param.UUID,i.LibvirtInst)
	if err!=nil{
		detailErr := virerr.ErrorGen(virerr.DeletionDomainError,fmt.Errorf("error getting domain while serving DeleteVM,UUID of %s, %w",param.UUID, err))
		resp.ResponseWriteErr(w,detailErr, http.StatusBadRequest)
		return
		//error handling 
	}

	DomainDeleter,_:=conn.DomainDeleterFactory(domain, param.DeletionType)
	domainInfo, err:=DomainDeleter.DeleteDomain(param.UUID)
	if err!=nil{
		encapsuledErr:= virerr.ErrorJoin(err,fmt.Errorf("error while Serving DeleteVM"))
		resp.ResponseWriteErr(w,encapsuledErr, http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	i.DomainControl.DeleteDomain(param.UUID, i.LibvirtInst)


	resp.ResponseWriteOK(w,domainInfo)
}

