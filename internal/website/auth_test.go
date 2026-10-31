package website

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
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

			RequireSignIn(w, req)

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

func TestUserIDMiddlewareAccessors(t *testing.T) {
	const signingKey = "test-signing-key"
	ac := &AuthConfig{Config: WebsiteConfiguration{JwtSigningKey: signingKey}, CookieKey: "jwt"}

	signed := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"email": "pie@fan.com"})
	signedJWT, err := signed.SignedString([]byte(signingKey))
	if err != nil {
		t.Fatalf("signing test jwt: %v", err)
	}

	tests := []struct {
		name         string
		cookies      []*http.Cookie
		wantUserID   string
		wantSignedIn bool
	}{
		{
			name:       "verified jwt carries email",
			cookies:    []*http.Cookie{{Name: "jwt", Value: signedJWT}},
			wantUserID: "pie@fan.com", wantSignedIn: true,
		},
		{
			name:       "anon cookie carries anon id",
			cookies:    []*http.Cookie{{Name: "anon_user_id", Value: "anon-42"}},
			wantUserID: "anon-42", wantSignedIn: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/years/2026", nil)
			for _, c := range tt.cookies {
				req.AddCookie(c)
			}
			w := httptest.NewRecorder()

			var gotUserID string
			var gotSignedIn bool
			h := ac.UserIdMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotUserID = UserIDFrom(r)
				gotSignedIn = SignedInFrom(r)
			}))
			h.ServeHTTP(w, req)

			if gotUserID != tt.wantUserID || gotSignedIn != tt.wantSignedIn {
				t.Errorf("accessors = (%q, %v), want (%q, %v)", gotUserID, gotSignedIn, tt.wantUserID, tt.wantSignedIn)
			}
		})
	}

	t.Run("no cookies mint a fresh anon id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/years/2026", nil)
		w := httptest.NewRecorder()

		var gotUserID string
		h := ac.UserIdMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotUserID = UserIDFrom(r)
			if SignedInFrom(r) {
				t.Error("anonymous request reported as signed in")
			}
		}))
		h.ServeHTTP(w, req)

		if gotUserID == "" {
			t.Error("expected a generated anon id")
		}
		var cookieSet bool
		for _, c := range w.Result().Cookies() {
			if c.Name == "anon_user_id" && c.Value == gotUserID {
				cookieSet = true
			}
		}
		if !cookieSet {
			t.Errorf("anon id %q was not sent back as a cookie", gotUserID)
		}
	})

	t.Run("accessors are zero without the middleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/years/2026", nil)
		if got := UserIDFrom(req); got != "" {
			t.Errorf("UserIDFrom = %q, want empty", got)
		}
		if SignedInFrom(req) {
			t.Error("SignedInFrom = true, want false")
		}
	})
}
