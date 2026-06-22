package external

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadToPresignedURL_Success(t *testing.T) {
	content := []byte("snapshot-data")

	var gotMethod string
	var gotContentLength int64
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotContentLength = r.ContentLength
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "snap.qcow2")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	if err := UploadToPresignedURL(context.Background(), path, srv.URL); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("expected PUT, got %s", gotMethod)
	}
	if gotContentLength != int64(len(content)) {
		t.Errorf("ContentLength: got %d, want %d", gotContentLength, len(content))
	}
	if string(gotBody) != string(content) {
		t.Errorf("body mismatch: got %q, want %q", gotBody, content)
	}
}

func TestUploadToPresignedURL_FileNotFound(t *testing.T) {
	err := UploadToPresignedURL(context.Background(), "/nonexistent/snap.qcow2", "http://127.0.0.1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUploadToPresignedURL_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "snap.qcow2")
	if err := os.WriteFile(path, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := UploadToPresignedURL(context.Background(), path, srv.URL); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSnapshotFilePath(t *testing.T) {
	got := SnapshotFilePath("/var/lib/kws", "vm-uuid", "snap1", "vda")
	want := "/var/lib/kws/vm-uuid/snapshots/snap1/vda.qcow2"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
