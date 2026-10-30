package website

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
)

// will abort with a 403 if not admin
func validateAdmin(wrpr *WebsiteWrapper, c *gin.Context) bool {
	userId := c.Keys["userid"].(string)
	isAdmin, err := wrpr.storage.GetQuerier().IsAdminId(c, userId)
	if err != nil {
		c.AbortWithError(500, err)
		return false
	}
	if !isAdmin {
		c.AbortWithStatus(403)
		return false
	}
	return true
}

// GetAdminIndex implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminIndex(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
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
			Title:       "Admin Index",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"PendingCount":    0,
			"ApprovedCount":   0,
			"PromotedCount":   0,
			"DuplicateCount":  0,
			"RejectedCount":   0,
			"UserPiesEnabled": piesEnabled,
		},
	}

	if piesEnabled {
		pendingCount, _ := wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(ay), "pending")
		approvedCount, _ := wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(ay), "approved")
		promotedCount, _ := wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(ay), "promoted")
		duplicateCount, _ := wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(ay), "duplicate")
		rejectedCount, _ := wrpr.storage.GetUserPieAdminQueueCountByStatus(c, int32(ay), "rejected")
		if pd, ok := vals.PageData.(map[string]interface{}); ok {
			pd["PendingCount"] = pendingCount
			pd["ApprovedCount"] = approvedCount
			pd["PromotedCount"] = promotedCount
			pd["DuplicateCount"] = duplicateCount
			pd["RejectedCount"] = rejectedCount
		}
	}

	c.HTML(http.StatusOK, "page:admin:index", vals)
}

func configPage(wrpr *WebsiteWrapper, c *gin.Context) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	config, err := wrpr.storage.GetQuerier().GetAllConfig(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	vals := templater.PageData{
		Head: templater.PageDataHead{
			CsrfToken:   c.GetString("csrf_token"),
			IsAdminPage: true,
			Title:       "Admin Config",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Config": config,
		},
	}

	c.HTML(http.StatusOK, "page:admin:config", vals)
}

// GetAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminConfig(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}

	configPage(wrpr, c)
}

// CreateAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) CreateAdminConfig(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}

	err := wrpr.storage.GetQuerier().InsertConfig(c, sqlcgen.InsertConfigParams{
		Key:   c.Request.FormValue("key"),
		Value: c.Request.FormValue("value"),
	})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	configPage(wrpr, c)
}

// DeleteAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminConfig(c *gin.Context, key string) {
	if !validateAdmin(wrpr, c) {
		return
	}

	err := wrpr.storage.GetQuerier().DeleteConfig(c, key)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	configPage(wrpr, c)
}

// UpdateAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) UpdateAdminConfig(c *gin.Context, key string) {
	if !validateAdmin(wrpr, c) {
		return
	}

	err := wrpr.storage.GetQuerier().UpdateConfig(c, sqlcgen.UpdateConfigParams{
		Key:   key,
		Value: c.Request.FormValue("value"),
	})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	configPage(wrpr, c)
}

