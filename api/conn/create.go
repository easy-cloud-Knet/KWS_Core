package conn

import (
	"fmt"

	virerr "github.com/easy-cloud-Knet/KWS_Core.git/api/error"
	"libvirt.org/go/libvirt"
)


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