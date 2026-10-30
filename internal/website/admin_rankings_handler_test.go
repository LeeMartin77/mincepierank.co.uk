package website

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

// fakeRankingStore covers just the methods the ranking handlers touch; the
// embedded nil interfaces panic if anything else is hit.
type fakeRankingStore struct {
	storage.Operations
	querier      fakeRankingQuerier
	piesEnabled  bool
	brandPieYear *types.MakerPieYearly
	userPie      *storage.UserPieWithOwner
}

type fakeRankingQuerier struct {
	sqlcgen.Querier
	isAdmin       bool
	activeYe      int64
	brandPies     []sqlcgen.ListMakerPieRankingSummariesRow
	userPies      []sqlcgen.ListUserPieRankingSummariesRow
	brandRankings []sqlcgen.MakerPieRankingYearly
	userRankings  []sqlcgen.UserPieRankingYearly
}

func newFakeRankingStore() *fakeRankingStore {
	return &fakeRankingStore{}
}

func (f fakeRankingQuerier) IsAdminId(ctx context.Context, userID string) (bool, error) {
	return f.isAdmin, nil
}

func (f fakeRankingQuerier) GetActiveYear(ctx context.Context) (int64, error) {
	return f.activeYe, nil
}

func (f fakeRankingQuerier) ListMakerPieRankingSummaries(ctx context.Context, year int32) ([]sqlcgen.ListMakerPieRankingSummariesRow, error) {
	return f.brandPies, nil
}

func (f fakeRankingQuerier) ListUserPieRankingSummaries(ctx context.Context, year int32) ([]sqlcgen.ListUserPieRankingSummariesRow, error) {
	return f.userPies, nil
}

func (f *fakeRankingQuerier) ListMakerPieRankingsForPie(ctx context.Context, arg sqlcgen.ListMakerPieRankingsForPieParams) ([]sqlcgen.MakerPieRankingYearly, error) {
	var out []sqlcgen.MakerPieRankingYearly
	for _, r := range f.brandRankings {
		if r.Year == arg.Year && r.Makerid == arg.Makerid && r.Pieid == arg.Pieid {
			out = append(out, r)
		}
	}
	return out, nil
}

func (f *fakeRankingQuerier) ListUserPieRankingsForPie(ctx context.Context, arg sqlcgen.ListUserPieRankingsForPieParams) ([]sqlcgen.UserPieRankingYearly, error) {
	var out []sqlcgen.UserPieRankingYearly
	for _, r := range f.userRankings {
		if r.Year == arg.Year && r.Pieid == arg.Pieid {
			out = append(out, r)
		}
	}
	return out, nil
}

func (f *fakeRankingQuerier) DeleteMakerPieRanking(ctx context.Context, arg sqlcgen.DeleteMakerPieRankingParams) error {
	kept := f.brandRankings[:0]
	for _, r := range f.brandRankings {
		if r.Year == arg.Year && r.Makerid == arg.Makerid && r.Pieid == arg.Pieid && r.Userid == arg.Userid {
			continue
		}
		kept = append(kept, r)
	}
	f.brandRankings = kept
	return nil
}

func (f *fakeRankingQuerier) DeleteUserPieRanking(ctx context.Context, arg sqlcgen.DeleteUserPieRankingParams) error {
	kept := f.userRankings[:0]
	for _, r := range f.userRankings {
		if r.Year == arg.Year && r.Pieid == arg.Pieid && r.Userid == arg.Userid {
			continue
		}
		kept = append(kept, r)
	}
	f.userRankings = kept
	return nil
}

func (f *fakeRankingStore) GetQuerier() sqlcgen.Querier { return &f.querier }

func (f *fakeRankingStore) GetUserPiesEnabled(c context.Context) (bool, error) { return f.piesEnabled, nil }

func (f *fakeRankingStore) GetMakerPieYearly(ctx context.Context, year int32, makerid string, id string) (*types.MakerPieYearly, error) {
	return f.brandPieYear, nil
}

func (f *fakeRankingStore) GetUserPie(c context.Context, year int32, id string) (*storage.UserPieWithOwner, error) {
	return f.userPie, nil
}

func newRankingTestRouter(t *testing.T, store *fakeRankingStore) *gin.Engine {
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

func TestAdminRankingsAccessControl(t *testing.T) {
	store := newFakeRankingStore()
	store.querier.isAdmin = false
	r := newRankingTestRouter(t, store)

	for _, req := range []struct{ method, path string }{
		{"GET", "/admin/rankings"},
		{"GET", "/admin/rankings/brands/greggs/sausage-roll"},
		{"DELETE", "/admin/rankings/brands/greggs/sausage-roll/users/someone@example.com"},
		{"GET", "/admin/rankings/userpies/pie1"},
		{"DELETE", "/admin/rankings/userpies/pie1/users/someone@example.com"},
	} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(req.method, req.path, nil))
		if w.Code != http.StatusForbidden {
			t.Errorf("%s %s = %d, want 403 for non-admin", req.method, req.path, w.Code)
		}
	}
}

