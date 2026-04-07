package external

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func freeSnapshotHandles(snaps []SnapshotHandle) {
	for _, s := range snaps {
		s.Free()
	}
}

func findExternalSnapshotByName(snaps []SnapshotHandle, snapName string) (SnapshotHandle, error) {
	for i := range snaps {
		name, err := snaps[i].Name()
		if err != nil || name != snapName {
			continue
		}

		isExternal, err := isExternalSnapshot(snaps[i])
		if err != nil {
			return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to inspect snapshot %s: %w", snapName, err))
		}
		if !isExternal {
			return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot %s is not external", snapName))
		}

		return snaps[i], nil
	}

	return nil, nil
}

func waitBlockJobReady(domain SnapshotDomain, disk string, timeout time.Duration) error {
	if domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}

	deadline := time.Now().Add(timeout)
	for {
		job, err := domain.BlockJobInfo(disk)
		if err != nil {
			if time.Now().After(deadline) {
				return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("timeout waiting for block job on disk %s: %w", disk, err))
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if job.End > 0 && job.Cur >= job.End {
			return nil
		}

		if time.Now().After(deadline) {
			return virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("timeout waiting for block commit to complete on disk %s", disk))
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func listFileDisksFromXMLDesc(xmlDesc string) ([]diskInfo, error) {

	var doc domainXML
	if err := xml.Unmarshal([]byte(xmlDesc), &doc); err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to parse domain xml: %w", err))
	}

	out := make([]diskInfo, 0, len(doc.Devices.Disks))
	for _, d := range doc.Devices.Disks {
		if d.Device != "disk" || d.Type != "file" {
			continue
		}
		if d.Source == nil || d.Target == nil || d.Target.Dev == "" || d.Source.File == "" {
			continue
		}
		driverType := ""
		driverName := ""
		if d.Driver != nil {
			driverType = d.Driver.Type
			driverName = d.Driver.Name
		}
		backingSource := ""
		if d.BackingStore != nil && d.BackingStore.Source != nil {
			backingSource = d.BackingStore.Source.File
		}
		out = append(out, diskInfo{
			TargetDev:     d.Target.Dev,
			TargetBus:     d.Target.Bus,
			Source:        d.Source.File,
			BackingSource: backingSource,
			Driver:        driverType,
			DriverName:    driverName,
		})
	}

	return out, nil
}

func buildDiskDeviceXML(info diskInfo, source string) string {
	driverXML := ""
	if info.Driver != "" || info.DriverName != "" {
		driverXML = "<driver"
		if info.DriverName != "" {
			driverXML += fmt.Sprintf(" name='%s'", info.DriverName)
		}
		if info.Driver != "" {
			driverXML += fmt.Sprintf(" type='%s'", info.Driver)
		}
		driverXML += "/>"
	}

	targetXML := fmt.Sprintf("<target dev='%s'", info.TargetDev)
	if info.TargetBus != "" {
		targetXML += fmt.Sprintf(" bus='%s'", info.TargetBus)
	}
	targetXML += "/>"

	return fmt.Sprintf("<disk type='file' device='disk'>%s<source file='%s'/>%s</disk>", driverXML, source, targetXML)
}

func isExternalSnapshot(snapshot SnapshotHandle) (bool, error) {
	if snapshot == nil {
		return false, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil snapshot"))
	}

	xmlDesc, err := snapshot.XMLDesc()
	if err != nil {
		return false, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to get snapshot xml: %w", err))
	}

	var doc snapshotXML
	if err := xml.Unmarshal([]byte(xmlDesc), &doc); err != nil {
		return false, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to parse snapshot xml: %w", err))
	}

	for _, d := range doc.Disks.Disks {
		if strings.EqualFold(d.Snapshot, "external") {
			return true, nil
		}
	}

	return false, nil
}

func isSafeSnapshotName(name string) bool {
	if name == "" {
		return false
	}
	clean := filepath.Clean(name)
	if clean != name {
		return false
	}
	if strings.Contains(name, "..") {
		return false
	}
	if strings.ContainsAny(name, `/\\`) {
		return false
	}
	return true
}

func extractExternalSnapshotSources(snapshot SnapshotHandle) (map[string]string, error) {
	if snapshot == nil {
		return nil, virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil snapshot"))
	}

	xmlDesc, err := snapshot.XMLDesc()
	if err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to get snapshot xml: %w", err))
	}

	var doc snapshotXML
	if err := xml.Unmarshal([]byte(xmlDesc), &doc); err != nil {
		return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to parse snapshot xml: %w", err))
	}

	out := make(map[string]string)
	for _, d := range doc.Disks.Disks {
		if !strings.EqualFold(d.Snapshot, "external") {
			continue
		}
		if d.Name == "" || d.Source == nil || d.Source.File == "" {
			continue
		}
		out[d.Name] = d.Source.File
	}

	return out, nil
}
