package website

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAbortStatus(t *testing.T) {
	w := httptest.NewRecorder()
	AbortStatus(w, http.StatusNotFound)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
	if got := w.Body.String(); got != "Not Found\n" {
		t.Errorf("body = %q, want the plain status text", got)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("content-type = %q, want text/plain", ct)
	}
}

func TestAbortErr(t *testing.T) {
	w := httptest.NewRecorder()
	AbortErr(w, http.StatusServiceUnavailable, errors.New("secret internal detail"))

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", w.Code)
	}
	body := w.Body.String()
	if body != "Service Unavailable\n" {
		t.Errorf("body = %q, want the plain status text", body)
	}
	// the underlying error must be logged (checked by eye in dev) but never
	// reach the client
	if strings.Contains(body, "secret internal detail") {
		t.Error("AbortErr leaked the error to the response body")
	}
}
