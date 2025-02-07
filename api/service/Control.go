package service

import (
	"fmt"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)



func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	
	param:= &DeleteDomain{}
	resp:=ResponseGen("Force Shutdown VM")
	if err:=HttpDecoder(w,r,param); err!=nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error booting vm",err), http.StatusInternalServerError)
		return
	}
	
	DomainSeeker := conn.DomSeekUUIDFactory(i.LibvirtInst, param.UUID)

	DomainTerminator,_:= conn.DomainTerminatorFactory(DomainSeeker)

	domainInfo,err:=DomainTerminator.ShutDownDomain()
	if err!= nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error booting vm",err), http.StatusInternalServerError)
		return
	}
	//uuid unparsing 중 에러, destroyDom 에서 에러
	// 수신시 에러 발생 가능 ,추후 에러 핸들링 
	
 
	resp.ResponseWriteOK(w,domainInfo)
	}




func (i *InstHandler)DeleteVM(w http.ResponseWriter, r *http.Request){
	param:=&DeleteDomain{}
	resp:= ResponseGen("Deleting Vm")
	if err:=HttpDecoder(w,r,param); err!=nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error booting vm",err), http.StatusInternalServerError)
		return
	}
 
	DomainSeeker:=conn.DomSeekUUIDFactory(i.LibvirtInst, param.UUID)
	DomainDeleter,_:=conn.DomainDeleterFactory(DomainSeeker, param.DeletionType, param.UUID)
	
	domainInfo, err:=DomainDeleter.DeleteDomain()
	if err!=nil{
		resp.ResponseWriteErr(w,fmt.Errorf("%w error booting vm",err), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
		//uuid unparsing 중 에러, undefine,destroyDom 에서 에러, 켜져 있지만 softdelete가 
	// 수신시 에러 발생 가능 ,추후 에러 핸들링 

	resp.ResponseWriteOK(w,domainInfo)
}

