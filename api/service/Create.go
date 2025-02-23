package service

import (
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)


func (i *InstHandler) CreateVMLocal(w http.ResponseWriter, r *http.Request) {

	resp:=ResponseGen[libvirt.DomainInfo]("CreateVm")
	param:=&parsor.VM_Init_Info{}
	
	if err:=HttpDecoder(r,param); err!=nil{
		resp.ResponseWriteErr(w,err, http.StatusBadRequest)
		i.Logger.Error("error occured while decoding user's parameter of requested creation")
		return
	}
	i.Logger.Info("Handling Create VM", zap.String("uuid", param.UUID))

	DomainParsor:= parsor.ParsorFactoryFromRequest(param, i.Logger)
	//추후에 요청별로 다른 방식으로 vm을 생성할 수 있음 
	// 예: snapshot 기반으로 생성, 사용자가 입력한 cloud-init 폼 등등
	DomainGenerator := &conn.DomainGenerator{
		DataParsor: DomainParsor,
	}
	// domain parsor 의 인터페이스를 사용하기 떄문에, Parsor Factory에서 반환만 잘하면 재사용 가능

	dom,err := DomainGenerator.DataParsor.Generate(i.LibvirtInst, i.Logger); 
	if err!=nil{
		i.Logger.Error("domain Generating failed, while encoding user's config", zap.String("uuid", param.UUID))
		resp.ResponseWriteErr(w,err, http.StatusInternalServerError)
		return
	}

	err = dom.Create()
	if err!= nil{
		i.Logger.Error("domain Generating failed, while booting created virtual machine", zap.String("uuid", param.UUID))
		resp.ResponseWriteErr(w,err, http.StatusInternalServerError)
		return 
	}
	newDomain := conn.NewDomainInstance(dom)

	i.DomainControl.AddNewDomain(newDomain,param.UUID)

	resp.ResponseWriteOK(w,nil)
}


