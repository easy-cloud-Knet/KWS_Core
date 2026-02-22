package creation

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"github.com/easy-cloud-Knet/KWS_Core/vm/parsor"
	userconfig "github.com/easy-cloud-Knet/KWS_Core/vm/parsor/cloud-init"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)


func LocalConfFactory(param *parsor.VM_Init_Info, logger *zap.Logger) *localConfigurer {
	return &localConfigurer{
		VMDescription:  param,
		YamlParsorUser: &userconfig.User_data_yaml{},
		YamlParsorMeta: &userconfig.Meta_data_yaml{},
		DeviceDefiner:  &parsor.VM_CREATE_XML{},
	}

}
func LocalCreatorFactory(confige *localConfigurer,libvirtInst *libvirt.Connect,logger *zap.Logger )(*LocalCreator){
	return &LocalCreator{
		DomainConfiger:confige,
		libvirtInst:libvirtInst,
		logger:logger,
	}
}




func (DCB *LocalCreator) CreateVM()(*domCon.Domain,error){
	//DomainConfiger 를 인터페이스로 둔다면 타입 체크로 분기 가능 
	err:= DCB.DomainConfiger.Generate(DCB.libvirtInst,DCB.logger)
	if err!=nil{
		DCB.logger.Warn("error whiling configuring base Configures, ", zap.Error(err))
	}

	output, err := xml.MarshalIndent(*DCB.DomainConfiger.DeviceDefiner, "", "  ")
	if err != nil {
		errDesc := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("in domain-Creator, XML marshaling error: %w", err))
		DCB.logger.Error(errDesc.Error())
		return nil,errDesc
	}

	domain,err :=CreateDomainWithXML(DCB.libvirtInst,output)
	fmt.Println(string(output))
	if err!=nil{
		errDesc := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("in domain-Creator, error occured from creating with libvirt: %w", err))
		DCB.logger.Error(errDesc.Error())
		return nil,errDesc
		
	}
	err=domain.Create()
	if err!=nil{
		errDesc := virerr.ErrorGen(virerr.DomainGenerationError, fmt.Errorf("in domain-Creator, XML marshaling error: %w", err))
		DCB.logger.Error(errDesc.Error())
		return nil,errDesc
	}
	domconDom:=domCon.NewDomainInstance(domain)
	return domconDom,nil
}




func (DB localConfigurer) Generate(LibvirtInst *libvirt.Connect, logger *zap.Logger) (error) {
	dirPath,err := parsor.GetSafeFilePath("/var/lib/kws", DB.VMDescription.UUID)
	if dirPath == "" {
		errDesc := fmt.Errorf("failed to generate safe file path for UUID %s %v", DB.VMDescription.UUID, err)
		logger.Error("failed to generate safe file path or some macilous attack happened. aborting", zap.Error(errDesc))
		return virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	
	

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		errDesc := fmt.Errorf("failed to create directory (%s)", dirPath)
		logger.Error("failed making directory", zap.Error(errDesc))
		return  virerr.ErrorGen(virerr.DomainGenerationError, errDesc)
	}

	// cloud-init 파일 처리
	if err := DB.processCloudInitFiles(dirPath); err != nil {
		errorEncapsed := virerr.ErrorJoin(err, fmt.Errorf("in domain-parsor,"))
		logger.Error(errorEncapsed.Error())
		return  errorEncapsed
	}
	logger.Info("generating configuration file successfully done", zap.String("filePath", dirPath))

	if err := DB.CreateDiskImage(dirPath); err != nil {
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

	DB.DeviceDefiner.XML_Parsor(DB.VMDescription)
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


func CreateDomainWithXML(LibvirtInst *libvirt.Connect ,config []byte) (*libvirt.Domain, error) {

	// DomainCreateXMLWithFiles를 호출하여 도메인을 생성합니다.
	domain, err := LibvirtInst.DomainDefineXML(string(config))
	if err != nil {
		return nil, virerr.ErrorGen(virerr.DomainGenerationError,fmt.Errorf("domain creating with libvirt daemon from xml err %w", err))
		// cpu나 ip 중복 등을 검사하는 코드를 삽입하고, 그에 맞는 에러 반환 필요
	} 
	//이전까지 생성 된 파일 삭제 해야됨.
  return domain ,nil
}
// local 파일에서 vm을 생성할 경우 사용