package website

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
	"github.com/rs/zerolog/log"
)

// rankingYear resolves the optional ?year= query parameter, falling back to
// the active year. The generated wrapper already rejects non-numeric years.
func rankingYear(wrpr *WebsiteWrapper, c *gin.Context, paramYear *int) (int32, bool) {
	if paramYear != nil {
		return int32(*paramYear), true
	}
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return 0, false
	}
	return int32(ay), true
}

// GetAdminRankings implements generated.ServerInterface. The picker lists
// every brand pie (and, when user pies are on, every public user pie) for the
// chosen year with a rating count, so admins can drill into a review screen.
func (wrpr *WebsiteWrapper) GetAdminRankings(c *gin.Context, params generated.GetAdminRankingsParams) {
	if !validateAdmin(wrpr, c) {
		return
	}
	year, ok := rankingYear(wrpr, c, params.Year)
	if !ok {
		return
	}

	brandPies, err := wrpr.storage.GetQuerier().ListMakerPieRankingSummaries(c, year)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	userPies := []sqlcgen.ListUserPieRankingSummariesRow{}
	if piesEnabled {
		userPies, err = wrpr.storage.GetQuerier().ListUserPieRankingSummaries(c, year)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
	}

	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			CsrfToken:   c.GetString("csrf_token"),
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
	c.HTML(http.StatusOK, "page:admin:rankings", vals)
}

// GetAdminBrandPieRankings implements generated.ServerInterface. The review
// screen lists each rating one user gave this brand pie, including the notes
// that are collected but never shown publicly.
func (wrpr *WebsiteWrapper) GetAdminBrandPieRankings(c *gin.Context, brand string, pie string, params generated.GetAdminBrandPieRankingsParams) {
	if !validateAdmin(wrpr, c) {
		return
	}
	year, ok := rankingYear(wrpr, c, params.Year)
	if !ok {
		return
	}

	pieYearly, err := wrpr.storage.GetMakerPieYearly(c, year, brand, pie)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if pieYearly == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	rankings, err := wrpr.storage.GetQuerier().ListMakerPieRankingsForPie(c, sqlcgen.ListMakerPieRankingsForPieParams{
		Year:    year,
		Makerid: brand,
		Pieid:   pie,
	})
	if err != nil {
		c.AbortWithError(500, err)
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

	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	vals := templater.PageData{
		Head: templater.PageDataHead{
			CsrfToken:   c.GetString("csrf_token"),
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
	c.HTML(http.StatusOK, "page:admin:rankreview", vals)
}

// DeleteAdminBrandPieRanking implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminBrandPieRanking(c *gin.Context, brand string, pie string, userid string, params generated.DeleteAdminBrandPieRankingParams) {
	if !validateAdmin(wrpr, c) {
		return
	}
	year, ok := rankingYear(wrpr, c, params.Year)
	if !ok {
		return
	}

	if err := wrpr.storage.GetQuerier().DeleteMakerPieRanking(c, sqlcgen.DeleteMakerPieRankingParams{
		Year:    year,
		Makerid: brand,
		Pieid:   pie,
		Userid:  userid,
	}); err != nil {
		c.AbortWithError(500, err)
		return
	}

	adminId := c.Keys["userid"].(string)
	log.Info().Str("pie", brand+"/"+pie).Str("rated_by", userid).Str("by", adminId).Msg("Deleted brand pie ranking")
	c.Redirect(http.StatusSeeOther, rankingReviewPath("brands/"+url.PathEscape(brand)+"/"+url.PathEscape(pie), year))
}

// GetAdminUserPieRankings implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminUserPieRankings(c *gin.Context, id string, params generated.GetAdminUserPieRankingsParams) {
	if !validateAdmin(wrpr, c) {
		return
	}
	year, ok := rankingYear(wrpr, c, params.Year)
	if !ok {
		return
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	if !piesEnabled {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	pieWithOwner, err := wrpr.storage.GetUserPie(c, year, id)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if pieWithOwner == nil || pieWithOwner.Pie.Id == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	rankings, err := wrpr.storage.GetQuerier().ListUserPieRankingsForPie(c, sqlcgen.ListUserPieRankingsForPieParams{
		Year:  year,
		Pieid: id,
	})
	if err != nil {
		c.AbortWithError(500, err)
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

	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			CsrfToken:   c.GetString("csrf_token"),
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
	c.HTML(http.StatusOK, "page:admin:rankreview", vals)
}

// DeleteAdminUserPieRanking implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminUserPieRanking(c *gin.Context, id string, userid string, params generated.DeleteAdminUserPieRankingParams) {
	if !validateAdmin(wrpr, c) {
		return
	}
	year, ok := rankingYear(wrpr, c, params.Year)
	if !ok {
		return
	}

	if err := wrpr.storage.GetQuerier().DeleteUserPieRanking(c, sqlcgen.DeleteUserPieRankingParams{
		Year:   year,
		Pieid:  id,
		Userid: userid,
	}); err != nil {
		c.AbortWithError(500, err)
		return
	}

	adminId := c.Keys["userid"].(string)
	log.Info().Str("pie", id).Str("rated_by", userid).Str("by", adminId).Msg("Deleted user pie ranking")
	c.Redirect(http.StatusSeeOther, rankingReviewPath("userpies/"+url.PathEscape(id), year))
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
