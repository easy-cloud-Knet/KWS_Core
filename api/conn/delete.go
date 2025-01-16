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



func DomainTerminatorFactory(DomainSeek DomainSeeker) (*DomainTerminator, error){
	return &DomainTerminator{
		DomainSeeker: DomainSeek,
	}, nil
}
func DomainDeleterFactory(DomainSeek DomainSeeker, DelType DomainDeleteType, uuid string)(*DomainDeleter, error){
	return &DomainDeleter{
		DomainSeeker: DomainSeek,
		DeletionType: DelType,
		DomainStatusManager:&DomainStatusManager{
			UUID:uuid,
		},}, nil
}

func (i *InstHandler)ForceShutDownVM(w http.ResponseWriter, r *http.Request){
	var param DeleteDomain
	if err:= json.NewDecoder(r.Body).Decode(&param); err!=nil{
		http.Error(w, "error decoding body", 1)
	}
	DomainSeeker:=&DomainSeekingByUUID{
		LibvirtInst: i.LibvirtInst,
		UUID: param.UUID,
		Domain: make([]*Domain, 0,1),
	}
	DomainTerminator,_:= DomainTerminatorFactory(DomainSeeker)
	domainInfo,err:=DomainTerminator.ShutDownDomain()
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

func (i *InstHandler)DeleteVM(w http.ResponseWriter, r *http.Request){
	var param DeleteDomain
	if err:= json.NewDecoder(r.Body).Decode(&param);err != nil{
		http.Error(w, "invalid parameter", http.StatusBadRequest)
	}
	DomainSeeker:= &DomainSeekingByUUID{
		LibvirtInst: i.LibvirtInst,
		UUID: param.UUID,
		Domain: make([]*Domain, 0,1),
	}
	DomainDeleter,_:=DomainDeleterFactory(DomainSeeker, param.DeletionType, param.UUID)
	
	domainInfo, err:=DomainDeleter.DeleteDomain()
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


func (DD *DomainDeleter) DeleteDomain() (*libvirt.DomainInfo, error){
	if err :=DD.DomainSeeker.SetDomain(); err!=nil{
		return &libvirt.DomainInfo{},err
	}
	dom,_ := DD.DomainSeeker.returnDomain()

	isRunning, _ :=dom[0].Domain.IsActive()
	if isRunning&& DD.DeletionType==SoftDelete {
		return &libvirt.DomainInfo{},fmt.Errorf("Domain Is Running %w", nil)
	}else if isRunning&&DD.DeletionType==HardDelete{
		domShut:=&DomainTerminator{DomainSeeker: DD.DomainSeeker}
		_,err := domShut.ShutDownDomain()
		if(err!=nil){
			return &libvirt.DomainInfo{},err
		}
	}
	basicFilePath:= "/var/lib/kws/"
	DomainPath,_:=ReturnUUID(DD.DomainStatusManager.UUID)
	FilePath := filepath.Join(basicFilePath,DomainPath.String())
	deleteCmd := exec.Command("rm", "-rf", FilePath)
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


func (DD *DomainTerminator) ShutDownDomain() (*libvirt.DomainInfo,error){
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



