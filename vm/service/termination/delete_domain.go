package termination

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"libvirt.org/go/libvirt"
)


func DomainDeleterFactory(Domain *domCon.Domain, DelType DomainDeleteType, uuid string) (*DomainDeleter, error) {
	return &DomainDeleter{
		uuid:uuid,
		domain: Domain,
		DeletionType: DelType,
		}, nil
}

func (DD *DomainDeleter) DeleteDomain() (*libvirt.Domain,error){
	dom := DD.domain

	isRunning, _ := dom.Domain.IsActive()
	if isRunning && DD.DeletionType == SoftDelete {

		return  nil,virerr.ErrorGen(virerr.DeletionDomainError, fmt.Errorf("domain is running, cannot softDelete running domain"))
	} else if isRunning && DD.DeletionType == HardDelete {
		DomainTerminator := &DomainTerminator{domain: dom}
		_,err := DomainTerminator.TerminateDomain()
		if err != nil {
			return nil,virerr.ErrorGen(virerr.DeletionDomainError, fmt.Errorf("failed deleting domain in libvirt instance: %w", err))
		}
	}
	basicFilePath := "/var/lib/kws/"
	FilePath := filepath.Join(basicFilePath, DD.uuid)
	deleteCmd := exec.Command("rm", "-rf", FilePath)
	deleteCmd.Stdout = os.Stdout
	deleteCmd.Stderr = os.Stderr

	if err := deleteCmd.Run(); err != nil {
		return nil,virerr.ErrorGen(virerr.DeletionDomainError, fmt.Errorf("failed deleting files in %s: %w", FilePath, err))
	}

	if err := dom.Domain.Undefine(); err != nil {
		return nil,virerr.ErrorGen(virerr.DeletionDomainError, fmt.Errorf("failed deleting domain in libvirt instance: %w", err))
	}


	return dom.Domain,nil
}
