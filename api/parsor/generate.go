package parsor

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)
func (DP DomainParsor) Generate(LibvirtInst *libvirt.Connect, logger *zap.SugaredLogger) (*libvirt.Domain, error) {
	dirPath := fmt.Sprintf("/var/lib/kws/%s", DP.VMDescription.UUID)
	
	// 디렉토리 생성
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		errDesc := fmt.Errorf("failed to create directory (%s)", dirPath) 
		logger.Errorf("%v", errDesc)
		return nil, virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}
	logger.Debugf("Directory created: %s", dirPath)

	// cloud-init 파일 처리
	if err := DP.processCloudInitFiles(dirPath, logger); err != nil {
		return nil, err
	}

	// 디스크 이미지 생성
	if err := DP.CreateDiskImage(dirPath); err != nil {
		logger.Errorf("Disk image creation failed: %v", err)
		return nil, err
	}

	// ISO 파일 생성
	if err := DP.CreateISOFile(dirPath); err != nil {
		logger.Errorf("ISO file creation failed: %v", err)
		return nil, err
	}

	// XML 생성 및 도메인 등록
	DP.DeviceDefiner.XML_Parsor(DP.VMDescription)
	output, err := xml.MarshalIndent(*DP.DeviceDefiner, "", "  ")
	if err != nil {
		errDesc := fmt.Errorf("XML marshaling error: %w", err)
		logger.Errorf("%v", errDesc)
		return nil, virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	dom, err := conn.CreateDomainWithXML(LibvirtInst, output)
	if err != nil {
		errDesc := fmt.Errorf("libvirt domain creation error: %w", err)
		// logger.Errorf("%v", errDesc)
		return nil, virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	logger.Infof("Domain created successfully: UUID=%s, CPU=%d, RAM=%d, IP=%s",
		DP.DeviceDefiner.UUID, DP.VMDescription.HardwardInfo.CPU, DP.VMDescription.HardwardInfo.Memory, DP.VMDescription.NetConf.Ips[0])

	return dom, nil
}

func (DP DomainParsor) processCloudInitFiles(dirPath string, logger *zap.SugaredLogger) error {
	if err := DP.YamlParsorUser.ParseData(DP.VMDescription); err != nil {
		errDesc := fmt.Errorf("failed to parse cloud-init user-data: %w", err)
		logger.Errorf("%v", errDesc)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}
	
	if err := DP.YamlParsorUser.WriteFile(dirPath); err != nil {
		errDesc := fmt.Errorf("failed to write user-data file: %w", err)
		logger.Errorf("%v", errDesc)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	if err := DP.YamlParsorMeta.ParseData(DP.VMDescription); err != nil {
		errDesc := fmt.Errorf("failed to parse cloud-init meta-data: %w", err)
		logger.Errorf("%v", errDesc)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}
	if err := DP.YamlParsorMeta.WriteFile(dirPath); err != nil {
		errDesc := fmt.Errorf("failed to write meta-data file: %w", err)
		logger.Errorf("%v", errDesc)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	logger.Debugf("Cloud-init files processed successfully in %s", dirPath)
	return nil
}

 
func (DP DomainParsor) CreateDiskImage(dirPath string) error{
	baseImage := fmt.Sprintf("/var/lib/kws/baseimg/%s",  DP.VMDescription.OS )
	targetImage := filepath.Join(dirPath, fmt.Sprintf("%s.qcow2",  DP.VMDescription.UUID))
	qemuImgCmd := exec.Command("qemu-img", "create",
		"-b", baseImage,
		"-f", "qcow2",
		"-F", "qcow2",
		targetImage, "10G",
	)
	if err := qemuImgCmd.Run(); err != nil {
		errorDescription:=fmt.Errorf("generating Disk image error, may have duplicdated uuid or lack of HD capacity %s", dirPath)
		// disk 크기 확인 후 부족할 시 LackCapacityHD 반환
		 return virerr.ErrorGen(virerr.DomainGenerationError, errorDescription)
	}


	return nil
}



func (DP DomainParsor) CreateISOFile(dirPath string)error{

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
		errorDescription:=fmt.Errorf("generating ISO image error, may have duplicdated uuid or wrong format of yaml file %s", dirPath)
		return virerr.ErrorGen(virerr.DomainGenerationError, errorDescription)
	}
	return nil
}



func ParsorFactoryFromRequest(param *VM_Init_Info, logger *zap.SugaredLogger)*DomainParsor{
	return &DomainParsor{
		YamlParsorUser: &User_data_yaml{},
		YamlParsorMeta: &Meta_data_yaml{},
		VMDescription:  param,
		DeviceDefiner: &VM_CREATE_XML{},
	}

}