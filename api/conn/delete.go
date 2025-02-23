package conn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
	"libvirt.org/go/libvirt"
)

func DomainTerminatorFactory(Domain *Domain) (*DomainTerminator, error) {
	return &DomainTerminator{
		domain: Domain,
	}, nil
}
func DomainDeleterFactory(Domain *Domain, DelType DomainDeleteType) (*DomainDeleter, error) {
	return &DomainDeleter{
		domain: Domain,
		DeletionType: DelType,
		}, nil
}

func (DD *DomainDeleter) DeleteDomain(uuid string) (*libvirt.DomainInfo, error) {
	dom := DD.domain

	isRunning, _ := dom.Domain.IsActive()
	if isRunning && DD.DeletionType == SoftDelete {
		return nil, virerr.ErrorGen(virerr.DeletionDomainError,fmt.Errorf("Domain Is Running,cannot softDelete running Domain")) 
	} else if isRunning && DD.DeletionType == HardDelete {
		domShut := &DomainTerminator{domain: dom}

		_, err := domShut.ShutDownDomain()
		if err != nil {
			return nil, virerr.ErrorGen(virerr.DeletionDomainError,fmt.Errorf("failed shutting domain down while hard deleting domain, %w",err)) 
		}
	}
	basicFilePath := "/var/lib/kws/"
 
	FilePath := filepath.Join(basicFilePath, uuid)
	deleteCmd := exec.Command("rm", "-rf", FilePath)
	deleteCmd.Stdout = os.Stdout
	deleteCmd.Stderr = os.Stderr

	if err := deleteCmd.Run(); err != nil {
		return nil, virerr.ErrorGen(virerr.DeletionDomainError,fmt.Errorf("failed deleting exsiting local file deleting domain, %w ",err)) 
	}

	if err := dom.Domain.Undefine(); err != nil {
		return nil, virerr.ErrorGen(virerr.DeletionDomainError,fmt.Errorf("failed deleting domain from libvirt daemon, %w ",err)) 
	}

	return nil, nil
}

func (DD *DomainTerminator) ShutDownDomain() (*libvirt.DomainInfo, error) {
 
	dom:= DD.domain

	isRunning, _ := dom.Domain.IsActive()
	if !isRunning {
		return nil, virerr.ErrorGen(virerr.DomainShutdownError, fmt.Errorf("domain is still running, can not soft shut running domain"))
	}

	if err := dom.Domain.Destroy(); err != nil {
		return nil, virerr.ErrorGen(virerr.DomainShutdownError, fmt.Errorf("domain shutting down error from libvirt daemon %w",err))
	}

	return nil, nil
}
