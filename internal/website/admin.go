package website

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
)

// will abort with a 403 if not admin
func validateAdmin(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request) bool {
	userId := UserIDFrom(r)
	isAdmin, err := wrpr.storage.GetQuerier().IsAdminId(r.Context(), userId)
	if err != nil {
		AbortErr(w, 500, err)
		return false
	}
	if !isAdmin {
		AbortStatus(w, 403)
		return false
	}
	return true
}

// GetAdminIndex implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminIndex(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
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
		pendingCount, _ := wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(ay), "pending")
		approvedCount, _ := wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(ay), "approved")
		promotedCount, _ := wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(ay), "promoted")
		duplicateCount, _ := wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(ay), "duplicate")
		rejectedCount, _ := wrpr.storage.GetUserPieAdminQueueCountByStatus(r.Context(), int32(ay), "rejected")
		if pd, ok := vals.PageData.(map[string]interface{}); ok {
			pd["PendingCount"] = pendingCount
			pd["ApprovedCount"] = approvedCount
			pd["PromotedCount"] = promotedCount
			pd["DuplicateCount"] = duplicateCount
			pd["RejectedCount"] = rejectedCount
		}
	}

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "admin:index", vals)
}

func configPage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	config, err := wrpr.storage.GetQuerier().GetAllConfig(r.Context())
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
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

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "admin:config", vals)
}

// GetAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminConfig(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	configPage(wrpr, w, r)
}

// CreateAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) CreateAdminConfig(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	err := wrpr.storage.GetQuerier().InsertConfig(r.Context(), sqlcgen.InsertConfigParams{
		Key:   r.FormValue("key"),
		Value: r.FormValue("value"),
	})
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	configPage(wrpr, w, r)
}

// DeleteAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminConfig(w http.ResponseWriter, r *http.Request, key string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	err := wrpr.storage.GetQuerier().DeleteConfig(r.Context(), key)
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	configPage(wrpr, w, r)
}

// UpdateAdminConfig implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) UpdateAdminConfig(w http.ResponseWriter, r *http.Request, key string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	err := wrpr.storage.GetQuerier().UpdateConfig(r.Context(), sqlcgen.UpdateConfigParams{
		Key:   key,
		Value: r.FormValue("value"),
	})
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	configPage(wrpr, w, r)
}

func piesPage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, pageErrs ...string) {
	pies, err := wrpr.storage.GetQuerier().GetAllMakerPies(r.Context())
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	brands, err := wrpr.storage.GetAllMakers(r.Context())
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
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

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "admin:pies", vals)
}

// brandExists reports whether a maker row exists for the id, so creates and
// promotions cannot leave pies pointing at a brand that has no page.
func (wrpr *WebsiteWrapper) brandExists(ctx context.Context, makerid string) (bool, error) {
	if makerid == "" {
		return false, nil
	}
	mkr, err := wrpr.storage.GetMaker(ctx, makerid)
	if err != nil {
		return false, err
	}
	return mkr != nil, nil
}

func piePage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, oid string, pageErrs ...string) {
	pie, err := wrpr.storage.GetQuerier().GetMakerPieByOid(r.Context(), oid)
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	uid, err := uuid.Parse(oid)
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	pcats, err := wrpr.storage.GetMakerPieCategoriesForMakerPieOid(r.Context(), uid)
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	cats, err := wrpr.storage.GetMakerPieCategories(r.Context())
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
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

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "admin:pie", vals)
}

// GetMakerPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetMakerPieAdmin(w http.ResponseWriter, r *http.Request, oid string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	piePage(wrpr, w, r, oid)
}

// GetMakerPiesAdmin implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetMakerPiesAdmin(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	piesPage(wrpr, w, r)
}

