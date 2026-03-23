package termination

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/easy-cloud-Knet/KWS_Core/config"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
)

func DomainDeleterFactory(d Domain, delType DomainDeleteType, uuid string) DomainDeletion {
	return &DomainDeleter{
		uuid:         uuid,
		domain:       d,
		DeletionType: delType,
	}
}

func (DD *DomainDeleter) DeleteDomain() error {
	isRunning, _ := DD.domain.IsActive()
	if isRunning && DD.DeletionType == SoftDelete {
		return virerr.ErrorGen(virerr.DeletionDomainError, fmt.Errorf("domain is running, cannot softDelete running domain"))
	} else if isRunning && DD.DeletionType == HardDelete {
		t := &DomainTerminator{domain: DD.domain}
		if err := t.TerminateDomain(); err != nil {
			return virerr.ErrorGen(virerr.DeletionDomainError, fmt.Errorf("failed deleting domain in libvirt instance: %w", err))
		}
	}

	FilePath := filepath.Join(config.StorageBase, DD.uuid)
	deleteCmd := exec.Command("rm", "-rf", FilePath)
	deleteCmd.Stdout = os.Stdout
	deleteCmd.Stderr = os.Stderr

	if err := deleteCmd.Run(); err != nil {
		return virerr.ErrorGen(virerr.DeletionDomainError, fmt.Errorf("failed deleting files in %s: %w", FilePath, err))
	}

	if err := DD.domain.Undefine(); err != nil {
		return virerr.ErrorGen(virerr.DeletionDomainError, fmt.Errorf("failed deleting domain in libvirt instance: %w", err))
	}

	return nil
}
