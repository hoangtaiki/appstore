package appstore

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// roundTripFunc adapts a function into an http.RoundTripper so tests can serve
// canned responses for the cert refresh without any network.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func newTestResponse(req *http.Request, status int, body []byte, header http.Header) *http.Response {
	if header == nil {
		header = make(http.Header)
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     header,
		Request:    req,
	}
}

// newTestCA builds an in-memory self-signed CA cert. No ExtKeyUsage is set, so it
// verifies against a pool containing it under Verify's default KeyUsages.
func newTestCA(t *testing.T) (*x509.Certificate, []byte) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "appstore test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return cert, der
}

// rootTrusted reports whether the embedded pinned Apple root (defaultRootPEM) is
// trusted by pool. Used instead of the deprecated x509.CertPool.Subjects().
func rootTrusted(pool *x509.CertPool) error {
	block, _ := pem.Decode([]byte(defaultRootPEM))
	if block == nil {
		return errors.New("failed to decode defaultRootPEM")
	}
	root, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	_, err = root.Verify(x509.VerifyOptions{Roots: pool})
	return err
}

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
	if err := rootTrusted(cp.GetCertPool()); err != nil {
		t.Fatalf("expected the embedded pinned root to be trusted: %v", err)
	}
}

// TestNewCertPool_DoesNotTouchCWDCerts locks in the H3 fix: NewCertPool() must
// never delete or modify a `certs/` directory in the process working directory.
//
// Not parallel-safe: it mutates the package var srcUrl and the process working
// directory. Do not add t.Parallel().
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
	defer func() {
		if err := os.Chdir(wd); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	}()

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
//
// Not parallel-safe: it mutates the package var srcUrl. Do not add t.Parallel().
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
	if err := rootTrusted(cp.GetCertPool()); err != nil {
		t.Fatalf("expected the embedded pinned root to remain trusted: %v", err)
	}
}

// TestNewCertPool_RefreshAddsDownloadedCert deterministically covers the in-memory
// success path with no network: a stub transport serves a listing and a cert, and
// we assert the downloaded cert becomes trusted by the pool.
//
// Not parallel-safe: it mutates the package var refreshTransport.
func TestNewCertPool_RefreshAddsDownloadedCert(t *testing.T) {
	ca, der := newTestCA(t)

	restore := refreshTransport
	refreshTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/certificateauthority/":
			body := []byte(`<a href="https://www.apple.com/testca.cer">cert</a>`)
			return newTestResponse(req, http.StatusOK, body, nil), nil
		case "/testca.cer":
			return newTestResponse(req, http.StatusOK, der, nil), nil
		default:
			return newTestResponse(req, http.StatusNotFound, nil, nil), nil
		}
	})
	defer func() { refreshTransport = restore }()

	cp, err := NewCertPool()
	if err != nil {
		t.Fatalf("NewCertPool() error = %v", err)
	}
	if _, err := ca.Verify(x509.VerifyOptions{Roots: cp.GetCertPool()}); err != nil {
		t.Fatalf("downloaded cert was not added to the pool: %v", err)
	}
}

// TestNewCertPool_RejectsRedirectToOtherHost verifies a cert link that redirects
// off the Apple allowlist is refused by CheckRedirect, so the redirected-to
// payload is never fetched or trusted.
//
// Not parallel-safe: it mutates the package var refreshTransport.
func TestNewCertPool_RejectsRedirectToOtherHost(t *testing.T) {
	ca, der := newTestCA(t)

	restore := refreshTransport
	refreshTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Host {
		case "www.apple.com":
			switch req.URL.Path {
			case "/certificateauthority/":
				body := []byte(`<a href="https://www.apple.com/redir.cer">cert</a>`)
				return newTestResponse(req, http.StatusOK, body, nil), nil
			case "/redir.cer":
				h := make(http.Header)
				h.Set("Location", "https://evil.example.com/evil.cer")
				return newTestResponse(req, http.StatusFound, nil, h), nil
			}
		case "evil.example.com":
			// Must never be reached: CheckRedirect blocks the hop first.
			return newTestResponse(req, http.StatusOK, der, nil), nil
		}
		return newTestResponse(req, http.StatusNotFound, nil, nil), nil
	})
	defer func() { refreshTransport = restore }()

	cp, err := NewCertPool()
	if err == nil {
		t.Fatal("expected a non-nil error when a cert link redirects off-allowlist")
	}
	if _, verr := ca.Verify(x509.VerifyOptions{Roots: cp.GetCertPool()}); verr == nil {
		t.Fatal("redirected-to cert was unexpectedly trusted")
	}
}
