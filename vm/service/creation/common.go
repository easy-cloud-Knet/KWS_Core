package creation

import (
	"fmt"
	"os/exec"
	"path/filepath"

	virerr "github.com/easy-cloud-Knet/KWS_Core.git/error"
)


func (DB localConfigurer)CreateDiskImage(dirPath string) error {
	baseImage := fmt.Sprintf("/var/lib/kws/baseimg/%s", DB.VMDescription.OS)
	targetImage := filepath.Join(dirPath, fmt.Sprintf("%s.qcow2", DB.VMDescription.UUID))
	qemuImgCmd := exec.Command("qemu-img", "create",
		"-b", baseImage,
		"-f", "qcow2",
		"-F", "qcow2",
		targetImage, "10G",
	)
	if err := qemuImgCmd.Run(); err != nil {
		errorDescription := fmt.Errorf("generating Disk image error, may have duplicdated uuid or lack of HD capacity %s, %v", dirPath, err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errorDescription)
	}

	return nil
}
