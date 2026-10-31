package website

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

// fakeBrandStore is an in-memory stand-in covering just the methods the brand
// handlers touch; the embedded nil interfaces panic if anything else is hit.
type fakeBrandStore struct {
	storage.Operations
	querier  fakeBrandQuerier
	brands   map[string]types.Maker
	pieCount map[string]int64
}

func newFakeBrandStore() *fakeBrandStore {
	return &fakeBrandStore{
		brands:   map[string]types.Maker{},
		pieCount: map[string]int64{},
	}
}

type fakeBrandQuerier struct {
	sqlcgen.Querier
	isAdmin  bool
	activeYe int64
}

func (f fakeBrandQuerier) IsAdminId(ctx context.Context, userID string) (bool, error) {
	return f.isAdmin, nil
}

func (f fakeBrandQuerier) GetActiveYear(ctx context.Context) (int64, error) {
	return f.activeYe, nil
}

func (f *fakeBrandStore) GetQuerier() sqlcgen.Querier { return f.querier }

func (f *fakeBrandStore) GetUserPiesEnabled(c context.Context) (bool, error) { return false, nil }

func (f *fakeBrandStore) GetMaker(ctx context.Context, id string) (*types.Maker, error) {
	m, ok := f.brands[id]
	if !ok {
		return nil, nil
	}
	return &m, nil
}

func (f *fakeBrandStore) GetAllMakers(c context.Context) (*[]storage.AdminMakerRow, error) {
	rows := []storage.AdminMakerRow{}
	for _, m := range f.brands {
		rows = append(rows, storage.AdminMakerRow{Maker: m, PieCount: f.pieCount[m.Id]})
	}
	return &rows, nil
}

func (f *fakeBrandStore) GetMakerPieCount(c context.Context, makerid string) (int64, error) {
	return f.pieCount[makerid], nil
}

func (f *fakeBrandStore) CreateMaker(c context.Context, m types.Maker) (*types.Maker, error) {
	if _, exists := f.brands[m.Id]; exists {
		return nil, storage.ErrBrandAlreadyExists
	}
	f.brands[m.Id] = m
	return &m, nil
}

func (f *fakeBrandStore) UpdateMaker(c context.Context, m types.Maker) (*types.Maker, error) {
	if _, exists := f.brands[m.Id]; !exists {
		return nil, storage.ErrBrandNotFound
	}
	f.brands[m.Id] = m
	return &m, nil
}

func (f *fakeBrandStore) DeleteMaker(c context.Context, makerid string) error {
	delete(f.brands, makerid)
	return nil
}

func newBrandTestRouter(t *testing.T, store *fakeBrandStore) http.Handler {
	t.Helper()
	mux := http.NewServeMux()
	generated.HandlerFromMux(NewWebsiteWrapper(store, WebsiteConfiguration{
		ImgprssrDir: t.TempDir(),
	}, templater.NewHtmlTemplater()), mux)
	// stand in for UserIdMiddleware: every request carries an admin user id
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKeyUserID, "admin@example.com")
		mux.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestAdminBrandsAccessControl(t *testing.T) {
	store := newFakeBrandStore()
	store.querier.isAdmin = false
	r := newBrandTestRouter(t, store)

	for _, req := range []struct{ method, path string }{
		{"GET", "/admin/brands"},
		{"GET", "/admin/brands/new"},
		{"GET", "/admin/brands/greggs"},
		{"POST", "/admin/brands"},
		{"PUT", "/admin/brands/greggs"},
		{"DELETE", "/admin/brands/greggs"},
	} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(req.method, req.path, nil))
		if w.Code != http.StatusForbidden {
			t.Errorf("%s %s = %d, want 403 for non-admin", req.method, req.path, w.Code)
		}
	}
}

