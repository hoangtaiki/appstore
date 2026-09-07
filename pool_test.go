package appstore

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewCertPool exercises the live download from apple.com. It is gated behind
// APPSTORE_NET=1 so the default `go test` run stays offline and deterministic.
func TestNewCertPool(t *testing.T) {
	if os.Getenv("APPSTORE_NET") != "1" {
		t.Skip("set APPSTORE_NET=1 to run the live apple.com cert download test")
	}
	cp, err := NewCertPool()
	if err != nil {
		t.Fatalf("NewCertPool() error = %v", err)
	}
	if cp == nil || cp.GetCertPool() == nil {
		t.Fatal("expected a usable cert pool")
	}
	if n := len(cp.GetCertPool().Subjects()); n < 1 {
		t.Fatalf("expected at least 1 cert in pool, got %d", n)
	}
}

// TestNewCertPool_DoesNotTouchCWDCerts locks in the H3 fix: NewCertPool() must
// never delete or modify a `certs/` directory in the process working directory.
func TestNewCertPool_DoesNotTouchCWDCerts(t *testing.T) {
	// Point the source at a fail-fast address so no real network call happens.
	restore := srcUrl
	srcUrl = "http://127.0.0.1:0/"
	defer func() { srcUrl = restore }()

	tmp := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	// Simulate a server that keeps its own TLS material under ./certs.
	if err := os.MkdirAll("certs", 0o755); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join("certs", "server.pem")
	if err := os.WriteFile(sentinel, []byte("keep me"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Download will fail (expected); the pool is still valid from the embed.
	if _, err := NewCertPool(); err == nil {
		t.Log("note: download unexpectedly succeeded; sentinel check still applies")
	}

	data, err := os.ReadFile(sentinel)
	if err != nil {
		t.Fatalf("H3 REGRESSION: sentinel under ./certs was destroyed: %v", err)
	}
	if string(data) != "keep me" {
		t.Fatalf("H3 REGRESSION: sentinel under ./certs was modified: %q", data)
	}
}

// TestNewCertPool_NeverEmptyOnDownloadFailure proves the embedded pinned root
// remains available (and the failure is surfaced) when the refresh fails.
func TestNewCertPool_NeverEmptyOnDownloadFailure(t *testing.T) {
	restore := srcUrl
	srcUrl = "http://127.0.0.1:0/"
	defer func() { srcUrl = restore }()

	cp, err := NewCertPool()
	if err == nil {
		t.Fatal("expected a non-nil error when the download fails")
	}
	if cp == nil || cp.GetCertPool() == nil {
		t.Fatal("expected a usable pool despite the download failure")
	}
	if n := len(cp.GetCertPool().Subjects()); n == 0 {
		t.Fatal("expected the embedded pinned root(s) to remain in the pool")
	}
}
