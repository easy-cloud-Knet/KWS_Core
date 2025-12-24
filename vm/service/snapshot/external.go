package snapshot

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
)

// minimal XML structs to locate disk source
type domainXML struct {
	Devices devices `xml:"devices"`
}

type devices struct {
	Disks []disk `xml:"disk"`
}

type disk struct {
	Source diskSource `xml:"source"`
	Target diskTarget `xml:"target"`
}

type diskSource struct {
	File string `xml:"file,attr"`
}

type diskTarget struct {
	Dev string `xml:"dev,attr"`
}

func primaryDiskPathFromDomainXML(xmlDesc string) (string, error) {
	var d domainXML
	if err := xml.Unmarshal([]byte(xmlDesc), &d); err != nil {
		return "", fmt.Errorf("failed to parse domain xml: %w", err)
	}
	// prefer first disk with a file source
	for _, dk := range d.Devices.Disks {
		if dk.Source.File != "" {
			return dk.Source.File, nil
		}
	}
	return "", fmt.Errorf("no disk file source found in domain xml")
}

// ExportSnapshot converts the domain's primary disk to a separate qcow2 image at dest.
// dest can be a directory or full target path. This uses `qemu-img convert` and
// does not create libvirt-managed external snapshots; it performs an image copy.
func ExportSnapshot(domain *domCon.Domain, name, dest string) error {
	if domain == nil || domain.Domain == nil {
		return fmt.Errorf("nil domain")
	}

	xmlDesc, err := domain.Domain.GetXMLDesc(0)
	if err != nil {
		return fmt.Errorf("failed to get domain xml: %w", err)
	}

	src, err := primaryDiskPathFromDomainXML(xmlDesc)
	if err != nil {
		return err
	}

	// if dest is a directory, build a filename
	fi, err := os.Stat(dest)
	var outPath string
	if err == nil && fi.IsDir() {
		// create filename: <domain-uuid>-<name>.qcow2
		uuid, _ := domain.Domain.GetUUIDString()
		outPath = filepath.Join(dest, fmt.Sprintf("%s-%s.qcow2", uuid, name))
	} else {
		outPath = dest
	}

	// ensure destination dir exists
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("failed to ensure dest dir: %w", err)
	}

	// Run qemu-img convert -O qcow2 src outPath
	cmd := exec.Command("qemu-img", "convert", "-O", "qcow2", src, outPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("qemu-img convert failed: %w: %s", err, string(out))
	}

	return nil
}

// ImportSnapshot copies an external image into the domain storage directory and
// optionally returns the path where it was copied. It does not modify libvirt domain XML.
func ImportSnapshot(domain *domCon.Domain, name, src string) error {
	if domain == nil || domain.Domain == nil {
		return fmt.Errorf("nil domain")
	}

	uuid, err := domain.Domain.GetUUIDString()
	if err != nil {
		return fmt.Errorf("failed to get domain uuid: %w", err)
	}

	destDir := filepath.Join("/var/lib/kws/snapshots", uuid)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create dest dir: %w", err)
	}

	destPath := filepath.Join(destDir, filepath.Base(src))
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open src: %w", err)
	}
	defer in.Close()
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create dest file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	return nil
}
