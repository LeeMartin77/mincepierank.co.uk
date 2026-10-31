package website

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
	"github.com/rs/zerolog/log"
)

// GetAdminUserPies implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminUserPies(w http.ResponseWriter, r *http.Request, params generated.GetAdminUserPiesParams) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	year := int64(0)
	if params.Year != nil {
		year = *params.Year
	} else {
		ay, err := wrpr.storage.GetActiveYear(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Error getting active year")
			abortStatus(w, 500)
			return
		}
		year = *ay
	}

	page := int64(1)
	limit := int64(20)
	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}
	if params.Limit != nil && *params.Limit > 0 && *params.Limit < 101 {
		limit = *params.Limit
	}
	pageZeroIdx := (page - 1) * limit

	statusFilter := ""
	if params.Status != nil {
		statusFilter = *params.Status
	}

	searchFilter := ""
	if params.Search != nil {
		searchFilter = *params.Search
	}

	var totalCount int64
	var pendingCount, approvedCount, promotedCount, duplicateCount, rejectedCount int64

	if statusFilter != "" {
		count, err := wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(year), statusFilter)
		if err != nil {
			log.Error().Err(err).Msg("Error getting queue count by status")
			abortStatus(w, 500)
			return
		}
		totalCount = count
	} else {
		pendingCount, _ = wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(year), "pending")
		approvedCount, _ = wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(year), "approved")
		promotedCount, _ = wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(year), "promoted")
		duplicateCount, _ = wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(year), "duplicate")
		rejectedCount, _ = wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(year), "rejected")
		totalCount = pendingCount + approvedCount + promotedCount + duplicateCount + rejectedCount
	}

	rows, err := wrpr.storage.GetUserPieAdminQueue(r.Context(), int32(year), int(limit), int(pageZeroIdx))
	if err != nil {
		log.Error().Err(err).Msg("Error getting admin queue")
		abortStatus(w, 500)
		return
	}

	// Apply search filter
	if searchFilter != "" {
		s := strings.ToLower(searchFilter)
		var filtered []storage.AdminUserPieRow
		for _, row := range rows {
			displayName := ""
			if row.DisplayName != nil {
				displayName = *row.DisplayName
			}
			maker := ""
			if row.Maker != nil {
				maker = *row.Maker
			}
			if strings.Contains(strings.ToLower(displayName), s) ||
				strings.Contains(strings.ToLower(maker), s) ||
				strings.Contains(row.ID, s) ||
				strings.Contains(strings.ToLower(row.OwnerUserid), s) {
				filtered = append(filtered, row)
			}
		}
		rows = filtered
	}

	// Apply status filter
	if statusFilter != "" {
		var filtered []storage.AdminUserPieRow
		for _, row := range rows {
			if row.Status == statusFilter {
				filtered = append(filtered, row)
			}
		}
		rows = filtered
	}

	ay, _ := wrpr.storage.GetActiveYear(r.Context())
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
			IsAdminPage: true,
			Title:       fmt.Sprintf("Admin User Pies - %d", year),
			Description: "Moderate user-submitted pies",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Rows":            rows,
			"Year":            year,
			"Page":            page,
			"Limit":           limit,
			"StatusFilter":    statusFilter,
			"SearchFilter":    searchFilter,
			"TotalCount":      totalCount,
			"TotalPages":      (int(totalCount) + int(limit) - 1) / int(limit),
			"PendingCount":    pendingCount,
			"ApprovedCount":   approvedCount,
			"PromotedCount":   promotedCount,
			"DuplicateCount":  duplicateCount,
			"RejectedCount":   rejectedCount,
			"ShouldLoadNext":  int64(len(rows)) == limit,
			"UserPiesEnabled": piesEnabled,
		},
	}

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "admin:userpies", vals)
}

// GetAdminUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminUserPie(w http.ResponseWriter, r *http.Request, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	pieWithOwner, ok := loadUserPieForOwner(wrpr, w, r, year, id, ownerid)
	if !ok {
		return
	}

	userPieReviewPage(wrpr, w, r, http.StatusOK, pieWithOwner)
}

