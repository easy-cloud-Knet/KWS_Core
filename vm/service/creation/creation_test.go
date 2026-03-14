package creation

import (
	"errors"
	"fmt"
	"testing"

	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type mockConfigurer struct {
	xml []byte
	err error
}

func (m *mockConfigurer) GenerateXML(_ *zap.Logger) ([]byte, error) {
	return m.xml, m.err
}

type mockLibvirtConnect struct {
	domain *libvirt.Domain
	err    error
}

func (m *mockLibvirtConnect) DomainDefineXML(_ string) (*libvirt.Domain, error) {
	return m.domain, m.err
}

func TestCreateVM_ConfigurerError(t *testing.T) {
	creator := LocalCreatorFactory(
		&mockConfigurer{err: fmt.Errorf("config error")},
		&mockLibvirtConnect{},
		zap.NewNop(),
	)
	_, err := creator.CreateVM()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateVM_DomainDefineXMLError(t *testing.T) {
	creator := LocalCreatorFactory(
		&mockConfigurer{xml: []byte("<domain/>")},
		&mockLibvirtConnect{err: fmt.Errorf("libvirt error")},
		zap.NewNop(),
	)
	_, err := creator.CreateVM()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.DomainGenerationError) {
		t.Errorf("expected DomainGenerationError, got %v", err)
	}
}

// TODO: processCloudInitFiles, GenerateXML 내부 로직 테스트는
// YamlParsor, XMLDefiner 인터페이스 + VM_CREATE_XML.MarshalXML() 추가 후 가능.
