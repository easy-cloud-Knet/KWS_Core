package qemuimg

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Runner abstracts qemu-img command execution for testability.
type Runner interface {
	Create(backingFile, backingFormat, overlayPath string) error
	// ReplaceOverlay removes an existing file at overlayPath if present, then creates a fresh overlay.
	ReplaceOverlay(backingFile, backingFormat, overlayPath string) error
	Info(diskPath string) (backingFile, backingFormat string, err error)
	// Convert flattens the full backing chain of src into a new standalone dst file.
	// Unlike Commit, it never writes to any backing file so shared base images stay untouched.
	Convert(src, dst string) error
	// Commit writes all changes recorded in overlay into base, collapsing every
	// intermediate layer between them. Only base needs a write lock; files above
	// base in the chain are read-only during the operation.
	Commit(overlay, base string) error
}

type real struct{}

func New() Runner {
	return &real{}
}

func (q *real) Create(backingFile, backingFormat, overlayPath string) error {
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

func (q *real) ReplaceOverlay(backingFile, backingFormat, overlayPath string) error {
	if err := os.Remove(overlayPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove stale overlay %s: %w", overlayPath, err)
	}
	return q.Create(backingFile, backingFormat, overlayPath)
}

type infoResult struct {
	BackingFilename   string `json:"backing-filename"`
	BackingFileFormat string `json:"backing-filename-format"`
}

func (q *real) Info(diskPath string) (backingFile, backingFormat string, err error) {
	out, execErr := exec.Command("qemu-img", "info", "--output=json", diskPath).CombinedOutput()
	if execErr != nil {
		return "", "", fmt.Errorf("qemu-img info failed: %w: %s", execErr, strings.TrimSpace(string(out)))
	}

	var result infoResult
	if jsonErr := json.Unmarshal(out, &result); jsonErr != nil {
		return "", "", fmt.Errorf("failed to parse qemu-img info output: %w", jsonErr)
	}

	return result.BackingFilename, result.BackingFileFormat, nil
}

func (q *real) Convert(src, dst string) error {
	out, err := exec.Command("qemu-img", "convert", "-f", "qcow2", "-O", "qcow2", src, dst).CombinedOutput()
	if err != nil {
		return fmt.Errorf("qemu-img convert failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (q *real) Commit(overlay, base string) error {
	out, err := exec.Command("qemu-img", "commit", "-b", base, overlay).CombinedOutput()
	if err != nil {
		return fmt.Errorf("qemu-img commit failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
