package conn

import (
	"fmt"
	"log"
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


func (DD *DomainDeleter) DeleteDomain() (*libvirt.DomainInfo, error){
	if err :=DD.DomainSeeker.SetDomain(); err!=nil{
		return &libvirt.DomainInfo{},err
	}
	dom,_ := DD.DomainSeeker.ReturnDomain()

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
	dom,_ := DD.DomainSeeker.ReturnDomain()
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



