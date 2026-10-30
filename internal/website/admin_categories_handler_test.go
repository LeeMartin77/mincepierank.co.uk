package website

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

// fakeCategoryStore is an in-memory stand-in covering just the methods the
// category handlers touch; the embedded nil interfaces panic if anything else
// is hit.
type fakeCategoryStore struct {
	storage.Operations
	querier   fakeCategoryQuerier
	cats      map[uuid.UUID]types.Category
	pieCount  map[uuid.UUID]int64
	nextCount int
}

func newFakeCategoryStore() *fakeCategoryStore {
	return &fakeCategoryStore{
		cats:     map[uuid.UUID]types.Category{},
		pieCount: map[uuid.UUID]int64{},
	}
}

type fakeCategoryQuerier struct {
	sqlcgen.Querier
	isAdmin  bool
	activeYe int64
}

func (f fakeCategoryQuerier) IsAdminId(ctx context.Context, userID string) (bool, error) {
	return f.isAdmin, nil
}

func (f fakeCategoryQuerier) GetActiveYear(ctx context.Context) (int64, error) {
	return f.activeYe, nil
}

func (f *fakeCategoryStore) GetQuerier() sqlcgen.Querier { return f.querier }

func (f *fakeCategoryStore) GetUserPiesEnabled(c context.Context) (bool, error) { return false, nil }

func (f *fakeCategoryStore) GetAllCategoriesAdmin(c context.Context) (*[]storage.AdminCategoryRow, error) {
	rows := []storage.AdminCategoryRow{}
	for _, cat := range f.cats {
		rows = append(rows, storage.AdminCategoryRow{Category: cat, PieCount: f.pieCount[cat.Id]})
	}
	return &rows, nil
}

func (f *fakeCategoryStore) GetCategoryById(c context.Context, id uuid.UUID) (*types.Category, error) {
	cat, ok := f.cats[id]
	if !ok {
		return nil, nil
	}
	return &cat, nil
}

func (f *fakeCategoryStore) GetCategoryPieCount(c context.Context, id uuid.UUID) (int64, error) {
	return f.pieCount[id], nil
}

func (f *fakeCategoryStore) CreateCategory(c context.Context, cat types.Category) (*types.Category, error) {
	for _, existing := range f.cats {
		if existing.Slug == cat.Slug {
			return nil, storage.ErrCategoryAlreadyExists
		}
	}
	f.nextCount++
	cat.Id = uuid.UUID{byte(f.nextCount)}
	f.cats[cat.Id] = cat
	return &cat, nil
}

func (f *fakeCategoryStore) UpdateCategory(c context.Context, cat types.Category) (*types.Category, error) {
	if _, exists := f.cats[cat.Id]; !exists {
		return nil, storage.ErrCategoryNotFound
	}
	f.cats[cat.Id] = cat
	return &cat, nil
}

func (f *fakeCategoryStore) DeleteCategory(c context.Context, id uuid.UUID) error {
	if _, exists := f.cats[id]; !exists {
		return storage.ErrCategoryNotFound
	}
	delete(f.cats, id)
	return nil
}

func newCategoryTestRouter(t *testing.T, store *fakeCategoryStore) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.SetHTMLTemplate(templater.NewHtmlTemplater().GetTemplater())
	r.Use(func(c *gin.Context) {
		c.Keys = map[any]any{"userid": "admin@example.com"}
	})
	generated.RegisterHandlers(r, NewWebsiteWrapper(store, WebsiteConfiguration{
		ImgprssrDir: t.TempDir(),
	}))
	return r
}

func TestAdminCategoriesAccessControl(t *testing.T) {
	store := newFakeCategoryStore()
	store.querier.isAdmin = false
	r := newCategoryTestRouter(t, store)

	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	for _, req := range []struct{ method, path string }{
		{"GET", "/admin/categories"},
		{"GET", "/admin/categories/new"},
		{"GET", "/admin/categories/" + id.String()},
		{"POST", "/admin/categories"},
		{"PUT", "/admin/categories/" + id.String()},
		{"DELETE", "/admin/categories/" + id.String()},
	} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(req.method, req.path, nil))
		if w.Code != http.StatusForbidden {
			t.Errorf("%s %s = %d, want 403 for non-admin", req.method, req.path, w.Code)
		}
	}
}

func TestAdminCategoryListAndNew(t *testing.T) {
	store := newFakeCategoryStore()
	store.querier.isAdmin, store.querier.activeYe = true, 2026
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	store.cats[id] = types.Category{Id: id, Slug: "vegetarian", Label: "Vegetarian"}
	store.pieCount[id] = 4
	r := newCategoryTestRouter(t, store)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/categories", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET /admin/categories = %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Vegetarian") || !strings.Contains(body, "vegetarian") || !strings.Contains(body, ">4<") {
		t.Errorf("list page should show the label, slug and pie count")
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/categories/new", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `action="/admin/categories"`) {
		t.Errorf("GET /admin/categories/new = %d, should render the create form", w.Code)
	}
}

