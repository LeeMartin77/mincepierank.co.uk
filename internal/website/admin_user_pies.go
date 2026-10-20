package website

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
	"github.com/rs/zerolog/log"
)

// GetAdminUserPies implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminUserPies(c *gin.Context, params generated.GetAdminUserPiesParams) {
	if !validateAdmin(wrpr, c) {
		return
	}

	year := int64(0)
	if params.Year != nil {
		year = *params.Year
	} else {
		ay, err := wrpr.storage.GetActiveYear(c)
		if err != nil {
			log.Error().Err(err).Msg("Error getting active year")
			c.AbortWithStatus(500)
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
		count, err := wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(year), statusFilter)
		if err != nil {
			log.Error().Err(err).Msg("Error getting queue count by status")
			c.AbortWithStatus(500)
			return
		}
		totalCount = count
	} else {
		pendingCount, _ = wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(year), "pending")
		approvedCount, _ = wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(year), "approved")
		promotedCount, _ = wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(year), "promoted")
		duplicateCount, _ = wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(year), "duplicate")
		rejectedCount, _ = wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(year), "rejected")
		totalCount = pendingCount + approvedCount + promotedCount + duplicateCount + rejectedCount
	}

	rows, err := wrpr.storage.GetUserPieAdminQueue(c, int32(year), int(limit), int(pageZeroIdx))
	if err != nil {
		log.Error().Err(err).Msg("Error getting admin queue")
		c.AbortWithStatus(500)
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

	ay, _ := wrpr.storage.GetActiveYear(c)
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("Admin User Pies - %d", year),
			Description: "Moderate user-submitted pies",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Rows":             rows,
			"Year":             year,
			"Page":             page,
			"Limit":            limit,
			"StatusFilter":     statusFilter,
			"SearchFilter":     searchFilter,
			"TotalCount":       totalCount,
			"TotalPages":       (int(totalCount) + int(limit) - 1) / int(limit),
			"PendingCount":     pendingCount,
			"ApprovedCount":    approvedCount,
			"PromotedCount":    promotedCount,
			"DuplicateCount":   duplicateCount,
			"RejectedCount":    rejectedCount,
			"ShouldLoadNext":   int64(len(rows)) == limit,
			"UserPiesEnabled":  piesEnabled,
		},
	}

	c.HTML(http.StatusOK, "page:admin:userpies", vals)
}

// GetAdminUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminUserPie(c *gin.Context, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, c) {
		return
	}

	pieWithOwner, err := wrpr.storage.GetUserPie(c, int32(year), id)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie")
		c.AbortWithStatus(500)
		return
	}
	if pieWithOwner.Pie.Id == "" {
		c.AbortWithStatus(404)
		return
	}

	// Verify owner matches
	if pieWithOwner.OwnerID != ownerid {
		c.AbortWithStatus(404)
		return
	}

	pie := &pieWithOwner.Pie
	status := storage.UserPieStatus(pie.Status)

	// Get duplicate candidates
	dupes, err := wrpr.storage.GetUserPieDuplicateGroups(
		c, pie.Year, pie.Id, pie.Maker, pie.DisplayName, pie.WebLink)
	if err != nil {
		log.Error().Err(err).Msg("Error getting duplicate groups")
	}

	ay, _ := wrpr.storage.GetActiveYear(c)
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	vals := templater.PageData{
		Head: templater.PageDataHead{
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
			"Status":            string(status),
			"DuplicateOf":       pie.DuplicateOf,
			"PromotedMakerPieO": pie.PromotedMakerPieO,
			"ModerationNote":    pie.ModerationNote,
			"DupeCandidates":    dupeCandidatesToPieList(dupes),
			"Year":              pie.Year,
			"UserPiesEnabled":   piesEnabled,
		},
	}

	c.HTML(http.StatusOK, "page:admin:userpie", vals)
}

