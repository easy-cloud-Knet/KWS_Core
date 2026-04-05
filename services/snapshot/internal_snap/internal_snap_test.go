package internal

import (
	"errors"
	"fmt"
	"testing"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

type mockInternalSnapshotDomain struct {
	createErr    error
	createHandle snapshotHandle
	listErr      error
	snapshots    []snapshotHandle

	lastCreateOpts snapshotCreateOptions
}

func (m *mockInternalSnapshotDomain) CreateSnapshot(_ string, opts snapshotCreateOptions) (snapshotHandle, error) {
	m.lastCreateOpts = opts
	if m.createErr != nil {
		return nil, m.createErr
	}
	if m.createHandle != nil {
		return m.createHandle, nil
	}
	return &mockInternalSnapshotHandle{name: "snap-created"}, nil
}

func (m *mockInternalSnapshotDomain) ListAllSnapshots() ([]snapshotHandle, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.snapshots, nil
}

type mockInternalSnapshotHandle struct {
	name      string
	nameErr   error
	deleteErr error
	revertErr error
}

func (m *mockInternalSnapshotHandle) Name() (string, error) {
	if m.nameErr != nil {
		return "", m.nameErr
	}
	return m.name, nil
}

func (m *mockInternalSnapshotHandle) Delete() error {
	return m.deleteErr
}

func (m *mockInternalSnapshotHandle) Revert() error {
	return m.revertErr
}

func (m *mockInternalSnapshotHandle) Free() error {
	return nil
}

func TestCreateSnapshot_NilDomain(t *testing.T) {
	_, err := createSnapshot(nil, "snap-1", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.InvalidParameter) {
		t.Fatalf("expected InvalidParameter, got %v", err)
	}
}

func TestCreateSnapshot_QuiesceOption(t *testing.T) {
	domain := &mockInternalSnapshotDomain{createHandle: &mockInternalSnapshotHandle{name: "snap-1"}}
	name, err := createSnapshot(domain, "snap-1", &SnapshotOptions{Description: "desc", Quiesce: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "snap-1" {
		t.Fatalf("expected snap-1, got %s", name)
	}
	if !domain.lastCreateOpts.Quiesce {
		t.Fatal("expected Quiesce=true to be mapped")
	}
}

func TestDeleteSnapshot_NotFound(t *testing.T) {
	domain := &mockInternalSnapshotDomain{snapshots: []snapshotHandle{}}
	err := deleteSnapshot(domain, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

func TestListSnapshots_CollectsNames(t *testing.T) {
	domain := &mockInternalSnapshotDomain{snapshots: []snapshotHandle{
		&mockInternalSnapshotHandle{name: "snap-a"},
		&mockInternalSnapshotHandle{name: "snap-b"},
	}}
	names, err := listSnapshots(domain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 2 || names[0] != "snap-a" || names[1] != "snap-b" {
		t.Fatalf("unexpected names: %v", names)
	}
}

func TestRevertToSnapshot_SnapshotNotFound(t *testing.T) {
	domain := &mockInternalSnapshotDomain{snapshots: []snapshotHandle{
		&mockInternalSnapshotHandle{name: "other"},
	}}
	err := revertToSnapshot(domain, "target")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}

func TestRevertToSnapshot_RevertFailure(t *testing.T) {
	domain := &mockInternalSnapshotDomain{snapshots: []snapshotHandle{
		&mockInternalSnapshotHandle{name: "target", revertErr: fmt.Errorf("revert failed")},
	}}
	err := revertToSnapshot(domain, "target")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, virerr.SnapshotError) {
		t.Fatalf("expected SnapshotError, got %v", err)
	}
}