func TestAdminBrandListAndNew(t *testing.T) {
	store := newFakeBrandStore()
	store.querier.isAdmin, store.querier.activeYe = true, 2026
	store.brands["greggs"] = types.Maker{Id: "greggs", Name: "Greggs"}
	store.pieCount["greggs"] = 4
	r := newBrandTestRouter(t, store)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/brands", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET /admin/brands = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Greggs") || !strings.Contains(w.Body.String(), ">4<") {
		t.Errorf("list page should show the brand and its pie count")
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/brands/new", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `action="/admin/brands"`) {
		t.Errorf("GET /admin/brands/new = %d, should render the create form", w.Code)
	}
}

func TestAdminBrandCreate(t *testing.T) {
	store := newFakeBrandStore()
	store.querier.isAdmin, store.querier.activeYe = true, 2026
	r := newBrandTestRouter(t, store)

	form := url.Values{"name": {"Bramwell's"}, "id": {"bramwell-s"}, "website": {"https://bramwells.example"}}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/admin/brands", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/admin/brands/bramwell-s" {
		t.Fatalf("create should redirect to the supplied id, got %d %q", w.Code, w.Header().Get("Location"))
	}
	if _, ok := store.brands["bramwell-s"]; !ok {
		t.Errorf("brand not persisted: %v", store.brands)
	}

	// missing id is rejected, never derived from the name
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/admin/brands", strings.NewReader(url.Values{"name": {"Bramwell's Other"}}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "Brand id is required") {
		t.Errorf("missing id = %d, want 400 with required-id message", w.Code)
	}

	// missing name is rejected with the form re-rendered
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/admin/brands", strings.NewReader(url.Values{"name": {""}}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "Brand name is required") {
		t.Errorf("empty name = %d, want 400 with error message", w.Code)
	}

	// unsafe id rejected
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/admin/brands", strings.NewReader(url.Values{"name": {"X"}, "id": {"../up"}}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("path traversal id = %d, want 400", w.Code)
	}

	// duplicate id rejected with a message, not a silent no-op
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/admin/brands", strings.NewReader(url.Values{"name": {"Bramwell Again"}, "id": {"bramwell-s"}}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "already exists") {
		t.Errorf("duplicate id = %d, want 400 with already-exists message", w.Code)
	}
}

func TestAdminBrandEditAndDeleteGuard(t *testing.T) {
	store := newFakeBrandStore()
	store.querier.isAdmin, store.querier.activeYe = true, 2026
	store.brands["greggs"] = types.Maker{Id: "greggs", Name: "Greggs"}
	store.brands["clean"] = types.Maker{Id: "clean", Name: "Clean Co"}
	store.pieCount["greggs"] = 2
	r := newBrandTestRouter(t, store)

	// edit form renders
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/brands/greggs", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `hx-put="/admin/brands/greggs"`) {
		t.Errorf("GET /admin/brands/greggs = %d", w.Code)
	}

	// unknown brand 404s
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/brands/nobody", nil))
	if w.Code != http.StatusNotFound {
		t.Errorf("GET unknown brand = %d, want 404", w.Code)
	}

	// update works
	form := url.Values{"name": {"Greggs plc"}, "website": {""}}
	w = httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/admin/brands/greggs", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusSeeOther {
		t.Fatalf("update = %d, want 303", w.Code)
	}
	if store.brands["greggs"].Name != "Greggs plc" {
		t.Errorf("brand not updated: %+v", store.brands["greggs"])
	}

	// delete refused while pies reference the brand
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("DELETE", "/admin/brands/greggs", nil))
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "Cannot delete") {
		t.Errorf("delete referenced brand = %d, want 400 with guard message", w.Code)
	}
	if _, ok := store.brands["greggs"]; !ok {
		t.Errorf("guarded delete should not remove the brand")
	}

	// delete works when unreferenced
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("DELETE", "/admin/brands/clean", nil))
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/admin/brands" {
		t.Errorf("delete clean brand = %d %q, want 303 to /admin/brands", w.Code, w.Header().Get("Location"))
	}
	if _, ok := store.brands["clean"]; ok {
		t.Errorf("clean brand should be gone")
	}
}