// PostAdminUserPieApprove implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) PostAdminUserPieApprove(c *gin.Context, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, c) {
		return
	}

	pieWithOwner, err := wrpr.storage.GetUserPie(c, int32(year), id)
	if err != nil || pieWithOwner.Pie.Id == "" || pieWithOwner.OwnerID != ownerid {
		c.AbortWithStatus(404)
		return
	}

	userId := c.Keys["userid"].(string)
	err = wrpr.storage.SetUserPieStatus(c, int32(year), id, storage.UserPieStatus(pieWithOwner.Pie.Status), storage.UserPieStatusApproved)
	if err != nil {
		log.Error().Err(err).Msg("Error approving user pie")
		c.AbortWithStatus(500)
		return
	}

	log.Info().Str("year", fmt.Sprintf("%d", year)).Str("id", id).Str("moderatedBy", userId).Msg("Approved user pie")
	c.Redirect(302, fmt.Sprintf("/admin/userpies/%d/%s/%s", year, id, ownerid))
}

// PostAdminUserPieReject implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) PostAdminUserPieReject(c *gin.Context, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, c) {
		return
	}

	pieWithOwner, err := wrpr.storage.GetUserPie(c, int32(year), id)
	if err != nil || pieWithOwner.Pie.Id == "" || pieWithOwner.OwnerID != ownerid {
		c.AbortWithStatus(404)
		return
	}

	note := c.Request.FormValue("note")
	userId := c.Keys["userid"].(string)
	err = wrpr.storage.SetUserPieStatusWithNote(c, int32(year), id, storage.UserPieStatus(pieWithOwner.Pie.Status), storage.UserPieStatusRejected, note, userId)
	if err != nil {
		log.Error().Err(err).Msg("Error rejecting user pie")
		c.AbortWithStatus(500)
		return
	}

	log.Info().Str("year", fmt.Sprintf("%d", year)).Str("id", id).Str("note", note).Str("moderatedBy", userId).Msg("Rejected user pie")
	c.Redirect(302, fmt.Sprintf("/admin/userpies/%d/%s/%s", year, id, ownerid))
}

// PostAdminUserPieDuplicate implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) PostAdminUserPieDuplicate(c *gin.Context, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, c) {
		return
	}

	pieWithOwner, err := wrpr.storage.GetUserPie(c, int32(year), id)
	if err != nil || pieWithOwner.Pie.Id == "" || pieWithOwner.OwnerID != ownerid {
		c.AbortWithStatus(404)
		return
	}

	duplicateOf := c.Request.FormValue("duplicateOf")
	if duplicateOf == "" {
		c.AbortWithStatus(400)
		return
	}

	userId := c.Keys["userid"].(string)
	err = wrpr.storage.SetUserPieDuplicateOf(c, int32(year), id, storage.UserPieStatus(pieWithOwner.Pie.Status), duplicateOf, userId)
	if err != nil {
		log.Error().Err(err).Msg("Error marking user pie as duplicate")
		c.AbortWithStatus(500)
		return
	}

	log.Info().Str("year", fmt.Sprintf("%d", year)).Str("id", id).Str("duplicateOf", duplicateOf).Str("moderatedBy", userId).Msg("Marked user pie as duplicate")
	c.Redirect(302, fmt.Sprintf("/admin/userpies/%d/%s/%s", year, id, ownerid))
}

