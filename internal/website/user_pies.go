package website

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
	"github.com/rs/zerolog/log"
)

// submitPieForm implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) SubmitPieForm(w http.ResponseWriter, r *http.Request, year int64) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	if !piesEnabled {
		abortStatus(w, 404)
		return
	}

	sin := signedInFrom(r)
	if !sin {
		requireSignIn(w, r)
		return
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       fmt.Sprintf("Submit a pie for %d", year),
			Description: "Submit your own mince pie for ranking",
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d, submit", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        sin,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/submit", year)),
			"Readonly":   ro,
			"Makers":     getMakersForYear(w, r, wrpr.storage, year),
		},
	}

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "submit", vals)
}

// submitPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) SubmitPie(w http.ResponseWriter, r *http.Request, year int64) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	if !piesEnabled {
		abortStatus(w, 404)
		return
	}

	sin := signedInFrom(r)
	if !sin {
		requireSignIn(w, r)
		return
	}

	if ro {
		wrpr.htmlTemplater.WritePage(w, http.StatusOK, "submit", templater.PageData{
			Head: templater.PageDataHead{
				Snow:         wrpr.snowing(r.Context()),
				Title:        "Submit a pie",
				Description:  "Submit your own mince pie for ranking",
				Keywords:     fmt.Sprintf("Mince pies, UK, ranking, %d, submit", year),
				MenuSettings: templater.MenuSettings{ActiveYear: *ay, SignedIn: sin, UserPiesEnabled: piesEnabled},
			},
			PageData: map[string]interface{}{
				"Year":       year,
				"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/submit", year)),
				"Readonly":   true,
				"Makers":     getMakersForYear(w, r, wrpr.storage, year),
				"Errors":     []string{"Site is in read-only mode"},
			},
		})
		return
	}

	displayName := r.PostFormValue("displayname")
	maker := r.PostFormValue("maker")
	location := r.PostFormValue("location")
	freshStr := r.PostFormValue("fresh")
	webLink := r.PostFormValue("webLink")
	packCountStr := r.PostFormValue("packCount")
	packPriceStr := r.PostFormValue("packPrice")

	errs, labels := validateUserPieForm(r.FormValue)

	fresh := isFresh(freshStr)

	packCount, pcErr := parsePackCount(packCountStr)
	if pcErr != nil {
		errs = append(errs, pcErr.Error())
	}

	packPriceInPence, ppErr := parsePackPrice(packPriceStr)
	if ppErr != nil {
		errs = append(errs, ppErr.Error())
	}

	if len(errs) > 0 {
		wrpr.htmlTemplater.WritePage(w, http.StatusOK, "submit", templater.PageData{
			Head: templater.PageDataHead{
				Snow:         wrpr.snowing(r.Context()),
				Title:        "Submit a pie",
				Description:  "Submit your own mince pie for ranking",
				Keywords:     fmt.Sprintf("Mince pies, UK, ranking, %d, submit", year),
				MenuSettings: templater.MenuSettings{ActiveYear: *ay, SignedIn: sin, UserPiesEnabled: piesEnabled},
			},
			PageData: map[string]interface{}{
				"Year":        year,
				"Breadcrumb":  templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/submit", year)),
				"Readonly":    false,
				"Makers":      getMakersForYear(w, r, wrpr.storage, year),
				"Errors":      errs,
				"DisplayName": displayName,
				"Maker":       maker,
				"Location":    location,
				"Fresh":       fresh,
				"LabelsStr":   r.PostFormValue("labels"),
				"WebLink":     webLink,
				"PackCount":   packCountStr,
				"PackPrice":   packPriceStr,
			},
		})
		return
	}

	// Handle image upload (max 500KB)
	const maxImageSize = 500 * 1024
	imageFile := ""
	imgFile, imgHdr, imgErr := r.FormFile("imageFile")
	if imgErr == nil && imgFile != nil {
		defer imgFile.Close()
		if imgHdr.Size > maxImageSize {
			errs = append(errs, "Image must be 500KB or smaller")
		} else {
			ext := strings.ToLower(imgHdr.Filename[strings.LastIndex(imgHdr.Filename, "."):])
			newFilename := uuid.NewString() + ext
			destPath := fmt.Sprintf("%s/%s", wrpr.config.ImgprssrDir, newFilename)
			outFile, err := os.Create(destPath)
			if err != nil {
				errs = append(errs, "Failed to save image")
			} else {
				defer outFile.Close()
				_, err = outFile.ReadFrom(imgFile)
				if err != nil {
					errs = append(errs, "Failed to save image")
				} else {
					imageFile = newFilename
				}
			}
		}
	}

	if len(errs) > 0 {
		wrpr.htmlTemplater.WritePage(w, http.StatusOK, "submit", templater.PageData{
			Head: templater.PageDataHead{
				Snow:         wrpr.snowing(r.Context()),
				Title:        "Submit a pie",
				Description:  "Submit your own mince pie for ranking",
				Keywords:     fmt.Sprintf("Mince pies, UK, ranking, %d, submit", year),
				MenuSettings: templater.MenuSettings{ActiveYear: *ay, SignedIn: sin, UserPiesEnabled: piesEnabled},
			},
			PageData: map[string]interface{}{
				"Year":        year,
				"Breadcrumb":  templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/submit", year)),
				"Readonly":    false,
				"Makers":      getMakersForYear(w, r, wrpr.storage, year),
				"Errors":      errs,
				"DisplayName": displayName,
				"Maker":       maker,
				"Location":    location,
				"Fresh":       fresh,
				"LabelsStr":   r.PostFormValue("labels"),
				"WebLink":     webLink,
				"PackCount":   packCountStr,
				"PackPrice":   packPriceStr,
			},
		})
		return
	}

	// Create the user pie
	uid := userIDFrom(r)
	pie := types.UserPieYearly{
		Year:             int32(year),
		Id:               uuid.NewString(),
		OwnerUserid:      uid,
		DisplayName:      displayName,
		Maker:            maker,
		Location:         location,
		Fresh:            fresh,
		Labels:           labels,
		ImageFile:        imageFile,
		WebLink:          webLink,
		PackCount:        packCount,
		PackPriceInPence: packPriceInPence,
		Clean:            true,
		Status:           string(storage.UserPieStatusPending),
	}

	_, err = wrpr.storage.CreateUserPie(r.Context(), pie)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user pie")
		abortStatus(w, 500)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/years/%d/user-pies", year), http.StatusSeeOther)
}

