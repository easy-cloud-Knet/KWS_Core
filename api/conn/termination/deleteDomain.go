package termination

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	domCon "github.com/easy-cloud-Knet/KWS_Core.git/api/conn/DomCon"
	"libvirt.org/go/libvirt"
)


func DomainDeleterFactory(Domain *domCon.Domain, DelType DomainDeleteType, uuid string) (*DomainDeleter, error) {
	return &DomainDeleter{
		uuid:uuid,
		domain: Domain,
		DeletionType: DelType,
		}, nil
}

func (DD *DomainDeleter) Operation() (*libvirt.Domain,error){
	dom := DD.domain

	isRunning, _ := dom.Domain.IsActive()
	if isRunning && DD.DeletionType == SoftDelete {
		
		return  nil,fmt.Errorf("domain Is Running %w, cannot softDelete running Domain", nil)
	} else if isRunning && DD.DeletionType == HardDelete {
		DomainTerminator := &DomainTerminator{domain: dom}
		_,err := DomainTerminator.Operation()
		if err != nil {
			return nil,fmt.Errorf("%w, failed deleting Domain in libvirt Instance, ", err)
		}
	}
	basicFilePath := "/var/lib/kws/"
	fmt.Println("domain uuid %s", DD.uuid)
	FilePath := filepath.Join(basicFilePath, DD.uuid)
	deleteCmd := exec.Command("rm", "-rf", FilePath)
	deleteCmd.Stdout = os.Stdout
	deleteCmd.Stderr = os.Stderr

	if err := deleteCmd.Run(); err != nil {
		return nil,fmt.Errorf("%w failed deleteing files in %s", err, FilePath)
	}

	if err := dom.Domain.Undefine(); err != nil {
		return nil,fmt.Errorf("%w, failed deleting Domain in libvirt Instance, ", err)
	}

	return dom.Domain,nil
}