// loadUserPieForOwner fetches a user pie and verifies the owner, aborting 404
// on a miss.
func loadUserPieForOwner(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, year int, id, ownerid string) (*storage.UserPieWithOwner, bool) {
	pieWithOwner, err := wrpr.storage.GetUserPie(r.Context(), int32(year), id)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie")
		abortStatus(w, 500)
		return nil, false
	}
	if pieWithOwner.Pie.Id == "" || pieWithOwner.OwnerID != ownerid {
		abortStatus(w, 404)
		return nil, false
	}
	return pieWithOwner, true
}

// userPieReviewPage renders the admin review screen for a user pie, with
// optional error messages (e.g. failed promote validation).
func userPieReviewPage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, status int, pieWithOwner *storage.UserPieWithOwner, pageErrs ...string) {
	pie := &pieWithOwner.Pie
	statusStr := pie.Status

	// Get duplicate candidates
	dupes, err := wrpr.storage.GetUserPieDuplicateGroups(
		r.Context(), pie.Year, pie.Id, pie.Maker, pie.DisplayName, pie.WebLink)
	if err != nil {
		log.Error().Err(err).Msg("Error getting duplicate groups")
	}

	brands, err := wrpr.storage.GetAllMakers(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting brands for promote form")
		abortStatus(w, 500)
		return
	}

	ay, _ := wrpr.storage.GetActiveYear(r.Context())
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
			IsAdminPage: true,
			Title:       fmt.Sprintf("Review: %s (Community) :: Admin", pie.DisplayName),
			Description: "Review user-submitted pie",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Pie":               pie,
			"OwnerID":           pieWithOwner.OwnerID,
			"Status":            statusStr,
			"DuplicateOf":       pie.DuplicateOf,
			"PromotedMakerPieO": pie.PromotedMakerPieO,
			"ModerationNote":    pie.ModerationNote,
			"DupeCandidates":    dupeCandidatesToPieList(dupes),
			"Year":              pie.Year,
			"AllBrands":         brands,
			"UserPiesEnabled":   piesEnabled,
			"PageErrors":        pageErrs,
		},
	}

	wrpr.htmlTemplater.WritePage(w, status, "admin:userpie", vals)
}

// PostAdminUserPieApprove implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) PostAdminUserPieApprove(w http.ResponseWriter, r *http.Request, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	pieWithOwner, ok := loadUserPieForOwner(wrpr, w, r, year, id, ownerid)
	if !ok {
		return
	}

	userId := userIDFrom(r)
	err := wrpr.storage.SetUserPieStatus(r.Context(), int32(year), id, storage.UserPieStatus(pieWithOwner.Pie.Status), storage.UserPieStatusApproved)
	if err != nil {
		log.Error().Err(err).Msg("Error approving user pie")
		abortStatus(w, 500)
		return
	}

	log.Info().Str("year", fmt.Sprintf("%d", year)).Str("id", id).Str("moderatedBy", userId).Msg("Approved user pie")
	http.Redirect(w, r, fmt.Sprintf("/admin/userpies/%d/%s/%s", year, id, ownerid), 302)
}

// PostAdminUserPieReject implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) PostAdminUserPieReject(w http.ResponseWriter, r *http.Request, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	pieWithOwner, ok := loadUserPieForOwner(wrpr, w, r, year, id, ownerid)
	if !ok {
		return
	}

	note := r.FormValue("note")
	userId := userIDFrom(r)
	err := wrpr.storage.SetUserPieStatusWithNote(r.Context(), int32(year), id, storage.UserPieStatus(pieWithOwner.Pie.Status), storage.UserPieStatusRejected, note, userId)
	if err != nil {
		log.Error().Err(err).Msg("Error rejecting user pie")
		abortStatus(w, 500)
		return
	}

	log.Info().Str("year", fmt.Sprintf("%d", year)).Str("id", id).Str("note", note).Str("moderatedBy", userId).Msg("Rejected user pie")
	http.Redirect(w, r, fmt.Sprintf("/admin/userpies/%d/%s/%s", year, id, ownerid), 302)
}

