package website

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCsrfMiddleware checks the CrossOriginProtection wiring in
// csrfMiddleware: safe methods always pass, state-changing methods only pass
// when the request carries same-origin evidence (Sec-Fetch-Site or an Origin
// whose host matches Host), and everything else is rejected with 403.
func TestCsrfMiddleware(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := csrfMiddleware(ok)

	// A request to this host; httptest.NewRequest sets Host from the URL.
	const path = "/years/2026/submit"

	request := func(method, secFetchSite, origin string) *http.Request {
		var body *strings.Reader = nil
		if method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions {
			body = strings.NewReader("name=test")
		}
		var req *http.Request
		if body != nil {
			req = httptest.NewRequest(method, "http://mincepierank.co.uk"+path, body)
		} else {
			req = httptest.NewRequest(method, "http://mincepierank.co.uk"+path, nil)
		}
		if secFetchSite != "" {
			req.Header.Set("Sec-Fetch-Site", secFetchSite)
		}
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		return req
	}

	tests := []struct {
		name         string
		method       string
		secFetchSite string
		origin       string
		want         int
	}{
		// Safe methods are never blocked, even when they look cross-site:
		// nothing on the site may change state on GET/HEAD/OPTIONS.
		{"get no headers", http.MethodGet, "", "", http.StatusOK},
		{"get cross-site", http.MethodGet, "cross-site", "http://evil.example", http.StatusOK},
		{"options preflight", http.MethodOptions, "cross-site", "http://evil.example", http.StatusOK},

		// Modern browsers: Sec-Fetch-Site decides.
		{"post same-origin", http.MethodPost, "same-origin", "", http.StatusOK},
		{"put same-origin", http.MethodPut, "same-origin", "", http.StatusOK},
		{"delete same-origin", http.MethodDelete, "same-origin", "", http.StatusOK},
		{"post cross-site", http.MethodPost, "cross-site", "http://evil.example", http.StatusForbidden},
		{"put cross-site", http.MethodPut, "cross-site", "http://evil.example", http.StatusForbidden},
		{"delete cross-site", http.MethodDelete, "cross-site", "http://evil.example", http.StatusForbidden},

		// Older browsers without Sec-Fetch-Site: Origin host vs Host decides.
		{"post origin matches host", http.MethodPost, "", "http://mincepierank.co.uk", http.StatusOK},
		{"delete origin matches host", http.MethodDelete, "", "http://mincepierank.co.uk", http.StatusOK},
		{"post origin mismatch", http.MethodPost, "", "http://evil.example", http.StatusForbidden},
		{"put origin mismatch", http.MethodPut, "", "http://evil.example", http.StatusForbidden},

		// No browser headers at all: the stdlib treats the request as
		// non-browser (curl etc) and lets it through; the protection targets
		// browser-driven CSRF.
		{"post no headers", http.MethodPost, "", "", http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, request(tc.method, tc.secFetchSite, tc.origin))
			if rec.Code != tc.want {
				t.Fatalf("%s (sec-fetch-site=%q origin=%q): got %d, want %d",
					tc.method, tc.secFetchSite, tc.origin, rec.Code, tc.want)
			}
		})
	}
}