func piesPage(wrpr *WebsiteWrapper, c *gin.Context, pageErrs ...string) {
	pies, err := wrpr.storage.GetQuerier().GetAllMakerPies(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	brands, err := wrpr.storage.GetAllMakers(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	vals := templater.PageData{
		Head: templater.PageDataHead{
			CsrfToken:   c.GetString("csrf_token"),
			IsAdminPage: true,
			Title:       "Admin Pies",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Pies":       pies,
			"AllBrands":  brands,
			"PageErrors": pageErrs,
		},
	}

	c.HTML(http.StatusOK, "page:admin:pies", vals)
}

// brandExists reports whether a maker row exists for the id, so creates and
// promotions cannot leave pies pointing at a brand that has no page.
func (wrpr *WebsiteWrapper) brandExists(c *gin.Context, makerid string) (bool, error) {
	if makerid == "" {
		return false, nil
	}
	mkr, err := wrpr.storage.GetMaker(c, makerid)
	if err != nil {
		return false, err
	}
	return mkr != nil, nil
}

func piePage(wrpr *WebsiteWrapper, c *gin.Context, oid string, pageErrs ...string) {
	pie, err := wrpr.storage.GetQuerier().GetMakerPieByOid(c, oid)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	uid, err := uuid.Parse(oid)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	pcats, err := wrpr.storage.GetMakerPieCategoriesForMakerPieOid(c, uid)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	cats, err := wrpr.storage.GetMakerPieCategories(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	vals := templater.PageData{
		Head: templater.PageDataHead{
			CsrfToken:   c.GetString("csrf_token"),
			IsAdminPage: true,
			Title:       "Admin Pie",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Pie":           pie,
			"PieCategories": pcats,
			"AllCategories": cats,
			"PageErrors":    pageErrs,
		},
	}

	c.HTML(http.StatusOK, "page:admin:pie", vals)
}

// GetMakerPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetMakerPieAdmin(c *gin.Context, oid string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	piePage(wrpr, c, oid)
}

// GetMakerPiesAdmin implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetMakerPiesAdmin(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}
	piesPage(wrpr, c)
}

// CreateMakerPieAdmin implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) CreateMakerPieAdmin(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}

	makerid := c.Request.FormValue("makerid")
	if ok, err := wrpr.brandExists(c, makerid); err != nil {
		c.AbortWithError(500, err)
		return
	} else if !ok {
		piesPage(wrpr, c, "Unknown brand id: "+makerid+" — create the brand first (admin brands page)")
		return
	}

	year, yearErr := parseFormYear(c.Request.FormValue("year"))
	packCount, packCountErr := parseFormInt("Pack count", c.Request.FormValue("pack_count"))
	packPrice, packPriceErr := parseFormInt("Pack price in pence", c.Request.FormValue("pack_price_in_pence"))
	fresh, freshErr := parseFormBool("Fresh", c.Request.FormValue("fresh"))
	validated, validatedErr := parseFormBool("Validated", c.Request.FormValue("validated"))
	if errs := formFieldErrors(yearErr, packCountErr, packPriceErr, freshErr, validatedErr); len(errs) > 0 {
		piesPage(wrpr, c, errs...)
		return
	}

	path, uploaded, err := saveUploadedFile(wrpr.config.ImgprssrDir, "img", c,
		c.Request.FormValue("year"), makerid, c.Request.FormValue("id"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if !uploaded {
		piesPage(wrpr, c, "An image file is required")
		return
	}

	oid, err := wrpr.storage.GetQuerier().CreateMakerPie(c, sqlcgen.CreateMakerPieParams{
		Year:    year,
		Makerid: makerid,
		ID:      c.Request.FormValue("id"),

		Displayname:      c.Request.FormValue("displayname"),
		Fresh:            fresh,
		ImageFile:        path,
		WebLink:          c.Request.FormValue("web_link"),
		PackCount:        packCount,
		PackPriceInPence: packPrice,
		Validated:        validated,
	})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Redirect(301, "/admin/makerpies/"+oid)
}

// DeleteMakerPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteMakerPie(c *gin.Context, oid string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	err := wrpr.storage.GetQuerier().DeleteMakerPieByOid(c, oid)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piesPage(wrpr, c)
}

// UpdateMakerPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) UpdateMakerPie(c *gin.Context, oid string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	uid, err := uuid.Parse(oid)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	// Handle optional image upload
	imageFile, _, err := saveUploadedFile(wrpr.config.ImgprssrDir, "img", c,
		c.Request.FormValue("year"), c.Request.FormValue("makerid"), c.Request.FormValue("id"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	packCount, packCountErr := parseFormInt("Pack count", c.Request.FormValue("pack_count"))
	packPrice, packPriceErr := parseFormInt("Pack price in pence", c.Request.FormValue("pack_price_in_pence"))
	fresh, freshErr := parseFormBool("Fresh", c.Request.FormValue("fresh"))
	validated, validatedErr := parseFormBool("Validated", c.Request.FormValue("validated"))
	if errs := formFieldErrors(packCountErr, packPriceErr, freshErr, validatedErr); len(errs) > 0 {
		piePage(wrpr, c, oid, errs...)
		return
	}

	err = wrpr.storage.GetQuerier().UpdateMakerPieByOid(c, sqlcgen.UpdateMakerPieByOidParams{
		Displayname:      c.Request.FormValue("displayname"),
		Fresh:            fresh,
		WebLink:          c.Request.FormValue("web_link"),
		PackCount:        packCount,
		PackPriceInPence: packPrice,
		Validated:        validated,
		Oid:              oid,
	})
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Only update the image if one was actually uploaded
	if imageFile != "" {
		err = wrpr.storage.GetQuerier().UpdateMakerPieImageByOid(c, sqlcgen.UpdateMakerPieImageByOidParams{
			ImageFile: imageFile,
			Oid:       oid,
		})
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
	}
	err = wrpr.storage.SetPieCategories(c, uid, c.Request.Form["categories"])
	if errors.Is(err, storage.ErrUnknownCategories) {
		// Unknown slugs type the transaction back, so the pie keeps its old
		// categories; show the problem on the edit page instead of a 500.
		piePage(wrpr, c, oid, "Unknown categories: "+strings.TrimPrefix(err.Error(), "unknown category slugs: "))
		return
	}
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piePage(wrpr, c, oid)
}

// formFieldErrors drops the empty strings from successful field parses so
// the caller can render just the problems.
func formFieldErrors(errs ...string) []string {
	problems := []string{}
	for _, e := range errs {
		if e != "" {
			problems = append(problems, e)
		}
	}
	return problems
}

// parseFormInt converts a form field to an int32. An empty value is stored
// as 0 (the number inputs may be left blank); anything that is not a whole
// number is rejected with a message naming the field so the admin can
// correct it, instead of being quietly replaced with 0.
func parseFormInt(label, value string) (int32, string) {
	if value == "" {
		return 0, ""
	}
	n, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, label + " must be a whole number (got \"" + value + "\")"
	}
	return int32(n), ""
}

// parseFormYear parses the year field, which unlike the pack fields is
// never allowed to be blank.
func parseFormYear(value string) (int32, string) {
	if strings.TrimSpace(value) == "" {
		return 0, "Year is required"
	}
	return parseFormInt("Year", value)
}

// parseFormBool parses a checkbox field. The templates render every
// checkbox with value="true", so "true" means ticked and blank means
// unticked; anything else is not something a checkbox can send and is
// rejected rather than guessed at.
func parseFormBool(label, value string) (bool, string) {
	switch value {
	case "":
		return false, ""
	case "true":
		return true, ""
	default:
		return false, label + " must be ticked or left unticked (got \"" + value + "\")"
	}
}
