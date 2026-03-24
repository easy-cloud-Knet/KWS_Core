package creation

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core/internal/config"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	uuid "github.com/easy-cloud-Knet/KWS_Core/pkg/UUID"
	"github.com/easy-cloud-Knet/KWS_Core/pkg/parsor"
	safepath "github.com/easy-cloud-Knet/KWS_Core/pkg/safePath"
	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	userconfig "github.com/easy-cloud-Knet/KWS_Core/pkg/yaml/cloud-init"
	"go.uber.org/zap"
)

func LocalConfFactory(param *vmtypes.VM_Init_Info, logger *zap.Logger) *localConfigurer {
	return &localConfigurer{
		VMDescription:  param,
		YamlParsorUser: &userconfig.User_data_yaml{},
		YamlParsorMeta: &userconfig.Meta_data_yaml{},
		DeviceDefiner:  &parsor.VM_CREATE_XML{},
	}

}
func LocalCreatorFactory(confige Configurer, libvirtInst LibvirtConnect, logger *zap.Logger) *LocalCreator {
	return &LocalCreator{
		DomainConfiger: confige,
		libvirtInst:    libvirtInst,
		logger:         logger,
	}
}

func (DCB *LocalCreator) CreateVM() (*domCon.Domain, error) {
	output, err := DCB.DomainConfiger.GenerateXML(DCB.logger)
	if err != nil {
		DCB.logger.Error("error while generating VM config", zap.Error(err))
		return nil, err
	}

	domain, err := DCB.libvirtInst.DomainDefineXML(string(output))
	if err != nil {
		errDesc := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("in domain-Creator, error defining domain via libvirt: %w", err))
		DCB.logger.Error(errDesc.Error())
		return nil, errDesc
	}

	if err := domain.Create(); err != nil {
		errDesc := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("in domain-Creator, error starting domain: %w", err))
		DCB.logger.Error(errDesc.Error())
		return nil, errDesc
	}

	return domCon.NewDomainInstance(domain), nil
}

func (DB localConfigurer) GenerateXML(logger *zap.Logger) ([]byte, error) {
	if err := DB.Generate(logger); err != nil {
		return nil, err
	}
	output, err := xml.MarshalIndent(*DB.DeviceDefiner, "", "  ")
	if err != nil {
		return nil, virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("XML marshaling error: %w", err))
	}
	return output, nil
}

func (DB localConfigurer) Generate(logger *zap.Logger) error {
	if _, err := uuid.ValidateAndReturnUUID(DB.VMDescription.UUID); err != nil {
		logger.Error("invalid UUID provided", zap.String("uuid", DB.VMDescription.UUID), zap.Error(err))
		return virerr.ErrorGen(virerr.InvalidUUID, err)
	}

	dirPath, err := safepath.GetSafeFilePath(config.StorageBase, DB.VMDescription.UUID)
	if dirPath == "" {
		errDesc := fmt.Errorf("failed to generate safe file path for UUID %s %v", DB.VMDescription.UUID, err)
		logger.Error("failed to generate safe file path or some macilous attack happened. aborting", zap.Error(errDesc))
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		errDesc := fmt.Errorf("failed to create directory (%s)", dirPath)
		logger.Error("failed making directory", zap.Error(errDesc))
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	// cloud-init 파일 처리
	if err := DB.processCloudInitFiles(dirPath); err != nil {
		errorEncapsed := virerr.ErrorJoin(err, fmt.Errorf("in domain-parsor,"))
		logger.Error(errorEncapsed.Error())
		return errorEncapsed
	}
	logger.Info("generating configuration file successfully done", zap.String("filePath", dirPath))

	if err := DB.CreateDiskImage(dirPath, DB.VMDescription.HardwardInfo.Disk); err != nil {
		errorEncapsed := virerr.ErrorJoin(err, fmt.Errorf("in domain-parsor,"))
		logger.Error(errorEncapsed.Error())
		return errorEncapsed
	}

	// ISO 파일 생성
	if err := DB.CreateISOFile(dirPath); err != nil {
		errorEncapsed := virerr.ErrorJoin(err, fmt.Errorf("in domain-parsor,"))
		logger.Error(errorEncapsed.Error())
		return err
	}

	if err := DB.DeviceDefiner.XML_Parsor(DB.VMDescription); err != nil {
		return virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("XML_Parsor error: %w", err))
	}
	return nil
}

func (DB localConfigurer) processCloudInitFiles(dirPath string) error {
	if err := DB.YamlParsorUser.ParseData(DB.VMDescription); err != nil {
		errDesc := fmt.Errorf("failed to parse cloud-init user-data: %w", err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	if err := DB.YamlParsorUser.WriteFile(dirPath); err != nil {
		errDesc := fmt.Errorf("failed to write user-data file: %w", err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	if err := DB.YamlParsorMeta.ParseData(DB.VMDescription); err != nil {
		errDesc := fmt.Errorf("failed to parse cloud-init meta-data: %w", err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}
	if err := DB.YamlParsorMeta.WriteFile(dirPath); err != nil {
		errDesc := fmt.Errorf("failed to write meta-data file: %w", err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
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

// local 파일에서 vm을 생성할 경우 사용
