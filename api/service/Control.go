package service

import (
	"fmt"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	domCon "github.com/easy-cloud-Knet/KWS_Core.git/api/conn/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn/termination"
	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
	"libvirt.org/go/libvirt"
)



func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	param:= &DeleteDomain{}
	resp:=ResponseGen[any]("domain number of"+param.UUID+", Force Shutdown VM")

	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,virerr.ErrorJoin(err,fmt.Errorf("error shutting down vm, from forceShutdown vm ")), http.StatusInternalServerError)
		return
	}
	fmt.Println(resp)
	dom,err:= i.DomainControl.GetDomain(param.UUID,i.LibvirtInst)
	if err!= nil{
		resp.ResponseWriteErr(w,virerr.ErrorJoin(err,fmt.Errorf("error shutting down vm, retreving Get domin error ")), http.StatusInternalServerError)
		return
	}

	DomainTerminator,_:= termination.DomainTerminatorFactory(dom)

	_,err=DomainTerminator.Operation()
	if err!= nil{
		resp.ResponseWriteErr(w,virerr.ErrorJoin(err,fmt.Errorf("error shutting down vm, retreving Get domin error ")), http.StatusInternalServerError)
		return
	}

	
	resp.ResponseWriteOK(w,nil)
	}




func (i *InstHandler)DeleteVM(w http.ResponseWriter, r *http.Request){
	param:=&DeleteDomain{}
	resp:= ResponseGen[libvirt.DomainInfo]("Deleting Vm")

	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,virerr.ErrorJoin(err,fmt.Errorf("error deleting vm, retreving Get domin error ")),http.StatusInternalServerError)
		return
	}
	if _, err := domCon.ReturnUUID(param.UUID); err!=nil{
		resp.ResponseWriteErr(w,virerr.ErrorJoin(err,fmt.Errorf("error deleting vm, retreving Get domin error ")),http.StatusInternalServerError)
		return
	}
	// uuid 가 적합한지 확인

	
	domain, err := i.DomainControl.GetDomain(param.UUID,i.LibvirtInst)
	if err!=nil{
		resp.ResponseWriteErr(w,virerr.ErrorJoin(err,fmt.Errorf("error deleting vm, retreving Get domin error ")),http.StatusInternalServerError)
		return
		//error handling 
	}

	DomainDeleter,_:=termination.DomainDeleterFactory(domain, param.DeletionType, param.UUID)
	_,err=DomainDeleter.Operation()
	if err!=nil{
		resp.ResponseWriteErr(w,virerr.ErrorJoin(err,fmt.Errorf("error deleting vm, retreving Get domin error ")),http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	DomCon:=conn.DomainControllerInjection(i.DomainControl,DomainDeleter)

	DomCon.DomainDeleteWithOperation(i.Logger,param.UUID)
	

	resp.ResponseWriteOK(w,nil)
}

