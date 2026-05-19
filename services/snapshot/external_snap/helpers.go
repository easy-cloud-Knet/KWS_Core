package external

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

// realQemuImg is the production implementation of QemuImg using exec.Command.
type realQemuImg struct{}

func newQemuImg() QemuImg {
	return &realQemuImg{}
}

func (q *realQemuImg) Create(backingFile, backingFormat, overlayPath string) error {
	args := []string{"create", "-f", "qcow2", "-b", backingFile}
	if backingFormat != "" {
		args = append(args, "-F", backingFormat)
	}
	args = append(args, overlayPath)

	out, err := exec.Command("qemu-img", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("qemu-img create failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

type qemuImgInfoResult struct {
	BackingFilename   string `json:"backing-filename"`
	BackingFileFormat string `json:"backing-filename-format"`
}

func (q *realQemuImg) Info(diskPath string) (backingFile, backingFormat string, err error) {
	out, execErr := exec.Command("qemu-img", "info", "--output=json", diskPath).CombinedOutput()
	if execErr != nil {
		return "", "", fmt.Errorf("qemu-img info failed: %w: %s", execErr, strings.TrimSpace(string(out)))
	}

	var result qemuImgInfoResult
	if jsonErr := json.Unmarshal(out, &result); jsonErr != nil {
		return "", "", fmt.Errorf("failed to parse qemu-img info output: %w", jsonErr)
	}

	return result.BackingFilename, result.BackingFileFormat, nil
}

func (q *realQemuImg) Convert(src, dst string) error {
	out, err := exec.Command("qemu-img", "convert", "-f", "qcow2", "-O", "qcow2", src, dst).CombinedOutput()
	if err != nil {
		return fmt.Errorf("qemu-img convert failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (q *realQemuImg) Commit(overlay, base string) error {
	out, err := exec.Command("qemu-img", "commit", "-b", base, overlay).CombinedOutput()
	if err != nil {
		return fmt.Errorf("qemu-img commit failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

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

// workingDiskPath derives the "working" overlay path from a snapshot overlay path.
// Snapshot layout: <root>/<uuid>/snapshots/<snapName>/<disk>.qcow2
// Working layout:  <root>/<uuid>/working/<disk>.qcow2
func workingDiskPath(snapOverlay, diskName string) string {
	snapNameDir := filepath.Dir(snapOverlay)
	snapshotsDir := filepath.Dir(snapNameDir)
	uuidDir := filepath.Dir(snapshotsDir)
	return filepath.Join(uuidDir, "working", diskName+".qcow2")
}

// isSnapshotOverlay reports whether path is one of the overlay files managed
// by this snapshot system (under snapshots/ or working/ directories).
func isSnapshotOverlay(path string) bool {
	return strings.Contains(path, "/snapshots/") || strings.Contains(path, "/working/")
}

// findOriginAndOverlays traverses the backing chain of topOverlay and returns
// the origin disk (the VM's own qcow2 that sits directly on top of the base
// image) together with the list of overlay paths above it.
// Returns an error if the chain has no origin (overlays backed directly by base).
func findOriginAndOverlays(qimg QemuImg, topOverlay string) (origin string, overlays []string, err error) {
	current := topOverlay
	for {
		backing, _, infoErr := qimg.Info(current)
		if infoErr != nil {
			return "", nil, fmt.Errorf("failed to query disk info for %s: %w", current, infoErr)
		}
		if backing == "" {
			return "", nil, fmt.Errorf("no VM origin disk found: chain root reached without an intermediate disk")
		}
		if !isSnapshotOverlay(backing) {
			// backing is a non-overlay file — verify it is origin (has its own backing)
			// rather than the base image itself (which would have no backing).
			backingOfBacking, _, infoErr := qimg.Info(backing)
			if infoErr != nil {
				return "", nil, fmt.Errorf("failed to query disk info for %s: %w", backing, infoErr)
			}
			if backingOfBacking == "" {
				return "", nil, fmt.Errorf("no VM origin disk found: overlays are backed directly by the base image")
			}
			overlays = append(overlays, current)
			return backing, overlays, nil
		}
		overlays = append(overlays, current)
		current = backing
	}
}
