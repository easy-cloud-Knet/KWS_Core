package status

import (
	"fmt"
	"testing"

	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type mockDomain struct {
	vcpus  uint
	vcpErr error
	xml    string
	xmlErr error
}

func (m *mockDomain) GetMaxVcpus() (uint, error)                          { return m.vcpus, m.vcpErr }
func (m *mockDomain) GetXMLDesc(_ libvirt.DomainXMLFlags) (string, error) { return m.xml, m.xmlErr }

var nopLogger = zap.NewNop()

const xmlTemplate = `<domain><vcpu>%d</vcpu><memory unit='KiB'>%d</memory></domain>`

func TestNew_Active(t *testing.T) {
	if _, ok := New(&mockDomain{}, true).(*LibvirtStatus); !ok {
		t.Error("expected *LibvirtStatus for active domain")
	}
}

func TestNew_Inactive(t *testing.T) {
	if _, ok := New(&mockDomain{}, false).(*XMLStatus); !ok {
		t.Error("expected *XMLStatus for inactive domain")
	}
}

func TestLibvirtStatus_CPU(t *testing.T) {
	dom := &mockDomain{vcpus: 8}
	result, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[vmtypes.SourceType]int{vmtypes.CPU: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[vmtypes.CPU] != 8 {
		t.Errorf("expected CPU=8, got %d", result[vmtypes.CPU])
	}
}

func TestLibvirtStatus_CPUError(t *testing.T) {
	dom := &mockDomain{vcpErr: fmt.Errorf("vcpu error")}
	_, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[vmtypes.SourceType]int{vmtypes.CPU: 0}, nopLogger)
	if err == nil {
		t.Error("expected error from GetMaxVcpus, got nil")
	}
}

func TestLibvirtStatus_UnknownSource(t *testing.T) {
	dom := &mockDomain{vcpus: 4}
	_, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[vmtypes.SourceType]int{"disk": 0}, nopLogger)
	if err == nil {
		t.Error("expected error for unknown source type, got nil")
	}
}

func TestXMLStatus_CPU(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 4, 1048576)}
	result, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[vmtypes.SourceType]int{vmtypes.CPU: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[vmtypes.CPU] != 4 {
		t.Errorf("expected CPU=4, got %d", result[vmtypes.CPU])
	}
}

func TestXMLStatus_XMLError(t *testing.T) {
	dom := &mockDomain{xmlErr: fmt.Errorf("xml error")}
	_, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[vmtypes.SourceType]int{vmtypes.CPU: 0}, nopLogger)
	if err == nil {
		t.Error("expected error from GetXMLDesc, got nil")
	}
}

func TestXMLStatus_Memory(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 2, 2097152)}
	result, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[vmtypes.SourceType]int{vmtypes.Memory: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[vmtypes.Memory] != 2097152 {
		t.Errorf("expected Memory=2097152, got %d", result[vmtypes.Memory])
	}
}

func TestXMLStatus_UnknownSource(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 4, 2097152)}
	_, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[vmtypes.SourceType]int{"disk": 0}, nopLogger)
	if err == nil {
		t.Error("expected error for unknown source type, got nil")
	}
}