// listUserPies implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) ListUserPies(w http.ResponseWriter, r *http.Request, year int64, params generated.ListUserPiesParams) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	if !piesEnabled {
		abortStatus(w, 404)
		return
	}

	limit := int64(20)
	pageZeroIdx := int64(0)

	if params.Limit != nil {
		l := int64(*params.Limit)
		if l > 0 && l < 101 {
			limit = l
		}
	}
	if params.Page != nil {
		p := int64(*params.Page)
		if p > 0 {
			pageZeroIdx = p - 1
		}
	}

	pies, err := wrpr.storage.GetAllApprovedUserPies(r.Context(), int32(year), int(limit), int(pageZeroIdx*limit))
	if err != nil {
		log.Error().Err(err).Msg("Error getting approved user pies")
		abortStatus(w, 500)
		return
	}

	count, err := wrpr.storage.GetUserPieCountByYear(r.Context(), int32(year))
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie count")
		abortStatus(w, 500)
		return
	}

	pieCards := []templater.UserPieCardData{}
	for _, pie := range pies {
		rs, err := wrpr.storage.GetUserPieRankingSummary(r.Context(), pie.Year, pie.Id)
		if err != nil {
			log.Error().Err(err).Msg("Error getting user pie ranking summary")
			rs = &types.RankingSummary{}
		}
		pieCards = append(pieCards, templater.UserPieCardData{
			PieLink:        fmt.Sprintf("/years/%d/user-pies/%s", pie.Year, pie.Id),
			Pie:            pie,
			RankingSummary: *rs,
		})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       fmt.Sprintf("Community Pies for %d", year),
			Description: fmt.Sprintf("Community-submitted mince pies for %d", year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d, community", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        signedInFrom(r),
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Heading":            fmt.Sprintf("Community Pies for %d", year),
			"Breadcrumb":         templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/user-pies", year)),
			"PieCards":           pieCards,
			"NextPage":           pageZeroIdx + 2,
			"TotalPages":         (int(count) + int(limit) - 1) / int(limit),
			"ShouldLoadNextPage": int64(len(pieCards)) == limit,
			"URL":                fmt.Sprintf("/years/%d/user-pies", year),
			"Query":              "",
		},
	}

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "userpielist", vals)
}

// getUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetUserPie(w http.ResponseWriter, r *http.Request, year int64, id string) {
	wrpr.getUserPiePage(w, r, year, id, nil)
}

// rateUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) RateUserPie(w http.ResponseWriter, r *http.Request, year int64, id string) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	if !piesEnabled {
		abortStatus(w, 404)
		return
	}

	pieWithOwner, err := wrpr.storage.GetUserPie(r.Context(), int32(year), id)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie")
		abortStatus(w, 500)
		return
	}
	if pieWithOwner.Pie.Id == "" {
		abortStatus(w, 404)
		return
	}

	pie := &pieWithOwner.Pie
	status := storage.UserPieStatus(pie.Status)
	if !status.IsRankable() {
		wrpr.getUserPiePage(w, r, year, id, []string{"This pie is not yet available for ranking"})
		return
	}

	if pie.Year != int32(*ay) || ro {
		wrpr.getUserPiePage(w, r, year, id, []string{"Rating is not open for this year or is read-only"})
		return
	}

	rnkingmap, validationErrs := parseRatings(r.FormValue)
	if len(validationErrs) > 0 {
		wrpr.getUserPiePage(w, r, year, id, validationErrs)
		return
	}

	uid := userIDFrom(r)
	wrpr.storage.GetQuerier().UpsertUserPieRanking(r.Context(), sqlcgen.UpsertUserPieRankingParams{
		Year:    pie.Year,
		Pieid:   pie.Id,
		Userid:  uid,
		Pastry:  rnkingmap["pastry"],
		Filling: rnkingmap["filling"],
		Topping: rnkingmap["topping"],
		Looks:   rnkingmap["looks"],
		Value:   rnkingmap["value"],
	})

	wrpr.getUserPiePage(w, r, year, id, nil)
}