// PostAdminUserPiePromote implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) PostAdminUserPiePromote(c *gin.Context, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, c) {
		return
	}

	pieWithOwner, err := wrpr.storage.GetUserPie(c, int32(year), id)
	if err != nil || pieWithOwner.Pie.Id == "" || pieWithOwner.OwnerID != ownerid {
		c.AbortWithStatus(404)
		return
	}

	pie := &pieWithOwner.Pie

	// Parse form values
	displayname := c.Request.FormValue("displayname")
	if displayname == "" {
		displayname = pie.DisplayName
	}
	makerid := c.Request.FormValue("makerid")
	if makerid == "" {
		makerid = pie.Maker
	}
	pieid := c.Request.FormValue("pieid")
	if pieid == "" {
		pieid = pie.Id
	}
	webLink := c.Request.FormValue("webLink")
	if webLink == "" {
		webLink = pie.WebLink
	}
	fresh := mustParseBool(c.Request.FormValue("fresh"))
	validated := mustParseBool(c.Request.FormValue("validated"))
	packCount := mustParseInt(c.Request.FormValue("packCount"))
	packPrice := mustParseInt(c.Request.FormValue("packPriceInPence"))

	// Parse categories
	var categories []string
	cats := c.Request.Form["categories"]
	for _, cat := range cats {
		if cat != "" {
			categories = append(categories, cat)
		}
	}

	// Parse labels
	var labels []string
	labelsStr := c.Request.FormValue("labels")
	if labelsStr != "" {
		for _, l := range strings.Split(labelsStr, ",") {
			l = strings.TrimSpace(l)
			if l != "" {
				labels = append(labels, l)
			}
		}
	}

	// Create brand pie
	newOid, err := wrpr.storage.GetQuerier().CreateMakerPie(c, sqlcgen.CreateMakerPieParams{
		Year:    pie.Year,
		Makerid: makerid,
		ID:      pieid,

		Displayname:      displayname,
		Fresh:            fresh,
		ImageFile:        pie.ImageFile,
		WebLink:          webLink,
		PackCount:        packCount,
		PackPriceInPence: packPrice,
		Validated:        validated,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error creating brand pie from user pie")
		c.AbortWithStatus(500)
		return
	}

	// Link categories to new brand pie
	if len(categories) > 0 {
		parsedOid, _ := uuid.Parse(newOid)
		wrpr.storage.SetPieCategories(c, parsedOid, categories)
	}

	// Update labels (labels are a text[] column)
	if len(labels) > 0 {
		parsedOid, _ := uuid.Parse(newOid)
		_ = wrpr.storage.UpdateMakerPieLabels(c, parsedOid, labels)
	}

	// Mark user pie as promoted
	userId := c.Keys["userid"].(string)
	parsedNewOid, _ := uuid.Parse(newOid)
	err = wrpr.storage.SetUserPiePromotedMakerPieOid(c, int32(year), id, storage.UserPieStatus(pie.Status), parsedNewOid.String(), userId)
	if err != nil {
		log.Error().Err(err).Msg("Error marking user pie as promoted")
	}

	log.Info().Str("userYear", fmt.Sprintf("%d", year)).Str("userId", id).Str("brandOid", newOid).Str("moderatedBy", userId).Msg("Promoted user pie to brand pie")
	c.Redirect(302, "/admin/makerpies/"+newOid)
}

// GetAdminUserPieDuplicates implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminUserPieDuplicates(c *gin.Context, params generated.GetAdminUserPieDuplicatesParams) {
	if !validateAdmin(wrpr, c) {
		return
	}

	year := int64(0)
	if params.Year != nil {
		year = *params.Year
	} else {
		ay, err := wrpr.storage.GetActiveYear(c)
		if err != nil {
			log.Error().Err(err).Msg("Error getting active year")
			c.AbortWithStatus(500)
			return
		}
		year = *ay
	}

	groups, err := wrpr.storage.GetUserPieDuplicateGroupsAll(c, int32(year))
	if err != nil {
		log.Error().Err(err).Msg("Error getting duplicate groups")
		c.AbortWithStatus(500)
		return
	}

	ay, _ := wrpr.storage.GetActiveYear(c)
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	vals := templater.PageData{
		Head: templater.PageDataHead{
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

	c.HTML(http.StatusOK, "page:admin:userpie_duplicates", vals)
}

// DeleteAdminUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminUserPie(c *gin.Context, year int, id string, ownerid string) {
	if !validateAdmin(wrpr, c) {
		return
	}

	pieWithOwner, err := wrpr.storage.GetUserPie(c, int32(year), id)
	if err != nil || pieWithOwner.Pie.Id == "" || pieWithOwner.OwnerID != ownerid {
		c.AbortWithStatus(404)
		return
	}

	pie := pieWithOwner.Pie
	status := storage.UserPieStatus(pie.Status)

	// If promoted, also delete the brand pie
	if status == storage.UserPieStatusPromoted && pie.PromotedMakerPieO != "" {
		delErr := wrpr.storage.GetQuerier().DeleteMakerPieByOid(c, pie.PromotedMakerPieO)
		if delErr != nil {
			log.Error().Err(delErr).Str("brandId", pie.PromotedMakerPieO).Msg("Error deleting brand pie during promotion removal")
		} else {
			log.Info().Str("brandId", pie.PromotedMakerPieO).Msg("Deleted brand pie during promotion removal")
		}
	}

	// Delete the user pie
	err = wrpr.storage.DeleteUserPie(c, int32(year), id)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting user pie")
		c.AbortWithStatus(500)
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
	c.Redirect(302, "/admin/userpies")
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
