package website

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
	"github.com/rs/zerolog/log"
)

// rankingYear resolves the optional ?year= query parameter, falling back to
// the active year. The generated wrapper already rejects non-numeric years.
func rankingYear(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, paramYear *int) (int32, bool) {
	if paramYear != nil {
		return int32(*paramYear), true
	}
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		abortErr(w, 500, err)
		return 0, false
	}
	return int32(ay), true
}

// GetAdminRankings implements generated.ServerInterface. The picker lists
// every brand pie (and, when user pies are on, every public user pie) for the
// chosen year with a rating count, so admins can drill into a review screen.
func (wrpr *WebsiteWrapper) GetAdminRankings(w http.ResponseWriter, r *http.Request, params generated.GetAdminRankingsParams) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	year, ok := rankingYear(wrpr, w, r, params.Year)
	if !ok {
		return
	}

	brandPies, err := wrpr.storage.GetQuerier().ListMakerPieRankingSummaries(r.Context(), year)
	if err != nil {
		abortErr(w, 500, err)
		return
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	userPies := []sqlcgen.ListUserPieRankingSummariesRow{}
	if piesEnabled {
		userPies, err = wrpr.storage.GetQuerier().ListUserPieRankingSummaries(r.Context(), year)
		if err != nil {
			abortErr(w, 500, err)
			return
		}
	}

	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		abortErr(w, 500, err)
		return
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			IsAdminPage: true,
			Title:       "Admin Rankings",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"BrandPies":  brandPies,
			"UserPies":   userPies,
			"PageErrors": []string{},
		},
	}
	wrpr.renderPage(w, http.StatusOK, "admin:rankings", vals)
}

// GetAdminBrandPieRankings implements generated.ServerInterface. The review
// screen lists each rating one user gave this brand pie, including the notes
// that are collected but never shown publicly.
func (wrpr *WebsiteWrapper) GetAdminBrandPieRankings(w http.ResponseWriter, r *http.Request, brand string, pie string, params generated.GetAdminBrandPieRankingsParams) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	year, ok := rankingYear(wrpr, w, r, params.Year)
	if !ok {
		return
	}

	pieYearly, err := wrpr.storage.GetMakerPieYearly(r.Context(), year, brand, pie)
	if err != nil {
		abortErr(w, 500, err)
		return
	}
	if pieYearly == nil {
		abortStatus(w, http.StatusNotFound)
		return
	}

	rankings, err := wrpr.storage.GetQuerier().ListMakerPieRankingsForPie(r.Context(), sqlcgen.ListMakerPieRankingsForPieParams{
		Year:    year,
		Makerid: brand,
		Pieid:   pie,
	})
	if err != nil {
		abortErr(w, 500, err)
		return
	}

	reviewPath := rankingReviewPath("brands/"+url.PathEscape(brand)+"/"+url.PathEscape(pie), year)
	rows := make([]adminRankingRow, 0, len(rankings))
	for _, r := range rankings {
		rows = append(rows, adminRankingRow{
			Userid:      r.Userid,
			Pastry:      r.Pastry,
			Filling:     r.Filling,
			Topping:     r.Topping,
			Looks:       r.Looks,
			Value:       r.Value,
			Notes:       r.Notes.String,
			LastUpdated: formatRankingTime(r.LastUpdated.Time),
			DeleteURL:   reviewPath + "/users/" + url.PathEscape(r.Userid) + "?year=" + strconv.Itoa(int(year)),
		})
	}

	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		abortErr(w, 500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
			IsAdminPage: true,
			Title:       pieYearly.DisplayName + " :: Rankings",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"PieName":    pieYearly.DisplayName,
			"PieSub":     brand + " / " + pie,
			"BackURL":    "/admin/rankings?year=" + strconv.Itoa(int(year)),
			"Rankings":   rows,
			"PageErrors": []string{},
		},
	}
	wrpr.renderPage(w, http.StatusOK, "admin:rankreview", vals)
}