// PostAdminUserPieDuplicate implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) PostAdminUserPieDuplicate(w http.ResponseWriter, r *http.Request, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	pieWithOwner, ok := loadUserPieForOwner(wrpr, w, r, year, id, ownerid)
	if !ok {
		return
	}

	duplicateOf := r.FormValue("duplicateOf")
	if duplicateOf == "" {
		abortStatus(w, 400)
		return
	}

	userId := userIDFrom(r)
	err := wrpr.storage.SetUserPieDuplicateOf(r.Context(), int32(year), id, storage.UserPieStatus(pieWithOwner.Pie.Status), duplicateOf, userId)
	if err != nil {
		log.Error().Err(err).Msg("Error marking user pie as duplicate")
		abortStatus(w, 500)
		return
	}

	log.Info().Str("year", fmt.Sprintf("%d", year)).Str("id", id).Str("duplicateOf", duplicateOf).Str("moderatedBy", userId).Msg("Marked user pie as duplicate")
	http.Redirect(w, r, fmt.Sprintf("/admin/userpies/%d/%s/%s", year, id, ownerid), 302)
}

// PostAdminUserPiePromote implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) PostAdminUserPiePromote(w http.ResponseWriter, r *http.Request, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	pieWithOwner, ok := loadUserPieForOwner(wrpr, w, r, year, id, ownerid)
	if !ok {
		return
	}

	pie := &pieWithOwner.Pie

	// Parse form values
	displayname := r.FormValue("displayname")
	if displayname == "" {
		displayname = pie.DisplayName
	}
	// The promote form must pick a real brand; falling back to the free-text
	// maker from the submission would orphan the pie from any brand page.
	makerid := strings.TrimSpace(r.FormValue("makerid"))
	if makerid == "" {
		userPieReviewPage(wrpr, w, r, http.StatusBadRequest, pieWithOwner,
			"Select the brand this pie belongs to before promoting")
		return
	}
	if exists, err := wrpr.brandExists(r.Context(), makerid); err != nil {
		log.Error().Err(err).Msg("Error checking brand existence for promote")
		abortStatus(w, 500)
		return
	} else if !exists {
		userPieReviewPage(wrpr, w, r, http.StatusBadRequest, pieWithOwner,
			"Unknown brand id: "+makerid+" — create the brand first (admin brands page)")
		return
	}
	pieid := r.FormValue("pieid")
	if pieid == "" {
		pieid = pie.Id
	}
	webLink := r.FormValue("webLink")
	if webLink == "" {
		webLink = pie.WebLink
	}
	// The promote form only supplies a brand id: everything else carries
	// over from the user pie being promoted rather than being read from
	// form fields the form does not even contain. The brand pie starts
	// unvalidated because an admin has not checked its details yet.

	// Parse categories
	var categories []string
	cats := r.Form["categories"]
	for _, cat := range cats {
		if cat != "" {
			categories = append(categories, cat)
		}
	}

	// Parse labels
	var labels []string
	labelsStr := r.FormValue("labels")
	if labelsStr != "" {
		for _, l := range strings.Split(labelsStr, ",") {
			l = strings.TrimSpace(l)
			if l != "" {
				labels = append(labels, l)
			}
		}
	}

	// Create brand pie
	newOid, err := wrpr.storage.GetQuerier().CreateMakerPie(r.Context(), sqlcgen.CreateMakerPieParams{
		Year:    pie.Year,
		Makerid: makerid,
		ID:      pieid,

		Displayname:      displayname,
		Fresh:            pie.Fresh,
		ImageFile:        pie.ImageFile,
		WebLink:          webLink,
		PackCount:        pie.PackCount,
		PackPriceInPence: pie.PackPriceInPence,
		Validated:        false,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error creating brand pie from user pie")
		abortStatus(w, 500)
		return
	}

	// Link categories to new brand pie
	if len(categories) > 0 {
		parsedOid, _ := uuid.Parse(newOid)
		wrpr.storage.SetPieCategories(r.Context(), parsedOid, categories)
	}

	// Update labels (labels are a text[] column)
	if len(labels) > 0 {
		parsedOid, _ := uuid.Parse(newOid)
		_ = wrpr.storage.UpdateMakerPieLabels(r.Context(), parsedOid, labels)
	}

	// Mark user pie as promoted
	userId := userIDFrom(r)
	parsedNewOid, _ := uuid.Parse(newOid)
	err = wrpr.storage.SetUserPiePromotedMakerPieOid(r.Context(), int32(year), id, storage.UserPieStatus(pie.Status), parsedNewOid.String(), userId)
	if err != nil {
		log.Error().Err(err).Msg("Error marking user pie as promoted")
	}

	log.Info().Str("userYear", fmt.Sprintf("%d", year)).Str("userId", id).Str("brandOid", newOid).Str("moderatedBy", userId).Msg("Promoted user pie to brand pie")
	http.Redirect(w, r, "/admin/makerpies/"+newOid, 302)
}

