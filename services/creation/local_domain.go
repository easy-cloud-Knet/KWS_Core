package creation

import (
	"encoding/xml"
	"fmt"
	"os"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"github.com/easy-cloud-Knet/KWS_Core/internal/config"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	uuid "github.com/easy-cloud-Knet/KWS_Core/pkg/UUID"
	"github.com/easy-cloud-Knet/KWS_Core/pkg/parsor"
	rollback "github.com/easy-cloud-Knet/KWS_Core/pkg/rollBack"
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
	output, dirPath, err := DCB.DomainConfiger.GenerateXML(DCB.logger)
	if err != nil {
		DCB.logger.Error("error while generating VM config", zap.Error(err))
		return nil, err
	}

	rbm := &rollback.RollBackManager{}
	rbm.Add(func() error { return os.RemoveAll(dirPath) })

	domain, err := DCB.libvirtInst.DomainDefineXML(string(output))
	if err != nil {
		if rbErr := rbm.Execute(); rbErr != nil {
			DCB.logger.Error("rollback failed after DomainDefineXML error", zap.Error(rbErr))
		}
		errDesc := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("in domain-Creator, error defining domain via libvirt: %w", err))
		DCB.logger.Error(errDesc.Error())
		return nil, errDesc
	}
	rbm.Add(func() error { return domain.Undefine() })

	if err := domain.Create(); err != nil {
		if rbErr := rbm.Execute(); rbErr != nil {
			DCB.logger.Error("rollback failed after domain.Create error", zap.Error(rbErr))
		}
		errDesc := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("in domain-Creator, error starting domain: %w", err))
		DCB.logger.Error(errDesc.Error())
		return nil, errDesc
	}

	return domCon.NewDomainInstance(domain), nil
}

func (DB localConfigurer) GenerateXML(logger *zap.Logger) ([]byte, string, error) {
	dirPath, err := safepath.GetSafeFilePath(config.StorageBase, DB.VMDescription.UUID)
	if err != nil {
		return nil, "", virerr.ErrorGen(virerr.DomainGenerationError, err)
	}
	if err := DB.Generate(logger); err != nil {
		return nil, "", err
	}
	output, err := xml.MarshalIndent(*DB.DeviceDefiner, "", "  ")
	if err != nil {
		os.RemoveAll(dirPath)
		return nil, "", virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("XML marshaling error: %w", err))
	}
	return output, dirPath, nil
}

func (DB localConfigurer) Generate(logger *zap.Logger) error {
	if _, err := uuid.ValidateAndReturnUUID(DB.VMDescription.UUID); err != nil {
		logger.Error("invalid UUID provided", zap.String("uuid", DB.VMDescription.UUID), zap.Error(err))
		return virerr.ErrorGen(virerr.InvalidUUID, err)
	}

	dirPath, err := safepath.GetSafeFilePath(config.StorageBase, DB.VMDescription.UUID)
	if err != nil {
		logger.Error("failed to generate safe file path", zap.String("uuid", DB.VMDescription.UUID), zap.Error(err))
		return virerr.ErrorGen(virerr.DomainGenerationError, err)
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		errDesc := fmt.Errorf("failed to create directory (%s)", dirPath)
		logger.Error("failed making directory", zap.Error(errDesc))
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	cleanup := func(err error) error {
		os.RemoveAll(dirPath)
		return err
	}

	// cloud-init 파일 처리
	if err := DB.processCloudInitFiles(dirPath); err != nil {
		errorEncapsed := virerr.ErrorJoin(err, fmt.Errorf("in domain-parsor,"))
		logger.Error(errorEncapsed.Error())
		return cleanup(errorEncapsed)
	}
	logger.Info("generating configuration file successfully done", zap.String("filePath", dirPath))

	if err := DB.CreateDiskImage(dirPath, DB.VMDescription.HardwardInfo.Disk); err != nil {
		errorEncapsed := virerr.ErrorJoin(err, fmt.Errorf("in domain-parsor,"))
		logger.Error(errorEncapsed.Error())
		return cleanup(errorEncapsed)
	}

	// ISO 파일 생성
	if err := DB.CreateISOFile(dirPath); err != nil {
		errorEncapsed := virerr.ErrorJoin(err, fmt.Errorf("in domain-parsor,"))
		logger.Error(errorEncapsed.Error())
		return cleanup(errorEncapsed)
	}

	if err := DB.DeviceDefiner.XML_Parsor(DB.VMDescription); err != nil {
		return cleanup(virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("XML_Parsor error: %w", err)))
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
