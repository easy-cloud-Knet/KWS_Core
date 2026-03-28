package fetcher

import (
	"fmt"
	"testing"

	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
	"libvirt.org/go/libvirt"
)

type mockDomain struct {
	xml string
	err error
}

func (m *mockDomain) GetXMLDesc(_ libvirt.DomainXMLFlags) (string, error) {
	return m.xml, m.err
}

const xmlTemplate = `<domain><vcpu>%d</vcpu><memory unit='KiB'>%d</memory></domain>`

func TestFetch_CPU(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 4, 2097152)}
	result, err := NewXMLFetcher().Fetch(dom, map[vmtypes.SourceType]int{vmtypes.CPU: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[vmtypes.CPU] != 4 {
		t.Errorf("expected CPU=4, got %d", result[vmtypes.CPU])
	}
}

func TestFetch_Memory(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 2, 2097152)}
	result, err := NewXMLFetcher().Fetch(dom, map[vmtypes.SourceType]int{vmtypes.Memory: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[vmtypes.Memory] != 2097152 {
		t.Errorf("expected Memory=2097152, got %d", result[vmtypes.Memory])
	}
}

func TestFetch_UnknownSource(t *testing.T) {
	dom := &mockDomain{xml: fmt.Sprintf(xmlTemplate, 4, 2097152)}
	_, err := NewXMLFetcher().Fetch(dom, map[vmtypes.SourceType]int{"disk": 0})
	if err == nil {
		t.Error("expected error for unknown source type, got nil")
	}
}

func TestFetch_XMLParseError(t *testing.T) {
	dom := &mockDomain{err: fmt.Errorf("libvirt error")}
	_, err := NewXMLFetcher().Fetch(dom, map[vmtypes.SourceType]int{vmtypes.CPU: 0})
	if err == nil {
		t.Error("expected error from GetXMLDesc, got nil")
	}
}
