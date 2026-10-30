package website

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireSignIn(t *testing.T) {
	tests := []struct {
		name string
		uri  string
		want string
	}{
		{name: "path with query preserved", uri: "/years/2026/submit?x=1", want: "/auth/signin?redirect=%2Fyears%2F2026%2Fsubmit%3Fx%3D1"},
		{name: "plain path preserved", uri: "/profile/pies/2026", want: "/auth/signin?redirect=%2Fprofile%2Fpies%2F2026"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.uri, nil)

			requireSignIn(w, req)

			if w.Code != http.StatusSeeOther {
				t.Errorf("status = %d, want %d", w.Code, http.StatusSeeOther)
			}
			if got := w.Header().Get("Location"); got != tt.want {
				t.Errorf("Location = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSafeRedirectPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain path", in: "/years/2026/brands/brand/pie", want: "/years/2026/brands/brand/pie"},
		{name: "path with query", in: "/years/2026/user-pies/abc?page=2", want: "/years/2026/user-pies/abc?page=2"},
		{name: "root", in: "/", want: "/"},
		{name: "empty rejected", in: "", want: ""},
		{name: "relative rejected", in: "years/2026", want: ""},
		{name: "protocol relative rejected", in: "//evil.example.com", want: ""},
		{name: "backslash protocol relative rejected", in: `/\evil.example.com`, want: ""},
		{name: "absolute url rejected", in: "https://evil.example.com/x", want: ""},
		{name: "scheme with path rejected", in: "/redirect?to=https://evil.example.com", want: "/redirect?to=https://evil.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safeRedirectPath(tt.in); got != tt.want {
				t.Errorf("safeRedirectPath(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
