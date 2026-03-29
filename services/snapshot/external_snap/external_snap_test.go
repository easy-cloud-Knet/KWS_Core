package external

import (
	"errors"
	"fmt"
	"testing"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

type mockExternalSnapshotDomain struct {
	isActive        bool
	isActiveErr     error
	uuid            string
	uuidErr         error
	xmlDesc         string
	xmlDescErr      error
	createErr       error
	createHandle    externalSnapshotHandle
	listErr         error
	snapshots       []externalSnapshotHandle
	blockJobInfo    externalBlockJobInfo
	blockJobInfoErr error

	lastCreateOpts externalSnapshotCreateExecOptions
}

func (m *mockExternalSnapshotDomain) IsActive() (bool, error) {
	return m.isActive, m.isActiveErr
}

func (m *mockExternalSnapshotDomain) CreateExternalSnapshot(_ string, opts externalSnapshotCreateExecOptions) (externalSnapshotHandle, error) {
	m.lastCreateOpts = opts
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createHandle != nil {
		return m.createHandle, nil
	}
	return &mockExternalSnapshotHandle{name: "created-snap"}, nil
}

func (m *mockExternalSnapshotDomain) ListAllSnapshots() ([]externalSnapshotHandle, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.snapshots, nil
}

func (m *mockExternalSnapshotDomain) StartBlockCommit(_, _, _ string) error {
	return nil
}

func (m *mockExternalSnapshotDomain) BlockJobInfo(_ string) (externalBlockJobInfo, error) {
	if m.blockJobInfoErr != nil {
		return externalBlockJobInfo{}, m.blockJobInfoErr
	}
	if m.blockJobInfo.End == 0 {
		return externalBlockJobInfo{Cur: 1, End: 1}, nil
	}
	return m.blockJobInfo, nil
}

func (m *mockExternalSnapshotDomain) AbortBlockJobPivot(_ string) error {
	return nil
}

func (m *mockExternalSnapshotDomain) UpdateDeviceConfig(_ string) error {
	return nil
}

func (m *mockExternalSnapshotDomain) UUIDString() (string, error) {
	if m.uuidErr != nil {
		return "", m.uuidErr
	}
	if m.uuid != "" {
		return m.uuid, nil
	}
	return "test-uuid", nil
}

func (m *mockExternalSnapshotDomain) XMLDesc() (string, error) {
	if m.xmlDescErr != nil {
		return "", m.xmlDescErr
	}
	return m.xmlDesc, nil
}

type mockExternalSnapshotHandle struct {
	name       string
	nameErr    error
	xmlDesc    string
	xmlDescErr error
	deleteErr  error
}

func (m *mockExternalSnapshotHandle) Name() (string, error) {
	if m.nameErr != nil {
		return "", m.nameErr
	}
	return m.name, nil
}

func (m *mockExternalSnapshotHandle) XMLDesc() (string, error) {
	if m.xmlDescErr != nil {
		return "", m.xmlDescErr
	}
	return m.xmlDesc, nil
}

func (m *mockExternalSnapshotHandle) Delete() error {
	return m.deleteErr
}

func (m *mockExternalSnapshotHandle) Free() error {
	return nil
}

func TestCreateExternalSnapshot_InvalidName(t *testing.T) {
	domain := &mockExternalSnapshotDomain{}
	_, err := createExternalSnapshot(domain, "", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.InvalidParameter) {
		t.Fatalf("expected InvalidParameter, got %v", err)
	}
}

func TestCreateExternalSnapshot_MapsCreateOptions(t *testing.T) {
	tmp := t.TempDir()
	domain := &mockExternalSnapshotDomain{
		isActive: true,
		xmlDesc: `<domain><devices>
			<disk device='disk' type='file'><driver type='qcow2'/><source file='/vm/vda.qcow2'/><target dev='vda' bus='virtio'/></disk>
			<disk device='disk' type='file'><driver type='qcow2'/><source file='/vm/vdb.qcow2'/><target dev='vdb' bus='virtio'/></disk>
		</devices></domain>`,
		createHandle: &mockExternalSnapshotHandle{name: "snap-a"},
	}

	name, err := createExternalSnapshot(domain, "snap-a", &ExternalSnapshotOptions{BaseDir: tmp, Live: true, Quiesce: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "snap-a" {
		t.Fatalf("expected snapshot name snap-a, got %s", name)
	}
	if !domain.lastCreateOpts.Live || !domain.lastCreateOpts.Quiesce || !domain.lastCreateOpts.Atomic {
		t.Fatalf("expected live/quiesce/atomic all true, got %+v", domain.lastCreateOpts)
	}
}

func TestCreateExternalSnapshot_XMLDescError(t *testing.T) {
	domain := &mockExternalSnapshotDomain{xmlDescErr: fmt.Errorf("xml fail")}
	_, err := createExternalSnapshot(domain, "snap", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

func TestDeleteExternalSnapshot_NotFound(t *testing.T) {
	domain := &mockExternalSnapshotDomain{
		snapshots: []externalSnapshotHandle{},
	}

	err := deleteExternalSnapshot(domain, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

func TestListExternalSnapshots_FiltersExternalOnly(t *testing.T) {
	domain := &mockExternalSnapshotDomain{
		snapshots: []externalSnapshotHandle{
			&mockExternalSnapshotHandle{
				name:    "ext-1",
				xmlDesc: `<domainsnapshot><name>ext-1</name><disks><disk name='vda' snapshot='external'><source file='/snap/vda.qcow2'/></disk></disks></domainsnapshot>`,
			},
			&mockExternalSnapshotHandle{
				name:    "int-1",
				xmlDesc: `<domainsnapshot><name>int-1</name><disks><disk name='vda' snapshot='internal'/></disks></domainsnapshot>`,
			},
		},
	}

	names, err := listExternalSnapshots(domain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 1 || names[0] != "ext-1" {
		t.Fatalf("expected [ext-1], got %v", names)
	}
}

func TestMergeExternalSnapshot_InactiveDomain(t *testing.T) {
	domain := &mockExternalSnapshotDomain{
		isActive: false,
		xmlDesc:  `<domain><devices></devices></domain>`,
	}

	_, err := mergeExternalSnapshot(domain, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.InvalidParameter) {
		t.Fatalf("expected InvalidParameter, got %v", err)
	}
}

func TestRevertExternalSnapshot_RequiresName(t *testing.T) {
	domain := &mockExternalSnapshotDomain{}
	err := revertExternalSnapshot(domain, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.InvalidParameter) {
		t.Fatalf("expected InvalidParameter, got %v", err)
	}
}
