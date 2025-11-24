package api

import (
	"fmt"
	"net/http"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"github.com/easy-cloud-Knet/KWS_Core/vm/service/termination"
	"libvirt.org/go/libvirt"
)



func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	param:= &DeleteDomain{}
	resp:=ResponseGen[any]("domain number of"+param.UUID+", Force Shutdown VM")

	if err:=HttpDecoder(r,param); err!=nil{
		ERR:=virerr.ErrorJoin(err,fmt.Errorf("error shutting down vm, from forceShutdown vm "))
		resp.ResponseWriteErr(w,ERR, http.StatusInternalServerError)
		i.Logger.Error(ERR.Error())
		return
	}
	fmt.Println(resp)
	dom,err:= i.DomainControl.GetDomain(param.UUID,i.LibvirtInst)
	if err!= nil{
		ERR:=virerr.ErrorJoin(err,fmt.Errorf("error shutting down vm, retreving Get domin error "))
		resp.ResponseWriteErr(w,ERR, http.StatusInternalServerError)
		i.Logger.Error(ERR.Error())
		return
	}
	vcpu, err := dom.Domain.GetMaxVcpus()
	if err!= nil{
		ERR:=virerr.ErrorJoin(err,fmt.Errorf("error shutting down vm, retreving Get domin error "))
		resp.ResponseWriteErr(w,ERR, http.StatusInternalServerError)
		i.Logger.Error(ERR.Error())
		return
	}
	i.DomainControl.DomainListStatus.AddSleepingCPU(int(vcpu))

	DomainTerminator,_:= termination.DomainTerminatorFactory(dom)

	_,err=DomainTerminator.TerminateDomain()
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
		ERR:=virerr.ErrorJoin(err,fmt.Errorf("error deleting vm, unparsing HTTP request "))
		resp.ResponseWriteErr(w,ERR,http.StatusInternalServerError)
		i.Logger.Error(ERR.Error())
		return
	}
	if _, err := domCon.ReturnUUID(param.UUID); err!=nil{
		ERR:=virerr.ErrorJoin(err,fmt.Errorf("error deleting vm,	invalid UUID "))
		resp.ResponseWriteErr(w,ERR,http.StatusInternalServerError)
		i.Logger.Error(ERR.Error())
		return
	}
	// uuid 가 적합한지 확인

	
	domain, err := i.DomainControl.GetDomain(param.UUID,i.LibvirtInst)
	if err!=nil{
		ERR:=virerr.ErrorJoin(err,fmt.Errorf("error deleting vm, retreving Get domin error "))
		resp.ResponseWriteErr(w,ERR,http.StatusInternalServerError)
		i.Logger.Error(ERR.Error())
		return
		//error handling 
	}

	vcpu, err :=domain.Domain.GetMaxVcpus()
	if err != nil {
		ERR:=virerr.ErrorJoin(err,fmt.Errorf("error deleting vm, retreving Get domin error "))
		resp.ResponseWriteErr(w,ERR,http.StatusInternalServerError)
		i.Logger.Error(ERR.Error())
		return
	}// 삭제된 도메인에서는 vcpu count 를 가져올 수 없으므로 미리 가져옴 . 맘에 안듦. 나중에 수정할 예정


	DomainDeleter,_:=termination.DomainDeleterFactory(domain, param.DeletionType, param.UUID)
	domDeleted,err:=DomainDeleter.DeleteDomain()
	if err!=nil{
		ERR:=virerr.ErrorJoin(err,fmt.Errorf("error deleting vm, retreving Get domin error "))
		resp.ResponseWriteErr(w,ERR,http.StatusInternalServerError)
		i.Logger.Error(ERR.Error())
		return
	}
	i.DomainControl.DeleteDomain(domDeleted,param.UUID, int(vcpu))
	

	resp.ResponseWriteOK(w,nil)
}

