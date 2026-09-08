package appstore

import (
	"context"
	"crypto/x509"
	"embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

//go:embed certs/*.cer
var certs embed.FS

// srcUrl is the Apple page listing the current certificate-authority certs.
// It is a var (not a const) so tests can point it at an unreachable host.
var srcUrl = "https://www.apple.com/certificateauthority/"

var certLinkPattern = regexp.MustCompile(`<a [^>]*href="([^"]+\.cer)"`)

// refreshTimeout bounds the entire refresh (listing fetch + all cert fetches),
// so NewCertPool cannot block for an unbounded time regardless of link count.
const refreshTimeout = 60 * time.Second

// refreshTransport is the RoundTripper used for the cert refresh. It is nil in
// production (so http.DefaultTransport is used) and injectable by tests.
var refreshTransport http.RoundTripper

// validateRedirect re-validates every redirect hop against the https + apple.com
// allowlist, so an allowlisted URL cannot redirect to an untrusted host whose
// response would then be parsed and trusted. It also caps the redirect chain.
func validateRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errors.New("appstore: stopped after 10 redirects")
	}
	return validateAppleURL(req.URL.String())
}

type CertPool struct {
	pool     *x509.CertPool
	initErr  error
	poolOnce sync.Once
}

// NewCertPool builds a trusted-root pool from the embedded Apple root cert(s) and
// then best-effort refreshes it with Apple's currently published CA certs.
//
// The returned *CertPool is always safe to use: even when the refresh fails
// (Apple unreachable, timeout, etc.) the pool still contains the embedded pinned
// root(s). In that case a non-nil error is returned alongside the usable pool so
// the failure is visible rather than silently swallowed - callers may log it and
// proceed, or treat it as fatal.
func NewCertPool() (*CertPool, error) {
	cp := &CertPool{}
	err := cp.Init()
	return cp, err
}

func (cp *CertPool) Init() error {
	// initErr is stored on the struct (not a local) so every caller - including
	// goroutines that block on poolOnce.Do rather than running it - observes the
	// same result. poolOnce.Do provides the happens-before for the read below.
	cp.poolOnce.Do(func() {
		cp.pool = x509.NewCertPool()
		// The embedded pinned root(s) are the guaranteed baseline. If none load,
		// the binary is broken - fail loudly rather than return an empty pool.
		if err := cp.loadEmbedded(); err != nil {
			cp.initErr = err
			return
		}
		// Best-effort refresh from Apple. On failure we keep the embedded-only
		// pool and surface the error.
		cp.initErr = cp.downloadCerts()
	})
	return cp.initErr
}

// loadEmbedded populates the pool from the compile-time embedded certs.
func (cp *CertPool) loadEmbedded() error {
	entries, err := certs.ReadDir("certs")
	if err != nil {
		return err
	}
	loaded := 0
	for _, entry := range entries {
		if entry.IsDir() || !entry.Type().IsRegular() {
			continue
		}
		raw, err := certs.ReadFile("certs/" + entry.Name())
		if err != nil {
			continue
		}
		if cp.pool.AppendCertsFromPEM(raw) {
			loaded++
			continue
		}
		if cert, err := x509.ParseCertificate(raw); err == nil {
			cp.pool.AddCert(cert)
			loaded++
		}
	}
	if loaded == 0 {
		return errors.New("appstore: no embedded certificates loaded")
	}
	return nil
}

// downloadCerts fetches Apple's current CA cert list and adds each cert to the
// pool. Certs are parsed in memory - nothing is written to disk. A single overall
// deadline bounds the whole refresh; a failure to fetch the listing page is
// returned, and per-cert failures are collected and surfaced (a single bad link
// cannot abort the whole refresh).
func (cp *CertPool) downloadCerts() error {
	ctx, cancel := context.WithTimeout(context.Background(), refreshTimeout)
	defer cancel()

	client := &http.Client{
		Transport:     refreshTransport,
		CheckRedirect: validateRedirect,
	}

	content, err := httpGetCtx(ctx, client, srcUrl)
	if err != nil {
		return fmt.Errorf("appstore: fetching cert list %q: %w", srcUrl, err)
	}

	matches := certLinkPattern.FindAllSubmatch(content, -1)
	if len(matches) == 0 {
		// A 200 with no matches usually means Apple changed the page layout;
		// surface it rather than silently refreshing nothing.
		return fmt.Errorf("appstore: cert list %q returned no .cer links", srcUrl)
	}

	// Best-effort per cert: a single bad cert must not drop the others, but the
	// failures are collected and surfaced so a partial refresh is not silent.
	seen := make(map[string]bool)
	var errs []string
	for _, match := range matches {
		certUrl, err := cp.constructCertUrl(string(match[1]))
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if seen[certUrl] {
			continue
		}
		seen[certUrl] = true
		if err := cp.downloadAndAddCert(ctx, client, certUrl); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("appstore: %d cert refresh error(s): %s", len(errs), strings.Join(errs, "; "))
	}
	return nil
}

// httpGetCtx performs a context-bound GET and returns the response body. A
// non-200 response is an error.
func httpGetCtx(ctx context.Context, client *http.Client, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("returned status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (cp *CertPool) constructCertUrl(certPath string) (string, error) {
	var raw string
	switch {
	case strings.HasPrefix(certPath, "/"):
		baseUrl, err := url.Parse(srcUrl)
		if err != nil {
			return "", err
		}
		baseUrl.Path = certPath
		raw = baseUrl.String()
	case strings.HasPrefix(certPath, "https://www.apple.com/"),
		strings.HasPrefix(certPath, "https://developer.apple.com/"):
		raw = certPath
	default:
		joined, err := url.JoinPath(srcUrl, certPath)
		if err != nil {
			return "", err
		}
		raw = joined
	}
	if err := validateAppleURL(raw); err != nil {
		return "", err
	}
	return raw, nil
}

// validateAppleURL enforces https and an Apple host allowlist so a tampered
// listing page cannot make us fetch (and trust) a cert from an arbitrary host.
func validateAppleURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme != "https" {
		return fmt.Errorf("appstore: refusing non-https cert url %q", raw)
	}
	switch strings.ToLower(u.Hostname()) {
	case "www.apple.com", "apple.com", "developer.apple.com":
		return nil
	default:
		return fmt.Errorf("appstore: refusing cert url from unexpected host %q", u.Hostname())
	}
}

// downloadAndAddCert fetches a single cert and adds it to the pool, parsing the
// bytes in memory (PEM first, DER fallback). Any failure (non-200, read, parse)
// is returned so the caller can surface it.
func (cp *CertPool) downloadAndAddCert(ctx context.Context, client *http.Client, certUrl string) error {
	raw, err := httpGetCtx(ctx, client, certUrl)
	if err != nil {
		return fmt.Errorf("appstore: fetching cert %q: %w", certUrl, err)
	}

	if cp.pool.AppendCertsFromPEM(raw) {
		return nil
	}
	cert, err := x509.ParseCertificate(raw)
	if err != nil {
		return fmt.Errorf("appstore: parsing cert from %q: %w", certUrl, err)
	}
	cp.pool.AddCert(cert)
	return nil
}

func (cp *CertPool) GetCertPool() *x509.CertPool {
	return cp.pool
}
