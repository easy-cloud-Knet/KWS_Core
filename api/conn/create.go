package conn

import (
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
		
		return err
	}
	return nil
}


func (DGL DomainGeneratorLocal) CloudInitConf(param *parsor.VM_Init_Info)error{
	DGL.DataParsor.YamlParsor.Parse_data(param)
	
	dirPath := fmt.Sprintf("/var/lib/kws/%s", DGL.DomainStatusManager.UUID)
	if err:= DGL.DataParsor.YamlParsor.FileConfig(dirPath); err!=nil{
		return err
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

	log.Println("Creating disk image...")
	if err := qemuImgCmd.Run(); err != nil {
		log.Printf("qemu-img command failed: %v", err)
		return err
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
		return err
	}
	return nil
}