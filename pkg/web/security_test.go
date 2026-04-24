package web

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeadersMiddlewareSetsHeaders(t *testing.T) {
	t.Parallel()

	handler := securityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	res := rr.Result()
	if got := res.Header.Get("Content-Security-Policy"); got != contentSecurityPolicy {
		t.Fatalf("unexpected CSP header: %q", got)
	}
	if got := res.Header.Get("Referrer-Policy"); got != "no-referrer" {
		t.Fatalf("unexpected Referrer-Policy header: %q", got)
	}
	if got := res.Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("unexpected X-Content-Type-Options header: %q", got)
	}
	if got := res.Header.Get("X-Frame-Options"); got != "DENY" {
		t.Fatalf("unexpected X-Frame-Options header: %q", got)
	}
	if got := res.Header.Get("Strict-Transport-Security"); got != "" {
		t.Fatalf("unexpected HSTS header on cleartext request: %q", got)
	}
}

func TestSecurityHeadersMiddlewareSetsHSTSForTLS(t *testing.T) {
	t.Parallel()

	handler := securityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.TLS = &tls.ConnectionState{}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if got := rr.Result().Header.Get("Strict-Transport-Security"); got != strictTransportPolicy {
		t.Fatalf("unexpected HSTS header: %q", got)
	}
}
