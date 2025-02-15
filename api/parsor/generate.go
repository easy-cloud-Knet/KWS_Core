package parsor

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
	"libvirt.org/go/libvirt"
)


func (DP DomainParsor)Generate(LibvirtInst *libvirt.Connect) (*libvirt.Domain,error){
	dirPath := fmt.Sprintf("/var/lib/kws/%s", DP.VMDescription.UUID)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil,virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("making directory error generating domain, may have duplicated uuid %w",err))
	}

	if err:= DP.YamlParsorUser.ParseData(DP.VMDescription); err!=nil {	
		return nil,virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("making directory error generating domain, may have duplicated uuid %w", err))
	}

	if err:= DP.YamlParsorUser.WriteFile(dirPath); err!=nil {	
		return nil,virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("making directory error generating domain, may have duplicated uuid %w",err))
	}

	if err:= DP.YamlParsorMeta.WriteFile(dirPath);err != nil {	
		return nil,virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("making directory error generating domain, may have duplicated uuid %w",err))
	}
	
	if err:= DP.CreateDiskImage();err!=nil{

	}

	if err:= DP.CreateISOFile();err!=nil{

	}

	DP.DeviceDefiner.XML_Parsor(DP.VMDescription)

	output, err := xml.MarshalIndent(DP.DeviceDefiner, "", "  ")
	if err!=nil{
		virerr.ErrorGen(virerr.DomainGenerationError,fmt.Errorf("XML deparsing Error, whilie defining hardware spec in Generation, %w", err))
	}

	dom,err := conn.CreateDomainWithXML(LibvirtInst, output)
	if err!=nil{
		
	}

	return dom,nil

}

func (DP DomainParsor) CreateDiskImage() error{
	dirPath := fmt.Sprintf("/var/lib/kws/%s", DP.VMDescription.UUID)
	baseImage := fmt.Sprintf("/var/lib/kws/baseimg/%s",  DP.VMDescription.OS )
	targetImage := filepath.Join(dirPath, fmt.Sprintf("%s.qcow2",  DP.VMDescription.UUID))
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
		return virerr.ErrorGen(virerr.DomainGenerationError, errorDescription)
	}
	return nil
}



func (DP DomainParsor) CreateISOFile()error{
	dirPath := fmt.Sprintf("/var/lib/kws/%s", DP.VMDescription.UUID)

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
		return virerr.ErrorGen(virerr.DomainGenerationError, errorDescription)
	}
	return nil
}



func ParsorFactoryFromRequest(param *VM_Init_Info)*DomainParsor{
	return &DomainParsor{
		YamlParsorUser: &User_data_yaml{},
		YamlParsorMeta: &Meta_data_yaml{},
		VMDescription:  param,
		DeviceDefiner: &VM_CREATE_XML{},
	}

}