package external

import (
	"errors"
	"testing"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

type mockExternalSnapshotDomain struct {
	createErr       error
	createHandle    SnapshotHandle
	listErr         error
	snapshots       []SnapshotHandle
	blockJobInfo    externalBlockJobInfo
	blockJobInfoErr error

	lastCreateOpts externalSnapshotCreateExecOptions
}

func (m *mockExternalSnapshotDomain) CreateExternalSnapshot(_ string, opts externalSnapshotCreateExecOptions) (SnapshotHandle, error) {
	m.lastCreateOpts = opts
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createHandle != nil {
		return m.createHandle, nil
	}
	return &mockExternalSnapshotHandle{name: "created-snap"}, nil
}

func (m *mockExternalSnapshotDomain) ListAllSnapshots() ([]SnapshotHandle, error) {
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
	_, err := createExternalSnapshot(domain, false, "test-uuid", `<domain><devices></devices></domain>`, "", nil)
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
		createHandle: &mockExternalSnapshotHandle{name: "snap-a"},
	}
	xmlDesc := `<domain><devices>
			<disk device='disk' type='file'><driver type='qcow2'/><source file='/vm/vda.qcow2'/><target dev='vda' bus='virtio'/></disk>
			<disk device='disk' type='file'><driver type='qcow2'/><source file='/vm/vdb.qcow2'/><target dev='vdb' bus='virtio'/></disk>
		</devices></domain>`

	name, err := createExternalSnapshot(domain, true, "test-uuid", xmlDesc, "snap-a", &ExternalSnapshotOptions{BaseDir: tmp, Live: true, Quiesce: true})
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
	domain := &mockExternalSnapshotDomain{}
	_, err := createExternalSnapshot(domain, false, "test-uuid", "<invalid", "snap", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

func TestDeleteExternalSnapshot_NotFound(t *testing.T) {
	domain := &mockExternalSnapshotDomain{
		snapshots: []SnapshotHandle{},
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
		snapshots: []SnapshotHandle{
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
	domain := &mockExternalSnapshotDomain{}

	_, err := mergeExternalSnapshot(domain, `<invalid`, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

func TestRevertExternalSnapshot_RequiresName(t *testing.T) {
	domain := &mockExternalSnapshotDomain{}
	err := revertExternalSnapshot(domain, `<domain><devices></devices></domain>`, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.InvalidParameter) {
		t.Fatalf("expected InvalidParameter, got %v", err)
	}
}
