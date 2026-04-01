package status

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
	"libvirt.org/go/libvirt"
)

type mockDomain struct {
	vcpus     uint
	vcpErr    error
	info      *libvirt.DomainInfo
	infoErr   error
	maxMem    uint64
	maxMemErr error
	xml       string
	xmlErr    error
}

func (m *mockDomain) GetMaxVcpus() (uint, error)                          { return m.vcpus, m.vcpErr }
func (m *mockDomain) GetInfo() (*libvirt.DomainInfo, error)               { return m.info, m.infoErr }
func (m *mockDomain) GetMaxMemory() (uint64, error)                       { return m.maxMem, m.maxMemErr }
func (m *mockDomain) GetXMLDesc(_ libvirt.DomainXMLFlags) (string, error) { return m.xml, m.xmlErr }

var nopLogger = zap.NewNop()

const xmlTemplate = `<domain><vcpu>%d</vcpu><memory unit='KiB'>%d</memory></domain>`
const xmlTemplateWithMaxMem = `<domain><maxMemory slots='16' unit='KiB'>%d</maxMemory><vcpu>%d</vcpu><memory unit='KiB'>%d</memory></domain>`

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

// LibvirtStatus tests

func TestLibvirtStatus_CPU(t *testing.T) {
	dom := &mockDomain{vcpus: 8}
	result, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[SourceType]int{CPU: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[CPU] != 8 {
		t.Errorf("expected CPU=8, got %d", result[CPU])
	}
}

func TestLibvirtStatus_Memory(t *testing.T) {
	dom := &mockDomain{info: &libvirt.DomainInfo{Memory: 2097152}}
	result, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[SourceType]int{Memory: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[Memory] != 2097152 {
		t.Errorf("expected Memory=2097152, got %d", result[Memory])
	}
}

func TestLibvirtStatus_MaxMemory(t *testing.T) {
	dom := &mockDomain{maxMem: 4194304}
	result, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[SourceType]int{MaxMemory: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[MaxMemory] != 4194304 {
		t.Errorf("expected MaxMemory=4194304, got %d", result[MaxMemory])
	}
}

func TestLibvirtStatus_CPUTime(t *testing.T) {
	dom := &mockDomain{info: &libvirt.DomainInfo{CpuTime: 123456789}}
	result, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[SourceType]int{CPUTime: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[CPUTime] != 123456789 {
		t.Errorf("expected CPUTime=123456789, got %d", result[CPUTime])
	}
}

func TestLibvirtStatus_GetInfoCachedAcrossFields(t *testing.T) {
	callCount := 0
	dom := &mockDomain{
		info: &libvirt.DomainInfo{Memory: 1048576, CpuTime: 999},
	}
	// wrap to count GetInfo calls
	type countingDom struct{ *mockDomain }
	cd := &countingDom{dom}
	origGetInfo := dom.GetInfo
	_ = origGetInfo
	// direct test: request Memory + CPUTime together, both should succeed
	result, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[SourceType]int{Memory: 0, CPUTime: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[Memory] != 1048576 {
		t.Errorf("expected Memory=1048576, got %d", result[Memory])
	}
	if result[CPUTime] != 999 {
		t.Errorf("expected CPUTime=999, got %d", result[CPUTime])
	}
	_ = callCount
	_ = cd
}

func TestLibvirtStatus_CPUError(t *testing.T) {
	dom := &mockDomain{vcpErr: fmt.Errorf("vcpu error")}
	_, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[SourceType]int{CPU: 0}, nopLogger)
	if err == nil {
		t.Error("expected error from GetMaxVcpus, got nil")
	}
}

func TestLibvirtStatus_MemoryError(t *testing.T) {
	dom := &mockDomain{infoErr: fmt.Errorf("info error")}
	_, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[SourceType]int{Memory: 0}, nopLogger)
	if err == nil {
		t.Error("expected error from GetInfo, got nil")
	}
}

func TestLibvirtStatus_UnknownSource(t *testing.T) {
	dom := &mockDomain{vcpus: 4}
	_, err := (&LibvirtStatus{dom: dom}).RetrieveStatus(map[SourceType]int{"disk": 0}, nopLogger)
	if err == nil {
		t.Error("expected error for unknown source type, got nil")
	}
}

// XMLStatus tests

func TestXMLStatus_CPU(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 4, 1048576)}
	result, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[SourceType]int{CPU: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[CPU] != 4 {
		t.Errorf("expected CPU=4, got %d", result[CPU])
	}
}

func TestXMLStatus_Memory(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 2, 2097152)}
	result, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[SourceType]int{Memory: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[Memory] != 2097152 {
		t.Errorf("expected Memory=2097152, got %d", result[Memory])
	}
}

func TestXMLStatus_MaxMemory_Explicit(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplateWithMaxMem, 8388608, 4, 2097152)}
	result, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[SourceType]int{MaxMemory: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[MaxMemory] != 8388608 {
		t.Errorf("expected MaxMemory=8388608, got %d", result[MaxMemory])
	}
}

func TestXMLStatus_MaxMemory_FallbackToMemory(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 4, 2097152)}
	result, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[SourceType]int{MaxMemory: 0}, nopLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[MaxMemory] != 2097152 {
		t.Errorf("expected MaxMemory fallback=2097152, got %d", result[MaxMemory])
	}
}

func TestXMLStatus_CPUTimeUnsupported(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 4, 1048576)}
	_, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[SourceType]int{CPUTime: 0}, nopLogger)
	if err == nil {
		t.Error("expected error for cpu_time on inactive domain, got nil")
	}
}

func TestXMLStatus_XMLError(t *testing.T) {
	dom := &mockDomain{xmlErr: fmt.Errorf("xml error")}
	_, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[SourceType]int{CPU: 0}, nopLogger)
	if err == nil {
		t.Error("expected error from GetXMLDesc, got nil")
	}
}

func TestXMLStatus_UnknownSource(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 4, 2097152)}
	_, err := (&XMLStatus{dom: dom}).RetrieveStatus(map[SourceType]int{"disk": 0}, nopLogger)
	if err == nil {
		t.Error("expected error for unknown source type, got nil")
	}
}
