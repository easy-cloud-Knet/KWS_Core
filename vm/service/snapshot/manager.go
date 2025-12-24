package snapshot

import (
	"fmt"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	"libvirt.org/go/libvirt"
)

// SnapshotManager defines high-level snapshot operations.
type SnapshotManager interface {
	Create(uuid, name string) (string, error)
	List(uuid string) ([]string, error)
	Revert(uuid, name string) error
	Delete(uuid, name string) error
	Export(uuid, name, dest string) error
	Import(uuid, name, src string) error
}

// NewDefaultManager returns a simple manager backed by filesystem store.
func NewDefaultManager(basePath string) SnapshotManager {
	return &defaultManager{basePath: basePath}
}

type defaultManager struct {
	basePath string
	domCtl   *domCon.DomListControl
	conn     *libvirt.Connect
}

// Note: defaultManager methods can be implemented to coordinate operations across
// operations.go / store.go / external.go. For now helper functions in operations.go
// and store.go are available for direct use.

// Basic defaultManager implementations that delegate to package helpers.
// These provide minimal functionality so the manager satisfies the interface.
func (m *defaultManager) Create(uuid, name string) (string, error) {
	if m.domCtl == nil || m.conn == nil {
		return "", fmt.Errorf("snapshot manager not initialized with domCtl/libvirt connection")
	}
	dom, err := m.domCtl.GetDomain(uuid, m.conn)
	if err != nil {
		return "", err
	}
	return CreateSnapshot(dom, name)
}

func (m *defaultManager) List(uuid string) ([]string, error) {
	if m.domCtl == nil || m.conn == nil {
		return nil, fmt.Errorf("snapshot manager not initialized with domCtl/libvirt connection")
	}
	dom, err := m.domCtl.GetDomain(uuid, m.conn)
	if err != nil {
		return nil, err
	}
	return ListSnapshots(dom)
}

func (m *defaultManager) Revert(uuid, name string) error {
	if m.domCtl == nil || m.conn == nil {
		return fmt.Errorf("snapshot manager not initialized with domCtl/libvirt connection")
	}
	dom, err := m.domCtl.GetDomain(uuid, m.conn)
	if err != nil {
		return err
	}
	return RevertToSnapshot(dom, name)
}

func (m *defaultManager) Delete(uuid, name string) error {
	if m.domCtl == nil || m.conn == nil {
		return fmt.Errorf("snapshot manager not initialized with domCtl/libvirt connection")
	}
	dom, err := m.domCtl.GetDomain(uuid, m.conn)
	if err != nil {
		return err
	}

	snaps, err := dom.Domain.ListAllSnapshots(0)
	if err != nil {
		return fmt.Errorf("failed to list snapshots: %w", err)
	}
	defer func() {
		for _, s := range snaps {
			s.Free()
		}
	}()

	for _, s := range snaps {
		n, err := s.GetName()
		if err != nil {
			continue
		}
		if n == name {
			if err := s.Delete(0); err != nil {
				return fmt.Errorf("failed to delete snapshot %s: %w", name, err)
			}
			return nil
		}
	}
	return fmt.Errorf("snapshot %s not found", name)
}

func (m *defaultManager) Export(uuid, name, dest string) error {
	if m.domCtl == nil || m.conn == nil {
		return fmt.Errorf("snapshot manager not initialized with domCtl/libvirt connection")
	}
	dom, err := m.domCtl.GetDomain(uuid, m.conn)
	if err != nil {
		return err
	}
	return ExportSnapshot(dom, name, dest)
}

func (m *defaultManager) Import(uuid, name, src string) error {
	if m.domCtl == nil || m.conn == nil {
		return fmt.Errorf("snapshot manager not initialized with domCtl/libvirt connection")
	}
	dom, err := m.domCtl.GetDomain(uuid, m.conn)
	if err != nil {
		return err
	}
	return ImportSnapshot(dom, name, src)
}

// ErrNotImplementedInManager is a simple sentinel for unimplemented manager ops.
var ErrNotImplementedInManager = fmt.Errorf("snapshot manager operation not implemented")

// NewManagerWithDeps creates a manager that has access to DomListControl and libvirt connection.
func NewManagerWithDeps(domCtl *domCon.DomListControl, conn *libvirt.Connect, basePath string) SnapshotManager {
	return &defaultManager{basePath: basePath, domCtl: domCtl, conn: conn}
}