func postCategoryForm(r *gin.Engine, method, path string, form url.Values) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	return w
}

func TestAdminCategoryCreate(t *testing.T) {
	store := newFakeCategoryStore()
	store.querier.isAdmin, store.querier.activeYe = true, 2026
	r := newCategoryTestRouter(t, store)

	// the supplied slug is kept and we redirect to the edit page
	w := postCategoryForm(r, "POST", "/admin/categories", url.Values{"label": {"Vegan Options"}, "slug": {"vegan-options"}})
	if w.Code != http.StatusSeeOther || !strings.HasPrefix(w.Header().Get("Location"), "/admin/categories/") {
		t.Fatalf("create should redirect to the new category, got %d %q", w.Code, w.Header().Get("Location"))
	}
	var created *types.Category
	for _, cat := range store.cats {
		if cat.Slug == "vegan-options" {
			created = &cat
		}
	}
	if created == nil {
		t.Fatalf("category not persisted with supplied slug: %v", store.cats)
	}

	// missing slug is rejected, never derived from the label
	w = postCategoryForm(r, "POST", "/admin/categories", url.Values{"label": {"Gluten Free"}})
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "Category slug is required") {
		t.Errorf("missing slug = %d, want 400 with required-slug message", w.Code)
	}

	// missing label is rejected with the form re-rendered
	w = postCategoryForm(r, "POST", "/admin/categories", url.Values{"label": {""}})
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "Category label is required") {
		t.Errorf("empty label = %d, want 400 with error message", w.Code)
	}

	// unsafe slug rejected
	w = postCategoryForm(r, "POST", "/admin/categories", url.Values{"label": {"X"}, "slug": {"../up"}})
	if w.Code != http.StatusBadRequest {
		t.Errorf("path traversal slug = %d, want 400", w.Code)
	}

	// duplicate slug rejected with a message, not a silent no-op
	w = postCategoryForm(r, "POST", "/admin/categories", url.Values{"label": {"Vegan Again"}, "slug": {"vegan-options"}})
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "already exists") {
		t.Errorf("duplicate slug = %d, want 400 with already-exists message", w.Code)
	}
}

func TestAdminCategoryEditAndDeleteGuard(t *testing.T) {
	store := newFakeCategoryStore()
	store.querier.isAdmin, store.querier.activeYe = true, 2026
	used := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	clean := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	store.cats[used] = types.Category{Id: used, Slug: "gluten-free", Label: "Gluten Free"}
	store.cats[clean] = types.Category{Id: clean, Slug: "spicy", Label: "Spicy"}
	store.pieCount[used] = 2
	r := newCategoryTestRouter(t, store)

	// edit form renders and the slug is displayed but disabled
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/categories/"+used.String(), nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `hx-put="/admin/categories/`+used.String()+`"`) {
		t.Errorf("GET edit page = %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `value="gluten-free" disabled`) {
		t.Errorf("edit page should show the slug as disabled, not editable")
	}

	// malformed id 400s, unknown id 404s
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/categories/not-a-uuid", nil))
	if w.Code != http.StatusBadRequest {
		t.Errorf("GET malformed id = %d, want 400", w.Code)
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/categories/44444444-4444-4444-4444-444444444444", nil))
	if w.Code != http.StatusNotFound {
		t.Errorf("GET unknown id = %d, want 404", w.Code)
	}

	// update works and keeps the slug
	w = postCategoryForm(r, "PUT", "/admin/categories/"+used.String(), url.Values{"label": {"Gluten Free (Baked)"}})
	if w.Code != http.StatusSeeOther {
		t.Fatalf("update = %d, want 303", w.Code)
	}
	if store.cats[used].Label != "Gluten Free (Baked)" || store.cats[used].Slug != "gluten-free" {
		t.Errorf("category not updated correctly: %+v", store.cats[used])
	}

	// empty label rejected
	w = postCategoryForm(r, "PUT", "/admin/categories/"+used.String(), url.Values{"label": {"  "}})
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "Category label is required") {
		t.Errorf("empty label update = %d, want 400 with error message", w.Code)
	}

	// delete refused while pies reference the category
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("DELETE", "/admin/categories/"+used.String(), nil))
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "Cannot delete") {
		t.Errorf("delete referenced category = %d, want 400 with guard message", w.Code)
	}
	if _, ok := store.cats[used]; !ok {
		t.Errorf("guarded delete should not remove the category")
	}

	// delete works when unreferenced
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("DELETE", "/admin/categories/"+clean.String(), nil))
	if w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/admin/categories" {
		t.Errorf("delete clean category = %d %q, want 303 to /admin/categories", w.Code, w.Header().Get("Location"))
	}
	if _, ok := store.cats[clean]; ok {
		t.Errorf("clean category should be gone")
	}
}