// DeleteAdminBrandPieRanking implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminBrandPieRanking(w http.ResponseWriter, r *http.Request, brand string, pie string, userid string, params generated.DeleteAdminBrandPieRankingParams) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	year, ok := rankingYear(wrpr, w, r, params.Year)
	if !ok {
		return
	}

	if err := wrpr.storage.GetQuerier().DeleteMakerPieRanking(r.Context(), sqlcgen.DeleteMakerPieRankingParams{
		Year:    year,
		Makerid: brand,
		Pieid:   pie,
		Userid:  userid,
	}); err != nil {
		abortErr(w, 500, err)
		return
	}

	adminId := userIDFrom(r)
	log.Info().Str("pie", brand+"/"+pie).Str("rated_by", userid).Str("by", adminId).Msg("Deleted brand pie ranking")
	http.Redirect(w, r, rankingReviewPath("brands/"+url.PathEscape(brand)+"/"+url.PathEscape(pie), year), http.StatusSeeOther)
}

// GetAdminUserPieRankings implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminUserPieRankings(w http.ResponseWriter, r *http.Request, id string, params generated.GetAdminUserPieRankingsParams) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	year, ok := rankingYear(wrpr, w, r, params.Year)
	if !ok {
		return
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	if !piesEnabled {
		abortStatus(w, http.StatusNotFound)
		return
	}

	pieWithOwner, err := wrpr.storage.GetUserPie(r.Context(), year, id)
	if err != nil {
		abortErr(w, 500, err)
		return
	}
	if pieWithOwner == nil || pieWithOwner.Pie.Id == "" {
		abortStatus(w, http.StatusNotFound)
		return
	}

	rankings, err := wrpr.storage.GetQuerier().ListUserPieRankingsForPie(r.Context(), sqlcgen.ListUserPieRankingsForPieParams{
		Year:  year,
		Pieid: id,
	})
	if err != nil {
		abortErr(w, 500, err)
		return
	}

	reviewPath := rankingReviewPath("userpies/"+url.PathEscape(id), year)
	rows := make([]adminRankingRow, 0, len(rankings))
	for _, r := range rankings {
		rows = append(rows, adminRankingRow{
			Userid:      r.Userid,
			Pastry:      r.Pastry.Int32,
			Filling:     r.Filling.Int32,
			Topping:     r.Topping.Int32,
			Looks:       r.Looks.Int32,
			Value:       r.Value.Int32,
			Notes:       r.Notes.String,
			LastUpdated: formatRankingTime(r.LastUpdated.Time),
			DeleteURL:   reviewPath + "/users/" + url.PathEscape(r.Userid) + "?year=" + strconv.Itoa(int(year)),
		})
	}

	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		abortErr(w, 500, err)
		return
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			IsAdminPage: true,
			Title:       pieWithOwner.Pie.DisplayName + " :: Rankings",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"PieName":    pieWithOwner.Pie.DisplayName,
			"PieSub":     id,
			"BackURL":    "/admin/rankings?year=" + strconv.Itoa(int(year)),
			"Rankings":   rows,
			"PageErrors": []string{},
		},
	}
	wrpr.renderPage(w, http.StatusOK, "admin:rankreview", vals)
}

// DeleteAdminUserPieRanking implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminUserPieRanking(w http.ResponseWriter, r *http.Request, id string, userid string, params generated.DeleteAdminUserPieRankingParams) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	year, ok := rankingYear(wrpr, w, r, params.Year)
	if !ok {
		return
	}

	if err := wrpr.storage.GetQuerier().DeleteUserPieRanking(r.Context(), sqlcgen.DeleteUserPieRankingParams{
		Year:   year,
		Pieid:  id,
		Userid: userid,
	}); err != nil {
		abortErr(w, 500, err)
		return
	}

	adminId := userIDFrom(r)
	log.Info().Str("pie", id).Str("rated_by", userid).Str("by", adminId).Msg("Deleted user pie ranking")
	http.Redirect(w, r, rankingReviewPath("userpies/"+url.PathEscape(id), year), http.StatusSeeOther)
}

// rankingReviewPath builds the review-screen redirect target, keeping the
// year the admin was looking at.
func rankingReviewPath(kindPath string, year int32) string {
	return "/admin/rankings/" + kindPath + "?year=" + strconv.Itoa(int(year))
}

// adminRankingRow is one rating on the review screen, pre-flattened so the
// template never has to touch pgtype wrappers or build escaped URLs.
type adminRankingRow struct {
	Userid      string
	Pastry      int32
	Filling     int32
	Topping     int32
	Looks       int32
	Value       int32
	Notes       string
	LastUpdated string
	DeleteURL   string
}

// formatRankingTime keeps zero timestamps (rows never updated) blank rather
// than printing year 1.
func formatRankingTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}
