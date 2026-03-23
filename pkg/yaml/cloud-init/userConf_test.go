package userconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	vmtypes "github.com/easy-cloud-Knet/KWS_Core/pkg/types"
)

func TestConfigNetworkIP(t *testing.T) {
	tests := []struct {
		name        string
		ips         []string
		wantLen     int
		wantPath    string
		wantAddr    string
		wantGateway string
	}{
		{
			name:        "single IP",
			ips:         []string{"192.168.1.100"},
			wantLen:     1,
			wantPath:    "/etc/systemd/network/10-enp0s3.network",
			wantAddr:    "192.168.1.100",
			wantGateway: "192.168.1.1",
		},
		{
			name:    "multiple IPs",
			ips:     []string{"10.0.0.5", "172.16.0.10"},
			wantLen: 2,
		},
		{
			name:    "empty",
			ips:     []string{},
			wantLen: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := &User_data_yaml{}
			result := u.configNetworkIP(tc.ips)

			if len(result) != tc.wantLen {
				t.Fatalf("expected %d entries, got %d", tc.wantLen, len(result))
			}
			if tc.wantPath != "" && result[0].Path != tc.wantPath {
				t.Errorf("path: expected %s, got %s", tc.wantPath, result[0].Path)
			}
			if tc.wantAddr != "" && !strings.Contains(result[0].Content, tc.wantAddr) {
				t.Errorf("content missing address %s", tc.wantAddr)
			}
			if tc.wantGateway != "" && !strings.Contains(result[0].Content, tc.wantGateway) {
				t.Errorf("content missing gateway %s", tc.wantGateway)
			}
		})
	}
}

func TestFetchos(t *testing.T) {
	u := &User_data_yaml{}
	result := u.fetchos()

	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Path != "/etc/profile.d/99-my-motd.sh" {
		t.Errorf("unexpected path: %s", result[0].Path)
	}
	if result[0].Content == "" {
		t.Error("fetchos content is empty")
	}
}

func TestUserDataWriteFile(t *testing.T) {
	dir := t.TempDir()
	u := &User_data_yaml{PackageUpdatable: true}

	if err := u.WriteFile(dir); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "user-data"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), "#cloud-config\n") {
		t.Error("missing #cloud-config header")
	}
}

func TestMetaDataWriteFile(t *testing.T) {
	dir := t.TempDir()
	m := &Meta_data_yaml{
		Instance_ID:   "test-uuid",
		Local_Host_Id: "test-host",
	}

	if err := m.WriteFile(dir); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "meta-data"))
	if err != nil {
		t.Fatalf("meta-data file not created: %v", err)
	}
	if !strings.Contains(string(content), "test-uuid") {
		t.Error("meta-data missing instance-id")
	}
}

func TestMetaDataParseData(t *testing.T) {
	m := &Meta_data_yaml{}
	param := &vmtypes.VM_Init_Info{
		UUID:    "abc-123",
		DomName: "test-vm",
	}

	if err := m.ParseData(param); err != nil {
		t.Fatal(err)
	}
	if m.Instance_ID != "abc-123" {
		t.Errorf("expected Instance_ID abc-123, got %s", m.Instance_ID)
	}
	if m.Local_Host_Id != "test-vm" {
		t.Errorf("expected Local_Host_Id test-vm, got %s", m.Local_Host_Id)
	}
}
