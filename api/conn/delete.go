package conn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"libvirt.org/go/libvirt"
)

func DomainTerminatorFactory(DomainSeek DomainSeeker) (*DomainTerminator, error) {
	return &DomainTerminator{
		DomainSeeker: DomainSeek,
	}, nil
}
func DomainDeleterFactory(DomainSeek DomainSeeker, DelType DomainDeleteType, uuid string) (*DomainDeleter, error) {
	return &DomainDeleter{
		DomainSeeker: DomainSeek,
		DeletionType: DelType,
		DomainStatusManager: &DomainStatusManager{
			UUID: uuid,
		}}, nil
}

func (DD *DomainDeleter) DeleteDomain() (*libvirt.DomainInfo, error) {
	dom, err := DD.DomainSeeker.ReturnDomain()
	if err != nil {
		return nil, err
	}

	isRunning, _ := dom[0].Domain.IsActive()
	if isRunning && DD.DeletionType == SoftDelete {
		return nil, fmt.Errorf("Domain Is Running %w, cannot softDelete running Domain", nil)
	} else if isRunning && DD.DeletionType == HardDelete {
		domShut := &DomainTerminator{DomainSeeker: DD.DomainSeeker}
		_, err := domShut.ShutDownDomain()
		if err != nil {
			return nil, fmt.Errorf("%w, failed deleting Domain in libvirt Instance, ", err)
		}
	}
	basicFilePath := "/var/lib/kws/"
	DomainPath, err := ReturnUUID(DD.DomainStatusManager.UUID)
	if err != nil {
		return nil, fmt.Errorf("error Occured while decoding UUID in DeleteDom, %w", err)
	}
	FilePath := filepath.Join(basicFilePath, DomainPath.String())
	deleteCmd := exec.Command("rm", "-rf", FilePath)
	deleteCmd.Stdout = os.Stdout
	deleteCmd.Stderr = os.Stderr

	if err := deleteCmd.Run(); err != nil {
		return &libvirt.DomainInfo{}, fmt.Errorf("%w failed deleteing files in %s", err, FilePath)
	}

	domainInfo, err := dom[0].Domain.GetInfo()
	if err != nil {
		return &libvirt.DomainInfo{}, fmt.Errorf("%w Error Retreving DomInfo in After Deleting Domain", err)
	} // 다른 정보 추가 고려

	if err := dom[0].Domain.Undefine(); err != nil {
		return &libvirt.DomainInfo{}, fmt.Errorf("%w, failed deleting Domain in libvirt Instance, ", err)
	}
	defer dom[0].Domain.Free()

	return domainInfo, nil
}

func (DD *DomainTerminator) ShutDownDomain() (*libvirt.DomainInfo, error) {
 
	dom, err := DD.DomainSeeker.ReturnDomain()
	if err!=nil{
		return nil, fmt.Errorf("%w, failed deleting Domain in libvirt Instance, ", err)

	}
	isRunning, _ := dom[0].Domain.IsActive()
	if !isRunning {
		return nil, fmt.Errorf("requested Domain to shutdown is already Dead ")
	}

	if err := dom[0].Domain.Destroy(); err != nil {
		fmt.Println("error occured while deleting Domain")
		return nil, fmt.Errorf("internal Error in Libvirt occured while shutting down domain")
	}
	defer dom[0].Domain.Free()

	domainInfo, err := dom[0].Domain.GetInfo()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return domainInfo, nil
}