func TestAdminRankingsPicker(t *testing.T) {
	store := newFakeRankingStore()
	store.querier.isAdmin, store.querier.activeYe = true, 2026
	store.querier.brandPies = []sqlcgen.ListMakerPieRankingSummariesRow{
		{Year: 2026, Makerid: "greggs", Pieid: "sausage-roll",
			Displayname: pgtype.Text{String: "Sausage Roll", Valid: true},
			Makername:   pgtype.Text{String: "Greggs", Valid: true}, RankingCount: 7},
	}
	store.piesEnabled = true
	store.querier.userPies = []sqlcgen.ListUserPieRankingSummariesRow{
		{Year: 2026, Pieid: "mypie", Displayname: "Nans Pie", RankingCount: 2},
	}
	r := newRankingTestRouter(t, store)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/rankings", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET /admin/rankings = %d", w.Code)
	}
	body := w.Body.String()
	for _, want := range []string{"Greggs", "Sausage Roll", ">7<", "Nans Pie", "/admin/rankings/userpies/mypie"} {
		if !strings.Contains(body, want) {
			t.Errorf("picker page missing %q", want)
		}
	}
}

func TestAdminBrandPieRankingsReviewAndDelete(t *testing.T) {
	store := newFakeRankingStore()
	store.querier.isAdmin, store.querier.activeYe = true, 2026
	store.brandPieYear = &types.MakerPieYearly{Year: 2026, MakerId: "greggs", Id: "sausage-roll", DisplayName: "Sausage Roll"}
	store.querier.brandRankings = []sqlcgen.MakerPieRankingYearly{
		{Year: 2026, Makerid: "greggs", Pieid: "sausage-roll", Userid: "ate@pie.com",
			Pastry: 4, Filling: 5, Topping: 3, Looks: 2, Value: 1,
			Notes: pgtype.Text{String: "Soggy bottom", Valid: true}},
		{Year: 2026, Makerid: "greggs", Pieid: "sausage-roll", Userid: "other pie man",
			Pastry: 1, Filling: 1, Topping: 1, Looks: 1, Value: 1},
	}
	r := newRankingTestRouter(t, store)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/rankings/brands/greggs/sausage-roll", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET review = %d", w.Code)
	}
	body := w.Body.String()
	// notes are shown to admins, and the awkward userid is escaped in the
	// delete link
	for _, want := range []string{"Soggy bottom", "ate@pie.com", "/users/other%20pie%20man?year=2026"} {
		if !strings.Contains(body, want) {
			t.Errorf("review page missing %q", want)
		}
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("DELETE", "/admin/rankings/brands/greggs/sausage-roll/users/ate@pie.com", nil))
	if w.Code != http.StatusSeeOther || !strings.HasPrefix(w.Header().Get("Location"), "/admin/rankings/brands/greggs/sausage-roll") {
		t.Fatalf("delete = %d %q, want 303 back to review", w.Code, w.Header().Get("Location"))
	}
	if len(store.querier.brandRankings) != 1 || store.querier.brandRankings[0].Userid != "other pie man" {
		t.Errorf("wrong ranking left after delete: %+v", store.querier.brandRankings)
	}

	// unknown pie is 404
	store.brandPieYear = nil
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/rankings/brands/greggs/nope", nil))
	if w.Code != http.StatusNotFound {
		t.Errorf("unknown pie review = %d, want 404", w.Code)
	}
}

func TestAdminUserPieRankingsReviewAndDelete(t *testing.T) {
	store := newFakeRankingStore()
	store.querier.isAdmin, store.querier.activeYe = true, 2026
	store.userPie = &storage.UserPieWithOwner{Pie: types.UserPieYearly{Year: 2026, Id: "mypie", DisplayName: "Nans Pie"}}
	store.querier.userRankings = []sqlcgen.UserPieRankingYearly{
		{Year: 2026, Pieid: "mypie", Userid: "rater@example.com", Notes: pgtype.Text{String: "Mushy", Valid: true}},
	}

	// user pie screens 404 while the feature is off
	store.piesEnabled = false
	r := newRankingTestRouter(t, store)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/rankings/userpies/mypie", nil))
	if w.Code != http.StatusNotFound {
		t.Errorf("pies disabled review = %d, want 404", w.Code)
	}

	store.piesEnabled = true
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/rankings/userpies/mypie", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("GET review = %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Nans Pie") || !strings.Contains(body, "Mushy") || !strings.Contains(body, "rater@example.com") {
		t.Errorf("review page missing pie name, note or rater")
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("DELETE", "/admin/rankings/userpies/mypie/users/rater@example.com?year=2026", nil))
	if w.Code != http.StatusSeeOther {
		t.Fatalf("delete = %d, want 303", w.Code)
	}
	if len(store.querier.userRankings) != 0 {
		t.Errorf("ranking not deleted: %+v", store.querier.userRankings)
	}

	// missing pie is 404
	store.userPie = nil
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/rankings/userpies/gone", nil))
	if w.Code != http.StatusNotFound {
		t.Errorf("missing pie review = %d, want 404", w.Code)
	}
}

func TestAdminRankingsYearParam(t *testing.T) {
	store := newFakeRankingStore()
	store.querier.isAdmin, store.querier.activeYe = true, 2026
	store.brandPieYear = &types.MakerPieYearly{Year: 2024, MakerId: "greggs", Id: "sausage-roll", DisplayName: "Sausage Roll"}
	r := newRankingTestRouter(t, store)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/rankings/brands/greggs/sausage-roll?year=2024", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "2024") {
		t.Errorf("?year= override not honoured: %d", w.Code)
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/admin/rankings/brands/greggs/sausage-roll?year=nonsense", nil))
	if w.Code != http.StatusBadRequest {
		t.Errorf("malformed year = %d, want 400", w.Code)
	}
}