// CreateMakerPieAdmin implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) CreateMakerPieAdmin(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	makerid := r.FormValue("makerid")
	if ok, err := wrpr.brandExists(r.Context(), makerid); err != nil {
		AbortErr(w, 500, err)
		return
	} else if !ok {
		piesPage(wrpr, w, r, "Unknown brand id: "+makerid+" — create the brand first (admin brands page)")
		return
	}

	year, yearErr := parseFormYear(r.FormValue("year"))
	packCount, packCountErr := parseFormInt("Pack count", r.FormValue("pack_count"))
	packPrice, packPriceErr := parseFormInt("Pack price in pence", r.FormValue("pack_price_in_pence"))
	fresh, freshErr := parseFormBool("Fresh", r.FormValue("fresh"))
	validated, validatedErr := parseFormBool("Validated", r.FormValue("validated"))
	if errs := formFieldErrors(yearErr, packCountErr, packPriceErr, freshErr, validatedErr); len(errs) > 0 {
		piesPage(wrpr, w, r, errs...)
		return
	}

	path, uploaded, err := saveUploadedFile(wrpr.config.ImgprssrDir, "img", r,
		r.FormValue("year"), makerid, r.FormValue("id"))
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	if !uploaded {
		piesPage(wrpr, w, r, "An image file is required")
		return
	}

	oid, err := wrpr.storage.GetQuerier().CreateMakerPie(r.Context(), sqlcgen.CreateMakerPieParams{
		Year:    year,
		Makerid: makerid,
		ID:      r.FormValue("id"),

		Displayname:      r.FormValue("displayname"),
		Fresh:            fresh,
		ImageFile:        path,
		WebLink:          r.FormValue("web_link"),
		PackCount:        packCount,
		PackPriceInPence: packPrice,
		Validated:        validated,
	})
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	http.Redirect(w, r, "/admin/makerpies/"+oid, 301)
}

// DeleteMakerPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteMakerPie(w http.ResponseWriter, r *http.Request, oid string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	err := wrpr.storage.GetQuerier().DeleteMakerPieByOid(r.Context(), oid)
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	piesPage(wrpr, w, r)
}

// UpdateMakerPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) UpdateMakerPie(w http.ResponseWriter, r *http.Request, oid string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	uid, err := uuid.Parse(oid)
	if err != nil {
		AbortStatus(w, 400)
		return
	}

	// Handle optional image upload
	imageFile, _, err := saveUploadedFile(wrpr.config.ImgprssrDir, "img", r,
		r.FormValue("year"), r.FormValue("makerid"), r.FormValue("id"))
	if err != nil {
		AbortErr(w, 500, err)
		return
	}

	packCount, packCountErr := parseFormInt("Pack count", r.FormValue("pack_count"))
	packPrice, packPriceErr := parseFormInt("Pack price in pence", r.FormValue("pack_price_in_pence"))
	fresh, freshErr := parseFormBool("Fresh", r.FormValue("fresh"))
	validated, validatedErr := parseFormBool("Validated", r.FormValue("validated"))
	if errs := formFieldErrors(packCountErr, packPriceErr, freshErr, validatedErr); len(errs) > 0 {
		piePage(wrpr, w, r, oid, errs...)
		return
	}

	err = wrpr.storage.GetQuerier().UpdateMakerPieByOid(r.Context(), sqlcgen.UpdateMakerPieByOidParams{
		Displayname:      r.FormValue("displayname"),
		Fresh:            fresh,
		WebLink:          r.FormValue("web_link"),
		PackCount:        packCount,
		PackPriceInPence: packPrice,
		Validated:        validated,
		Oid:              oid,
	})
	if err != nil {
		AbortErr(w, 500, err)
		return
	}

	// Only update the image if one was actually uploaded
	if imageFile != "" {
		err = wrpr.storage.GetQuerier().UpdateMakerPieImageByOid(r.Context(), sqlcgen.UpdateMakerPieImageByOidParams{
			ImageFile: imageFile,
			Oid:       oid,
		})
		if err != nil {
			AbortErr(w, 500, err)
			return
		}
	}
	err = wrpr.storage.SetPieCategories(r.Context(), uid, r.Form["categories"])
	if errors.Is(err, storage.ErrUnknownCategories) {
		// Unknown slugs type the transaction back, so the pie keeps its old
		// categories; show the problem on the edit page instead of a 500.
		piePage(wrpr, w, r, oid, "Unknown categories: "+strings.TrimPrefix(err.Error(), "unknown category slugs: "))
		return
	}
	if err != nil {
		AbortErr(w, 500, err)
		return
	}
	piePage(wrpr, w, r, oid)
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