// GetAdminUserPieDuplicates implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminUserPieDuplicates(w http.ResponseWriter, r *http.Request, params generated.GetAdminUserPieDuplicatesParams) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	year := int64(0)
	if params.Year != nil {
		year = *params.Year
	} else {
		ay, err := wrpr.storage.GetActiveYear(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("Error getting active year")
			abortStatus(w, 500)
			return
		}
		year = *ay
	}

	groups, err := wrpr.storage.GetUserPieDuplicateGroupsAll(r.Context(), int32(year))
	if err != nil {
		log.Error().Err(err).Msg("Error getting duplicate groups")
		abortStatus(w, 500)
		return
	}

	ay, _ := wrpr.storage.GetActiveYear(r.Context())
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
			IsAdminPage: true,
			Title:       fmt.Sprintf("User Pie Duplicates - %d", year),
			Description: "Review potential duplicate user pies",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Pies":            groups,
			"Year":            year,
			"UserPiesEnabled": piesEnabled,
		},
	}

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "admin:userpie_duplicates", vals)
}

// DeleteAdminUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminUserPie(w http.ResponseWriter, r *http.Request, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	pieWithOwner, ok := loadUserPieForOwner(wrpr, w, r, year, id, ownerid)
	if !ok {
		return
	}

	pie := pieWithOwner.Pie
	status := storage.UserPieStatus(pie.Status)

	// If promoted, also delete the brand pie
	if status == storage.UserPieStatusPromoted && pie.PromotedMakerPieO != "" {
		delErr := wrpr.storage.GetQuerier().DeleteMakerPieByOid(r.Context(), pie.PromotedMakerPieO)
		if delErr != nil {
			log.Error().Err(delErr).Str("brandId", pie.PromotedMakerPieO).Msg("Error deleting brand pie during promotion removal")
		} else {
			log.Info().Str("brandId", pie.PromotedMakerPieO).Msg("Deleted brand pie during promotion removal")
		}
	}

	// Delete the user pie
	err := wrpr.storage.DeleteUserPie(r.Context(), int32(year), id)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting user pie")
		abortStatus(w, 500)
		return
	}

	// Delete the stored image file
	if pie.ImageFile != "" {
		imgPath := fmt.Sprintf("%s/%s", wrpr.config.ImgprssrDir, pie.ImageFile)
		if delErr := os.Remove(imgPath); delErr != nil {
			log.Error().Err(delErr).Str("imagePath", imgPath).Msg("Error deleting user pie image")
		}
	}

	log.Info().Str("year", fmt.Sprintf("%d", year)).Str("id", id).Str("status", string(status)).Msg("Deleted user pie")
	http.Redirect(w, r, "/admin/userpies", 302)
}

// dupeCandidatesToPieList converts duplicate candidate pies to a simple list
// for template rendering.
func dupeCandidatesToPieList(pies []types.UserPieYearly) []struct {
	ID   string
	Name string
} {
	var result []struct {
		ID   string
		Name string
	}
	for _, p := range pies {
		result = append(result, struct {
			ID   string
			Name string
		}{ID: p.Id, Name: p.DisplayName})
	}
	return result
}
