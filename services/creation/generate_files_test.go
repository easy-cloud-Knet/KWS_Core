package creation

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureBaseImage_AlreadyExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "base.img")
	if err := os.WriteFile(path, []byte("existing"), 0644); err != nil {
		t.Fatal(err)
	}

	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer srv.Close()

	if err := ensureBaseImage(path, srv.URL); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("HTTP call made when file already exists")
	}
}

func TestEnsureBaseImage_MissingNoURL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.img")
	if err := ensureBaseImage(path, ""); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestEnsureBaseImage_Downloads(t *testing.T) {
	content := []byte("fake-image-data")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))
	defer srv.Close()

	path := filepath.Join(t.TempDir(), "ubuntu-22.04")

	if err := ensureBaseImage(path, srv.URL); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("content mismatch: got %q, want %q", got, content)
	}
}

func TestEnsureBaseImage_NoTmpFileLeftOnServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "ubuntu-22.04")

	if err := ensureBaseImage(path, srv.URL); err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Error("tmp file not cleaned up after server error")
	}
}
