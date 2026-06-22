package creation

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/easy-cloud-Knet/KWS_Core/internal/config"
	virerr "github.com/easy-cloud-Knet/KWS_Core/internal/error"
)

func (DB localConfigurer) CreateDiskImage(dirPath string, diskSize int) error {
	baseImage := fmt.Sprintf("%s/baseimg/%s", config.StorageBase, DB.VMDescription.OS)
	targetImage := filepath.Join(dirPath, fmt.Sprintf("%s.qcow2", DB.VMDescription.UUID))
	qemuImgCmd := exec.Command("qemu-img", "create",
		"-b", baseImage,
		"-f", "qcow2",
		"-F", "qcow2",
		targetImage, fmt.Sprintf("%dG", diskSize),
	)
	if err := qemuImgCmd.Run(); err != nil {
		errorDescription := fmt.Errorf("generating Disk image error, check duplicated uuid or lack of HD capacity, or validity for base img %s, %v", dirPath, err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errorDescription)
	}

	return nil
}

// ensureBaseImage checks if the base image exists locally.
// If not, it downloads it from presignedURL to the target path atomically.
func ensureBaseImage(path, presignedURL string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	if presignedURL == "" {
		return fmt.Errorf("base image not found at %s and no presignedImageUrl provided", path)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create baseimg directory: %w", err)
	}

	tmpPath := path + ".tmp"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, presignedURL, nil)
	if err != nil {
		return fmt.Errorf("failed to build download request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download base image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download base image returned status %d", resp.StatusCode)
	}

	tmp, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to write base image: %w", err)
	}
	tmp.Close()

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to place base image: %w", err)
	}

	return nil
}

func (DB localConfigurer) CreateISOFile(dirPath string) error {

	isoOutput := filepath.Join(dirPath, "cidata.iso")
	userDataPath := filepath.Join(dirPath, "user-data")
	metaDataPath := filepath.Join(dirPath, "meta-data")

	genisoCmd := exec.Command("genisoimage",
		"--output", isoOutput,
		"-V", "cidata",
		"-r", "-J",
		userDataPath, metaDataPath,
	)

	if err := genisoCmd.Run(); err != nil {
		errorDescription := fmt.Errorf("generating ISO image error, may have duplicated uuid or wrong format of yaml file %s, %v", dirPath, err)
		return virerr.ErrorGen(virerr.DomainGenerationError, errorDescription)
	}
	return nil
}

