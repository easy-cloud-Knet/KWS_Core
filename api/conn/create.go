package conn

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
)

func (DGL DomainGeneratorLocal) CreateFolder()error{

	dirPath := fmt.Sprintf("/var/lib/kws/%s", DGL.DomainStatusManager.UUID)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		
		return ErrorGen(DomainGenerationError, errors.New("making directory error generating domain, may have duplicated uuid"))
	}
	return nil
}


func (DGL DomainGeneratorLocal) CloudInitConf(param *parsor.VM_Init_Info)error{
	err :=DGL.DataParsor.YamlParsor.Parse_data(param)
	if err!= nil{
		return err
	}
	
	dirPath := fmt.Sprintf("/var/lib/kws/%s", DGL.DomainStatusManager.UUID)
	if err:= DGL.DataParsor.YamlParsor.FileConfig(dirPath); err!=nil{
		errorDescription:=fmt.Errorf("writing Cloud Init Config File Failed, may Have Duplicated File in %s, %w" ,dirPath,err)
		return ErrorGen(DomainGenerationError, errorDescription)
	}

	return nil
}

func (DGL DomainGeneratorLocal) CreateDiskImage() error{
	dirPath := fmt.Sprintf("/var/lib/kws/%s", DGL.DomainStatusManager.UUID)
	baseImage := fmt.Sprintf("/var/lib/kws/baseimg/%s", DGL.OS )
	targetImage := filepath.Join(dirPath, fmt.Sprintf("%s.qcow2", DGL.DomainStatusManager.UUID))
	qemuImgCmd := exec.Command("qemu-img", "create",
		"-b", baseImage,
		"-f", "qcow2",
		"-F", "qcow2",
		targetImage, "10G",
	)

	qemuImgCmd.Stdout = os.Stdout
	qemuImgCmd.Stderr = os.Stderr

	if err := qemuImgCmd.Run(); err != nil {
		errorDescription:=fmt.Errorf("generating Disk image error, may have duplicdated uuid or lack of HD capacity %s, %w", dirPath,err)
		// disk 크기 확인 후 부족할 시 LackCapacityHD 반환
		return ErrorGen(DomainGenerationError, errorDescription)
	}
	return nil
}

func (DGL DomainGeneratorLocal) CreateISOFile()error{
	dirPath := fmt.Sprintf("/var/lib/kws/%s", DGL.DomainStatusManager.UUID)

	isoOutput := filepath.Join(dirPath, "cidata.iso")
	userDataPath := filepath.Join(dirPath, "user-data")
	metaDataPath := filepath.Join(dirPath, "meta-data")

	genisoCmd := exec.Command("genisoimage",
		"--output", isoOutput,
		"-V", "cidata",
		"-r", "-J",
		userDataPath, metaDataPath,
	)

	genisoCmd.Stdout = os.Stdout
	genisoCmd.Stderr = os.Stderr

	log.Println("Generating ISO image...")
	if err := genisoCmd.Run(); err != nil {
		log.Printf("genisoimage command failed: %v", err)
		errorDescription:=fmt.Errorf("generating ISO image error, may have duplicdated uuid or wrong format of yaml file %s, %w", dirPath,err)
		return ErrorGen(DomainGenerationError, errorDescription)
	}
	return nil
}