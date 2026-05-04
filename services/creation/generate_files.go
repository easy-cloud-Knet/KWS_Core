package creation

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/easy-cloud-Knet/KWS_Core/internal/config"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func (DB localConfigurer) CreateDiskImage(dirPath string, diskSize int) error {
	baseImage := fmt.Sprintf("%s/baseimg/%s", config.StorageBase, DB.VMDescription.OS)
	targetImage := filepath.Join(dirPath, fmt.Sprintf("%s.qcow2", DB.VMDescription.UUID))
	qemuImgCmd := exec.Command("qemu-img", "create",
		"-b", baseImage,
		"-f", "qcow2",
		"-F", "qcow2",
		targetImage, fmt.Sprintf("%dG", diskSize), // 10G
	)
	if err := qemuImgCmd.Run(); err != nil {
		errorDescription := fmt.Errorf("generating Disk image error, check duplicdated uuid or lack of HD capacity, or validity for base img %s, %v", dirPath, err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errorDescription)
	}

	return nil
}

func (DB localConfigurer) CreateISOFile(dirPath string) error {

	isoOutput := filepath.Join(dirPath, "cidata.iso")
	userDataPath := filepath.Join(dirPath, "user-data")
	metaDataPath := filepath.Join(dirPath, "meta-data")

	genisoCmd := exec.Command("genisoimage",
		"--output", isoOutput,
		"-V", "cidata",
		"-r", "-J",
		userDataPath, metaDataPath,
	)

	if err := genisoCmd.Run(); err != nil {
		errorDescription := fmt.Errorf("generating ISO image error, may have duplicdated uuid or wrong format of yaml file %s, %v", dirPath, err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errorDescription)
	}
	return nil
}

