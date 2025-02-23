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
func (DP DomainParsor) Generate(LibvirtInst *libvirt.Connect, logger *zap.Logger) (*libvirt.Domain, error) {
	dirPath := fmt.Sprintf("/var/lib/kws/%s", DP.VMDescription.UUID)
	
	// 디렉토리 생성
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		errDesc := fmt.Errorf("failed to create directory (%s)", dirPath) 
		logger.Error("failed making directory", zap.Error(errDesc))
		return nil, virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	// cloud-init 파일 처리
	if err := DP.processCloudInitFiles(dirPath); err != nil {
		errorEncapsed := virerr.ErrorJoin(err,fmt.Errorf("in domain-parsor,"))
		logger.Error(errorEncapsed.Error())
		return nil, errorEncapsed
	}
		logger.Info("generating configuration file successfully done", zap.String("filePath", dirPath))

	if err := DP.CreateDiskImage(dirPath); err != nil {
		errorEncapsed := virerr.ErrorJoin(err,fmt.Errorf("in domain-parsor,"))
		logger.Error(errorEncapsed.Error())		
		return nil, errorEncapsed
	}

	// ISO 파일 생성
	if err := DP.CreateISOFile(dirPath); err != nil {
		errorEncapsed := virerr.ErrorJoin(err,fmt.Errorf("in domain-parsor,"))
		logger.Error(errorEncapsed.Error())		
		return nil, err
	}

	// XML 생성 및 도메인 등록
	DP.DeviceDefiner.XML_Parsor(DP.VMDescription)
	output, err := xml.MarshalIndent(*DP.DeviceDefiner, "", "  ")
	if err != nil {
		errDesc := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("in domain-parsor, XML marshaling error: %w", err))
		logger.Error(errDesc.Error())
		return nil, errDesc
	}

	dom, err := conn.CreateDomainWithXML(LibvirtInst, output)
	if err != nil {
		errDesc := virerr.ErrorJoin(err, fmt.Errorf("in domain-parsor"))
		logger.Error(errDesc.Error())
		return nil, errDesc
	}

	logger.Info("Domain created successfully: ", 
	zap.String("UUID",DP.DeviceDefiner.UUID ),
 	zap.Int("CPU",DP.VMDescription.HardwardInfo.CPU),
	zap.Int("Memory",DP.VMDescription.HardwardInfo.Memory),
	zap.String("IP", DP.VMDescription.NetConf.Ips[0]),
)
	return dom, nil
}

func (DP DomainParsor) processCloudInitFiles(dirPath string) error {
	if err := DP.YamlParsorUser.ParseData(DP.VMDescription); err != nil {
		errDesc := fmt.Errorf("failed to parse cloud-init user-data: %w", err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}
	
	if err := DP.YamlParsorUser.WriteFile(dirPath); err != nil {
		errDesc := fmt.Errorf("failed to write user-data file: %w", err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	if err := DP.YamlParsorMeta.ParseData(DP.VMDescription); err != nil {
		errDesc := fmt.Errorf("failed to parse cloud-init meta-data: %w", err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}
	if err := DP.YamlParsorMeta.WriteFile(dirPath); err != nil {
		errDesc := fmt.Errorf("failed to write meta-data file: %w", err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

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
		errorDescription:=fmt.Errorf("generating Disk image error, may have duplicdated uuid or lack of HD capacity %s, %v", dirPath, err)
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
		errorDescription:=fmt.Errorf("generating ISO image error, may have duplicdated uuid or wrong format of yaml file %s, %v", dirPath, err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errorDescription)
	}
	return nil
}



func ParsorFactoryFromRequest(param *VM_Init_Info, logger *zap.Logger)*DomainParsor{
	return &DomainParsor{
		YamlParsorUser: &User_data_yaml{},
		YamlParsorMeta: &Meta_data_yaml{},
		VMDescription:  param,
		DeviceDefiner: &VM_CREATE_XML{},
	}

}