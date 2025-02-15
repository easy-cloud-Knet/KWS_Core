package conn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
		return nil, fmt.Errorf("Domain Is Running %w, cannot softDelete running Domain", nil)
	} else if isRunning && DD.DeletionType == HardDelete {
		domShut := &DomainTerminator{domain: dom}
		_, err := domShut.ShutDownDomain()
		if err != nil {
			return nil, fmt.Errorf("%w, failed deleting Domain in libvirt Instance, ", err)
		}
	}
	basicFilePath := "/var/lib/kws/"
 
	FilePath := filepath.Join(basicFilePath, uuid)
	deleteCmd := exec.Command("rm", "-rf", FilePath)
	deleteCmd.Stdout = os.Stdout
	deleteCmd.Stderr = os.Stderr

	if err := deleteCmd.Run(); err != nil {
		return nil, fmt.Errorf("%w failed deleteing files in %s", err, FilePath)
	}

	// domainInfo, err := dom.Domain.GetInfo()
	// if err != nil {
	// 	return nil, fmt.Errorf("%w Error Retreving DomInfo in After Deleting Domain", err)
	// } // 다른 정보 추가 고려

	if err := dom.Domain.Undefine(); err != nil {
		return nil, fmt.Errorf("%w, failed deleting Domain in libvirt Instance, ", err)
	}

	return nil, nil
}

func (DD *DomainTerminator) ShutDownDomain() (*libvirt.DomainInfo, error) {
 
	dom:= DD.domain

	isRunning, _ := dom.Domain.IsActive()
	if !isRunning {
		return nil, fmt.Errorf("requested Domain to shutdown is already Dead ")
	}

	if err := dom.Domain.Destroy(); err != nil {
		fmt.Println("error occured while deleting Domain")
		return nil, fmt.Errorf("internal Error in Libvirt occured while shutting down domain")
	}
	// domainInfo, err := dom.Domain.GetInfo()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	return nil, nil
}
