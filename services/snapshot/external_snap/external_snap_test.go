package external

import (
	"errors"
	"testing"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

// --- mock SnapshotDomain ---

type mockExternalSnapshotDomain struct {
	registerErr    error
	registerHandle SnapshotHandle
	listErr        error
	snapshots      []SnapshotHandle
	updateErr      error
}

func (m *mockExternalSnapshotDomain) RegisterExternalSnapshot(_ string) (SnapshotHandle, error) {
	if m.registerErr != nil {
		return nil, m.registerErr
	}
	if m.registerHandle != nil {
		return m.registerHandle, nil
	}
	return &mockExternalSnapshotHandle{name: "created-snap"}, nil
}

func (m *mockExternalSnapshotDomain) ListAllSnapshots() ([]SnapshotHandle, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.snapshots, nil
}

func (m *mockExternalSnapshotDomain) UpdateDeviceConfig(_ string) error {
	return m.updateErr
}

// --- mock SnapshotHandle ---

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

// --- mock QemuImg ---

type mockQemuImg struct {
	createErr  error
	infoErr    error
	infoFunc   func(string) (string, string, error)
	backing    string
	backingFmt string
	convertErr error
	commitErr  error
}

func (q *mockQemuImg) Create(_, _, _ string) error {
	return q.createErr
}

func (q *mockQemuImg) Info(path string) (string, string, error) {
	if q.infoFunc != nil {
		return q.infoFunc(path)
	}
	if q.infoErr != nil {
		return "", "", q.infoErr
	}
	return q.backing, q.backingFmt, nil
}

func (q *mockQemuImg) Convert(_, _ string) error {
	return q.convertErr
}

func (q *mockQemuImg) Commit(_, _ string) error {
	return q.commitErr
}

// --- Create tests ---

func TestCreateExternalSnapshot_InvalidName(t *testing.T) {
	domain := &mockExternalSnapshotDomain{}
	qimg := &mockQemuImg{}
	_, err := createExternalSnapshot(domain, qimg, "test-uuid", `<domain><devices></devices></domain>`, "", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.InvalidParameter) {
		t.Fatalf("expected InvalidParameter, got %v", err)
	}
}

func TestCreateExternalSnapshot_XMLDescError(t *testing.T) {
	domain := &mockExternalSnapshotDomain{}
	qimg := &mockQemuImg{}
	_, err := createExternalSnapshot(domain, qimg, "test-uuid", "<invalid", "snap", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

func TestCreateExternalSnapshot_CreatesOverlayAndRegisters(t *testing.T) {
	tmp := t.TempDir()
	domain := &mockExternalSnapshotDomain{
		registerHandle: &mockExternalSnapshotHandle{name: "snap-a"},
	}
	qimg := &mockQemuImg{}
	xmlDesc := `<domain><devices>
		<disk device='disk' type='file'><driver type='qcow2'/><source file='/vm/vda.qcow2'/><target dev='vda' bus='virtio'/></disk>
	</devices></domain>`

	name, err := createExternalSnapshot(domain, qimg, "test-uuid", xmlDesc, "snap-a", &ExternalSnapshotOptions{BaseDir: tmp})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "snap-a" {
		t.Fatalf("expected snapshot name snap-a, got %s", name)
	}
}

func TestCreateExternalSnapshot_QemuImgError(t *testing.T) {
	tmp := t.TempDir()
	domain := &mockExternalSnapshotDomain{}
	qimg := &mockQemuImg{createErr: errors.New("qemu-img failed")}
	xmlDesc := `<domain><devices>
		<disk device='disk' type='file'><driver type='qcow2'/><source file='/vm/vda.qcow2'/><target dev='vda' bus='virtio'/></disk>
	</devices></domain>`

	_, err := createExternalSnapshot(domain, qimg, "test-uuid", xmlDesc, "snap-a", &ExternalSnapshotOptions{BaseDir: tmp})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

// --- Revert tests ---

func TestRevertExternalSnapshot_RequiresName(t *testing.T) {
	domain := &mockExternalSnapshotDomain{}
	qimg := &mockQemuImg{}
	err := revertExternalSnapshot(domain, qimg, `<domain><devices></devices></domain>`, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.InvalidParameter) {
		t.Fatalf("expected InvalidParameter, got %v", err)
	}
}

func TestRevertExternalSnapshot_SnapshotNotFound(t *testing.T) {
	domain := &mockExternalSnapshotDomain{snapshots: []SnapshotHandle{}}
	qimg := &mockQemuImg{}
	err := revertExternalSnapshot(domain, qimg, `<domain><devices></devices></domain>`, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

// --- Delete tests ---

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

// --- List tests ---

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

// --- Merge tests ---

func TestMergeExternalSnapshot_InvalidXML(t *testing.T) {
	domain := &mockExternalSnapshotDomain{}
	qimg := &mockQemuImg{}

	_, err := mergeExternalSnapshot(domain, qimg, `<invalid`, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

func TestMergeExternalSnapshot_NoBacking(t *testing.T) {
	domain := &mockExternalSnapshotDomain{}
	// Info returns empty backing — disk is already a base, nothing to merge.
	qimg := &mockQemuImg{backing: ""}
	xmlDesc := `<domain><devices>
		<disk device='disk' type='file'><driver type='qcow2'/><source file='/vm/vda.qcow2'/><target dev='vda' bus='virtio'/></disk>
	</devices></domain>`

	_, err := mergeExternalSnapshot(domain, qimg, xmlDesc, "")
	if err == nil {
		t.Fatal("expected error for no mergeable disks, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

func TestMergeExternalSnapshot_NoOrigin(t *testing.T) {
	domain := &mockExternalSnapshotDomain{}
	// Chain: base ← /snapshots/ex-1/vda.qcow2 (overlay backed directly by base, no origin disk)
	infoResults := map[string]struct{ b, f string }{
		"/vm/snapshots/ex-1/vda.qcow2": {"/vm/base.qcow2", "qcow2"},
		"/vm/base.qcow2":               {"", ""},
	}
	qimg := &mockQemuImg{
		infoFunc: func(path string) (string, string, error) {
			r := infoResults[path]
			return r.b, r.f, nil
		},
	}
	xmlDesc := `<domain><devices>
		<disk device='disk' type='file'><driver type='qcow2'/><source file='/vm/snapshots/ex-1/vda.qcow2'/><target dev='vda' bus='virtio'/></disk>
	</devices></domain>`

	_, err := mergeExternalSnapshot(domain, qimg, xmlDesc, "")
	if err == nil {
		t.Fatal("expected error for chain with no origin, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

func TestMergeExternalSnapshot_CommitsIntoOrigin(t *testing.T) {
	domain := &mockExternalSnapshotDomain{}
	// Chain: base ← /vm/origin.qcow2 ← /vm/snapshots/ex-1/vda.qcow2 ← /vm/snapshots/ex-2/vda.qcow2
	infoResults := map[string]struct{ b, f string }{
		"/vm/snapshots/ex-2/vda.qcow2": {"/vm/snapshots/ex-1/vda.qcow2", "qcow2"},
		"/vm/snapshots/ex-1/vda.qcow2": {"/vm/origin.qcow2", "qcow2"},
		"/vm/origin.qcow2":             {"/vm/base.qcow2", "qcow2"},
		"/vm/base.qcow2":               {"", ""},
	}
	qimg := &mockQemuImg{
		infoFunc: func(path string) (string, string, error) {
			r := infoResults[path]
			return r.b, r.f, nil
		},
	}
	xmlDesc := `<domain><devices>
		<disk device='disk' type='file'><driver type='qcow2'/><source file='/vm/snapshots/ex-2/vda.qcow2'/><target dev='vda' bus='virtio'/></disk>
	</devices></domain>`

	disks, err := mergeExternalSnapshot(domain, qimg, xmlDesc, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(disks) != 1 || disks[0] != "vda" {
		t.Fatalf("expected [vda], got %v", disks)
	}
}