// getUserPiePage renders the user pie detail page with optional errors.
func (wrpr *WebsiteWrapper) getUserPiePage(w http.ResponseWriter, r *http.Request, year int64, id string, errs []string) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	if !piesEnabled {
		abortStatus(w, 404)
		return
	}

	pieWithOwner, err := wrpr.storage.GetUserPie(r.Context(), int32(year), id)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie")
		abortStatus(w, 500)
		return
	}
	if pieWithOwner.Pie.Id == "" {
		abortStatus(w, 404)
		return
	}

	pie := &pieWithOwner.Pie

	// Check if pie is approved/promoted (publicly visible)
	status := storage.UserPieStatus(pie.Status)
	if !status.IsPublic() && errs == nil {
		abortStatus(w, 404)
		return
	}

	rankingSummary, err := wrpr.storage.GetUserPieRankingSummary(r.Context(), pie.Year, pie.Id)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie ranking summary")
		abortStatus(w, 500)
		return
	}

	uid := userIDFrom(r)
	sin := signedInFrom(r)

	var userRanking *types.UserPieRankingYearly
	if sin {
		ur, err := wrpr.storage.GetUserPieRanking(r.Context(), pie.Year, pie.Id, uid)
		if err == nil && ur != nil {
			userRanking = ur
		}
	}

	isOwner := sin && uid == pie.OwnerUserid

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       fmt.Sprintf("%s (Community) :: %.2f/5", pie.DisplayName, rankingSummary.Average),
			Description: fmt.Sprintf("%s from %s — community-submitted mince pie scored %.2f/5 in %d", pie.DisplayName, pie.Maker, rankingSummary.Average, pie.Year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d, %s, community", pie.Year, pie.DisplayName),
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        sin,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"ImgprssrPrefix": wrpr.config.ImgprssrPrefix,
			"Pie":            pie,
			"RankingSummary": rankingSummary,
			"Breadcrumb":     templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/user-pies/%s", pie.Year, pie.Id)),
			"Readonly":       ro,
			"UserRanking":    userRanking,
			"IsOwner":        isOwner,
			"Status":         string(status),
			"Errors":         errs,
			"RankingFormURL": fmt.Sprintf("/years/%d/user-pies/%s/rating", pie.Year, pie.Id),
			"SignedIn":       sin,
			"SignInRedirect": fmt.Sprintf("/years/%d/user-pies/%s", pie.Year, pie.Id),
		},
	}

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "userpie", vals)
}

// getMakersForYear is a helper to fetch makers for the dropdown.
func getMakersForYear(w http.ResponseWriter, r *http.Request, storage storage.Operations, year int64) *[]types.Maker {
	makers, err := storage.GetMakersForYear(r.Context(), year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting makers for year")
		return nil
	}
	return makers
}

// validateUserPieForm validates the core required fields and parses labels.
// form is a getter function (e.g. r.FormValue or url.Values.Get).
// Returns (errors, labels).
func validateUserPieForm(form func(string) string) ([]string, []string) {
	var errs []string

	displayName := form("displayname")
	maker := form("maker")
	webLink := form("webLink")

	if displayName == "" {
		errs = append(errs, "Display name is required")
	}
	if maker == "" {
		errs = append(errs, "Maker/brand is required")
	}
	if webLink == "" {
		errs = append(errs, "Website link is required")
	}

	labelsStr := form("labels")
	labels := []string{}
	if labelsStr != "" {
		for _, l := range strings.Split(labelsStr, ",") {
			trimmed := strings.TrimSpace(l)
			if trimmed != "" {
				labels = append(labels, trimmed)
			}
		}
	}

	return errs, labels
}

// parsePackCount parses a pack count string. Returns (value, nil) or (0, error).
func parsePackCount(value string) (int32, error) {
	if value == "" {
		return 0, nil
	}
	pc, err := strconv.ParseInt(value, 10, 32)
	if err != nil || pc < 0 {
		return 0, fmt.Errorf("pack count must be a non-negative number")
	}
	return int32(pc), nil
}

// parsePackPrice parses a pack price string (e.g. "3.60"). Returns pence value or error.
func parsePackPrice(value string) (int32, error) {
	if value == "" {
		return 0, nil
	}
	pp, err := strconv.ParseFloat(value, 64)
	if err != nil || pp < 0 {
		return 0, fmt.Errorf("pack price must be a non-negative number (e.g. 3.60)")
	}
	return int32(pp*100 + 0.5), nil
}

// isFresh returns true if the fresh checkbox was checked.
func isFresh(val string) bool {
	return val == "on" || val == "true"
}
