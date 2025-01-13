package conn

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"libvirt.org/go/libvirt"
)


func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	var param DeleteDomain
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		http.Error(w, "error decoding body", 1)
	}
	DomainTurner := &DomainController{
		DomainSeeker: &DomainSeekingByUUID{
			LibvirtInst: i.LibvirtInst,
			UUID: param.UUID,
			Domain: make([]*Domain, 0,1),
		},
	}
	domainInfo,err:=DomainTurner.ShutDownDomain()
	if err!= nil{
		http.Error(w,"error while shutting down Domain", http.StatusBadRequest)
		fmt.Println(err)
	}
	//uuid unparsing 중 에러, destroyDom 에서 에러
	// 수신시 에러 발생 가능 ,추후 에러 핸들링 


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"message":   fmt.Sprintf("VM with UUID %s Shutdown successfully.", param.UUID),
		"domainInfo": domainInfo,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

// func (i *InstHandler)StartVM(w http.ResponseWriter, r *http.Request){
// 	var param StartDomain
// 	if err:= json.NewDecoder(r.Body).Decode(&param);err != nil{
// 		http.Error(w, "invalid parameter", http.StatusBadRequest)
// 	}

// 	DomainStarer := &DomainController{
// 		DomainSeeker: &DomainSeekingByUUID{
// 			LibvirtInst: i.LibvirtInst,
// 			UUID: param.UUID,
// 			Domain: make([]*Domain, 0,1),
// 		},
// 	}
// 	dom,err:=DomainStarer.StartDomain()
// }
func (DD *DomainController)StartDomain()(*libvirt.DomainInfo,error){
	if err :=DD.DomainSeeker.SetDomain(); err!=nil{
		return &libvirt.DomainInfo{},err
	}
	dom,_ := DD.DomainSeeker.returnDomain()
	isRunning, _ :=dom[0].Domain.IsActive()
	if !isRunning {
		return &libvirt.DomainInfo{},fmt.Errorf("Domain Is already running %w", nil)
	}

	if err:= dom[0].Domain.Destroy(); err!=nil{
		fmt.Println("error occured while deleting Domain")
		return &libvirt.DomainInfo{},err
	}
	defer dom[0].Domain.Free()
	domainInfo,err:= dom[0].Domain.GetInfo()
	if err!=nil{
		fmt.Println(err)	
		//error handler needed
		return &libvirt.DomainInfo{}, err
	}
	return domainInfo,nil
}

func (i *InstHandler)DeleteVM(w http.ResponseWriter, r *http.Request){
	var param DeleteDomain
	if err:= json.NewDecoder(r.Body).Decode(&param);err != nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}

	DomainDeleter := &DomainController{
		DomainSeeker: &DomainSeekingByUUID{
			LibvirtInst: i.LibvirtInst,
			UUID: param.UUID,
			Domain: make([]*Domain, 0,1),
		},
	}
	domainInfo, err:=DomainDeleter.DeleteDomain(param.DeletionType)
	if err!=nil{
		fmt.Println(err)
		http.Error(w,"error while destroying vm", http.StatusInternalServerError)
	}
		//uuid unparsing 중 에러, undefine,destroyDom 에서 에러, 켜져 있지만 softdelete가 
	// 수신시 에러 발생 가능 ,추후 에러 핸들링 

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"message":   fmt.Sprintf("VM with UUID %s Deletion successfully.", param.UUID),
		"domainInfo": domainInfo,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
	
}


func (DD *DomainController) DeleteDomain(deleteType DomainDeleteType) (*libvirt.DomainInfo, error){
	if err :=DD.DomainSeeker.SetDomain(); err!=nil{
		return &libvirt.DomainInfo{},err
	}
	dom,_ := DD.DomainSeeker.returnDomain()

	isRunning, _ :=dom[0].Domain.IsActive()
	if isRunning&& deleteType==SoftDelete {
		return &libvirt.DomainInfo{},fmt.Errorf("Domain Is Running %w", nil)
	}else if isRunning&&deleteType==HardDelete{
		_,err := DD.ShutDownDomain()
		if(err!=nil){
			return &libvirt.DomainInfo{},err
		}
	}


	basicFilePath:= "/var/lib/kws/"
	DomainPath,_ :=DD.DomainSeeker.ReturnUUID()
	DD.DomainStatusManager.FilePath = filepath.Join(basicFilePath,DomainPath.String())

	deleteCmd := exec.Command("rm", "-rf", DD.DomainStatusManager.FilePath)
	deleteCmd.Stdout = os.Stdout
	deleteCmd.Stderr = os.Stderr


	if err := deleteCmd.Run(); err !=nil{
		log.Printf("qemu-img command failed: %v", err)
		return &libvirt.DomainInfo{},nil
	}
	domainInfo,err:= dom[0].Domain.GetInfo()
	if err:= dom[0].Domain.Undefine(); err!=nil{
		fmt.Println("error occured while deleting Domain")
		return &libvirt.DomainInfo{},nil
	}
	defer dom[0].Domain.Free()
	
	
	if err!=nil{
		fmt.Println(err)	
		//error handler needed
		return &libvirt.DomainInfo{}, nil
	}
	return domainInfo,nil
}


func (DD *DomainController) ShutDownDomain() (*libvirt.DomainInfo,error){
	if err :=DD.DomainSeeker.SetDomain(); err!=nil{
		return &libvirt.DomainInfo{},err
	}
	dom,_ := DD.DomainSeeker.returnDomain()
	isRunning, _ :=dom[0].Domain.IsActive()
	if !isRunning {
		return &libvirt.DomainInfo{},fmt.Errorf("Domain Is Not Running %w", nil)
	}

	if err:= dom[0].Domain.Destroy(); err!=nil{
		fmt.Println("error occured while deleting Domain")
		return &libvirt.DomainInfo{},err
	}
	defer dom[0].Domain.Free()
	domainInfo,err:= dom[0].Domain.GetInfo()
	if err!=nil{
		fmt.Println(err)	
		//error handler needed
		return &libvirt.DomainInfo{}, err
	}
	return domainInfo,nil
}



