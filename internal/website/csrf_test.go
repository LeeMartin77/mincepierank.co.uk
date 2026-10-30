package website

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// fakeCsrfStore keeps tokens in a map so the middleware can be tested
// without a live redis.
type fakeCsrfStore struct {
	tokens map[string]string
	err    error
}

func (f *fakeCsrfStore) Get(_ context.Context, sid string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.tokens[sid], nil
}

func (f *fakeCsrfStore) Set(_ context.Context, sid string, token string) error {
	if f.err != nil {
		return f.err
	}
	f.tokens[sid] = token
	return nil
}

func newCsrfTestRouter(store csrfStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ch := &csrfHandler{
		store:  store,
		config: WebsiteConfiguration{MincepierankDomain: "mincepierank.co.uk", MincepierankSecureDomain: false},
	}
	r.Use(ch.Middleware())

	r.GET("/page", func(c *gin.Context) {
		token, _ := c.Keys["csrf_token"].(string)
		c.String(http.StatusOK, token)
	})
	r.POST("/mutate", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.PUT("/mutate", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.DELETE("/mutate", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return r
}

// getSessionToken performs a GET and returns the session id cookie plus the
// token the middleware stored for it.
func getSessionToken(t *testing.T, r *gin.Engine, store csrfStore) (sid string, token string) {
	t.Helper()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/page", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET /page status = %d, want %d", w.Code, http.StatusOK)
	}

	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == csrfCookieName {
			sid = cookie.Value
		}
	}
	if sid == "" {
		t.Fatal("no csrf session cookie set")
	}
	token = w.Body.String()
	if token == "" {
		t.Fatal("no csrf token handed to the template")
	}
	stored, err := store.Get(context.Background(), sid)
	if err != nil {
		t.Fatalf("store.Get: %v", err)
	}
	if stored != token {
		t.Fatalf("stored token = %q, want %q", stored, token)
	}
	return sid, token
}

func postWithCookie(r *gin.Engine, method, path, sid string, body *strings.Reader, contentType string, header string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, body)
		req.Header.Set("Content-Type", contentType)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: sid})
	if header != "" {
		req.Header.Set(csrfHeader, header)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCsrfMiddlewareIssuesCookieAndToken(t *testing.T) {
	store := &fakeCsrfStore{tokens: map[string]string{}}
	r := newCsrfTestRouter(store)

	sid, token := getSessionToken(t, r, store)

	// A second GET with the same cookie must reuse the same session id and
	// token rather than minting new ones.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/page", nil)
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: sid})
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET /page status = %d, want %d", w.Code, http.StatusOK)
	}
	if w.Body.String() != token {
		t.Errorf("token on second GET = %q, want %q", w.Body.String(), token)
	}
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == csrfCookieName && cookie.Value != sid {
			t.Errorf("session cookie changed to %q, want %q", cookie.Value, sid)
		}
	}
}

func TestCsrfMiddlewareRejectsMutationsWithoutValidToken(t *testing.T) {
	store := &fakeCsrfStore{tokens: map[string]string{}}
	r := newCsrfTestRouter(store)
	sid, _ := getSessionToken(t, r, store)
	// A token that exists in the store but was minted for another session.
	store.tokens["another-sid"] = "another-token"

	tests := []struct {
		name        string
		method      string
		header      string
		body        *strings.Reader
		contentType string
	}{
		{name: "POST no token", method: http.MethodPost},
		{name: "PUT no token", method: http.MethodPut},
		{name: "DELETE no token", method: http.MethodDelete},
		{name: "POST empty token field", method: http.MethodPost,
			body: strings.NewReader("csrf_token="), contentType: "application/x-www-form-urlencoded"},
		{name: "POST wrong header token", method: http.MethodPost, header: "not-the-token"},
		{name: "POST wrong form field token", method: http.MethodPost,
			body: strings.NewReader("csrf_token=not-the-token"), contentType: "application/x-www-form-urlencoded"},
		{name: "POST token for another session", method: http.MethodPost, header: "another-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := postWithCookie(r, tt.method, "/mutate", sid, tt.body, tt.contentType, tt.header)
			if w.Code != http.StatusForbidden {
				t.Errorf("%s status = %d, want %d", tt.name, w.Code, http.StatusForbidden)
			}
		})
	}
}

func TestCsrfMiddlewareAcceptsValidToken(t *testing.T) {
	store := &fakeCsrfStore{tokens: map[string]string{}}
	r := newCsrfTestRouter(store)
	sid, token := getSessionToken(t, r, store)

	// header path: used by htmx via the hx-headers attribute on body
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		w := postWithCookie(r, method, "/mutate", sid, nil, "", token)
		if w.Code != http.StatusOK {
			t.Errorf("%s with header token status = %d, want %d", method, w.Code, http.StatusOK)
		}
	}

	// form field path: used by plain <form method="post"> submissions
	form := url.Values{"csrf_token": []string{token}, "note": []string{"hello"}}
	w := postWithCookie(r, http.MethodPost, "/mutate", sid,
		strings.NewReader(form.Encode()), "application/x-www-form-urlencoded", "")
	if w.Code != http.StatusOK {
		t.Errorf("POST with form field token status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCsrfMiddlewareStoreFailure(t *testing.T) {
	store := &fakeCsrfStore{tokens: map[string]string{}, err: errors.New("redis down")}
	r := newCsrfTestRouter(store)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/page", nil))
	if w.Code != http.StatusInternalServerError {
		t.Errorf("GET /page status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
