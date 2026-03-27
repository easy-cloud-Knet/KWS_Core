package virtxml

import (
	"fmt"
	"testing"
)

func TestNew_ReturnsNonNil(t *testing.T) {
	d := New()
	if d == nil {
		t.Fatal("New() returned nil")
	}
}

func TestConvertExistingDomain_Success(t *testing.T) {
	xmlStr := `<domain type="kvm"><name>test-vm</name><uuid>123e4567-e89b-12d3-a456-426614174000</uuid></domain>`

	domain, err := ConvertExistingDomain(0, func(_ DomainXMLFlags) (string, error) {
		return xmlStr, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if domain.Name != "test-vm" {
		t.Errorf("name: expected test-vm, got %s", domain.Name)
	}
	if domain.UUID != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("uuid: expected 123e4567-e89b-12d3-a456-426614174000, got %s", domain.UUID)
	}
}

func TestConvertExistingDomain_GetXMLDescError(t *testing.T) {
	domain, err := ConvertExistingDomain(0, func(_ DomainXMLFlags) (string, error) {
		return "", fmt.Errorf("libvirt error")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if domain != nil {
		t.Fatal("expected nil domain on error")
	}
}

func TestConvertExistingDomain_InvalidXML(t *testing.T) {
	domain, err := ConvertExistingDomain(0, func(_ DomainXMLFlags) (string, error) {
		return "not valid xml <<<", nil
	})
	if err == nil {
		t.Fatal("expected error for invalid XML")
	}
	if domain != nil {
		t.Fatal("expected nil domain on parse error")
	}
}

func TestConvertExistingDomain_FlagPassthrough(t *testing.T) {
	var received DomainXMLFlags
	_, err := ConvertExistingDomain(42, func(f DomainXMLFlags) (string, error) {
		received = f
		return `<domain><name>x</name></domain>`, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if received != 42 {
		t.Errorf("flag not passed through: expected 42, got %d", received)
	}
}
