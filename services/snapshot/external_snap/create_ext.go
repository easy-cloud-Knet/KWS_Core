package external

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
	"github.com/easy-cloud-Knet/KWS_Core/internal/qemuimg"
)

func validateCreateParams(domain *domCon.Domain, name string) error {
	if domain == nil || domain.Domain == nil {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("nil domain"))
	}
	if name == "" {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("snapshot name required"))
	}
	if !isSafeSnapshotName(name) {
		return virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("invalid snapshot name"))
	}
	return nil
}

func CreateExternalSnapshot(domain *domCon.Domain, name string, opts *ExternalSnapshotOptions) (string, error) {
	if err := validateCreateParams(domain, name); err != nil {
		return "", err
	}

	xmlDesc, err := domain.Domain.GetXMLDesc(0)
	if err != nil {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to get domain xml: %w", err))
	}

	domainUUID, err := domain.Domain.GetUUIDString()
	if err != nil {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to resolve domain uuid: %w", err))
	}

	return createExternalSnapshot(newExternalSnapshotDomain(domain.Domain), qemuimg.New(), domainUUID, xmlDesc, name, opts)
}

func createExternalSnapshot(domain SnapshotDomain, qimg qemuimg.Runner, domainUUID, xmlDesc, name string, opts *ExternalSnapshotOptions) (string, error) {
	disks, err := listFileDisksFromXMLDesc(xmlDesc)
	if err != nil {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to list file disks: %w", err))
	}
	if len(disks) == 0 {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("no file-based disks found for snapshot"))
	}

	snapshotDir, err := prepareSnapshotDir(opts, domainUUID, name)
	if err != nil {
		return "", err
	}

	snapDisks, err := createOverlays(qimg, disks, snapshotDir)
	if err != nil {
		return "", err
	}

	description := "external snapshot created by DODDLE"
	if opts != nil && opts.Description != "" {
		description = opts.Description
	}

	snapXML, err := buildSnapshotXML(name, description, snapDisks)
	if err != nil {
		return "", err
	}

	return registerSnapshot(domain, snapXML)
}

func prepareSnapshotDir(opts *ExternalSnapshotOptions, domainUUID, name string) (string, error) {
	snapshotRoot, err := resolveSnapshotRoot(opts)
	if err != nil {
		return "", virerr.ErrorGen(virerr.InvalidParameter, fmt.Errorf("failed to resolve snapshot root: %w", err))
	}

	snapshotDir := filepath.Join(snapshotRoot, domainUUID, "snapshots", name)
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to create snapshot directory: %w", err))
	}

	return snapshotDir, nil
}

func createOverlays(qimg qemuimg.Runner, disks []diskInfo, snapshotDir string) ([]snapshotDisk, error) {
	snapDisks := make([]snapshotDisk, 0, len(disks))

	for _, d := range disks {
		overlayPath := filepath.Join(snapshotDir, fmt.Sprintf("%s.qcow2", d.TargetDev))

		backingFormat := d.Driver
		if backingFormat == "" {
			backingFormat = "qcow2"
		}

		if err := qimg.Create(d.Source, backingFormat, overlayPath); err != nil {
			var rollbackErrs []string
			for _, created := range snapDisks {
				if removeErr := qimg.RemoveOverlay(snapshotDir, created.Source.File); removeErr != nil {
					rollbackErrs = append(rollbackErrs, fmt.Sprintf("%s: %v", created.Source.File, removeErr))
				}
			}
			if len(rollbackErrs) > 0 {
				return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to create overlay for disk %s: %w; rollback also failed: %s", d.TargetDev, err, strings.Join(rollbackErrs, ", ")))
			}
			return nil, virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to create overlay for disk %s: %w", d.TargetDev, err))
		}

		var driver *snapshotDriver
		if d.Driver != "" {
			driver = &snapshotDriver{Type: d.Driver}
		}
		snapDisks = append(snapDisks, snapshotDisk{
			Name:     d.TargetDev,
			Snapshot: "external",
			Driver:   driver,
			Source:   &snapshotSource{File: overlayPath},
		})
	}

	return snapDisks, nil
}

func buildSnapshotXML(name, description string, snapDisks []snapshotDisk) (string, error) {
	snapXML := snapshotXML{
		Name:        name,
		Description: description,
		Disks:       snapshotDisks{Disks: snapDisks},
	}

	snapBytes, err := xml.Marshal(snapXML)
	if err != nil {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to build snapshot xml: %w", err))
	}

	return string(snapBytes), nil
}

func registerSnapshot(domain SnapshotDomain, snapXML string) (string, error) {
	snap, err := domain.RegisterExternalSnapshot(snapXML)
	if err != nil {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("failed to register external snapshot: %w", err))
	}
	defer snap.Free()

	snapName, err := snap.Name()
	if err != nil {
		return "", virerr.ErrorGen(virerr.SnapshotError, fmt.Errorf("snapshot registered but failed to read name: %w", err))
	}

	return snapName, nil
}
