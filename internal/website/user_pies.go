package website

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

// submitPieForm implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) SubmitPieForm(c *gin.Context, year int64) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
	}

	if !wrpr.config.UserPiesEnabled {
		c.AbortWithStatus(404)
		return
	}

	sin, _ := c.Keys["signedin"].(bool)
	if !sin {
		c.Redirect(http.StatusSeeOther, "/signin")
		return
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("Submit a pie for %d", year),
			Description: "Submit your own mince pie for ranking",
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d, submit", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:        *ay,
				SignedIn:          sin,
				UserPiesEnabled:   wrpr.config.UserPiesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/submit", year)),
			"Readonly":   ro,
			"Makers":     getMakersForYear(c, wrpr.storage, year),
		},
	}

	c.HTML(http.StatusOK, "page:submit", vals)
}

// submitPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) SubmitPie(c *gin.Context, year int64) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
	}

	if !wrpr.config.UserPiesEnabled {
		c.AbortWithStatus(404)
		return
	}

	sin, _ := c.Keys["signedin"].(bool)
	if !sin {
		c.Redirect(http.StatusSeeOther, "/signin")
		return
	}

	if ro {
		c.HTML(http.StatusOK, "page:submit", templater.PageData{
			Head: templater.PageDataHead{
				Title:        "Submit a pie",
				Description:  "Submit your own mince pie for ranking",
				Keywords:     fmt.Sprintf("Mince pies, UK, ranking, %d, submit", year),
				MenuSettings: templater.MenuSettings{ActiveYear: *ay, SignedIn: sin, UserPiesEnabled: wrpr.config.UserPiesEnabled},
			},
			PageData: map[string]interface{}{
				"Year":       year,
				"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/submit", year)),
				"Readonly":   true,
				"Makers":     getMakersForYear(c, wrpr.storage, year),
				"Errors":     []string{"Site is in read-only mode"},
			},
		})
		return
	}

	displayName := c.PostForm("displayname")
	maker := c.PostForm("maker")
	location := c.PostForm("location")
	freshStr := c.PostForm("fresh")
	labelsStr := c.PostForm("labels")
	webLink := c.PostForm("webLink")
	packCountStr := c.PostForm("packCount")
	packPriceStr := c.PostForm("packPrice")

	var errs []string

	if displayName == "" {
		errs = append(errs, "Display name is required")
	}
	if maker == "" {
		errs = append(errs, "Maker/brand is required")
	}
	if webLink == "" {
		errs = append(errs, "Website link is required")
	}

	fresh := freshStr == "on" || freshStr == "true"

	labels := []string{}
	if labelsStr != "" {
		for _, l := range strings.Split(labelsStr, ",") {
			trimmed := strings.TrimSpace(l)
			if trimmed != "" {
				labels = append(labels, trimmed)
			}
		}
	}

	packCount := int32(0)
	if packCountStr != "" {
		pc, err := strconv.ParseInt(packCountStr, 10, 32)
		if err != nil || pc < 0 {
			errs = append(errs, "Pack count must be a non-negative number")
		} else {
			packCount = int32(pc)
		}
	}

	packPriceInPence := int32(0)
	if packPriceStr != "" {
		pp, err := strconv.ParseFloat(packPriceStr, 64)
		if err != nil || pp < 0 {
			errs = append(errs, "Pack price must be a non-negative number (e.g. 3.60)")
		} else {
			packPriceInPence = int32(pp*100 + 0.5) // round to nearest pence
		}
	}

	if len(errs) > 0 {
		c.HTML(http.StatusOK, "page:submit", templater.PageData{
			Head: templater.PageDataHead{
				Title:        "Submit a pie",
				Description:  "Submit your own mince pie for ranking",
				Keywords:     fmt.Sprintf("Mince pies, UK, ranking, %d, submit", year),
				MenuSettings: templater.MenuSettings{ActiveYear: *ay, SignedIn: sin, UserPiesEnabled: wrpr.config.UserPiesEnabled},
			},
			PageData: map[string]interface{}{
				"Year":        year,
				"Breadcrumb":  templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/submit", year)),
				"Readonly":    false,
				"Makers":      getMakersForYear(c, wrpr.storage, year),
				"Errors":      errs,
				"DisplayName": displayName,
				"Maker":       maker,
				"Location":    location,
				"Fresh":       fresh,
				"LabelsStr":   labelsStr,
				"WebLink":     webLink,
				"PackCount":   packCountStr,
				"PackPrice":   packPriceStr,
			},
		})
		return
	}

	// Handle image upload
	imageFile := ""
	imgFile, imgHdr, imgErr := c.Request.FormFile("imageFile")
	if imgErr == nil && imgFile != nil {
		defer imgFile.Close()
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

	if len(errs) > 0 {
		c.HTML(http.StatusOK, "page:submit", templater.PageData{
			Head: templater.PageDataHead{
				Title:        "Submit a pie",
				Description:  "Submit your own mince pie for ranking",
				Keywords:     fmt.Sprintf("Mince pies, UK, ranking, %d, submit", year),
				MenuSettings: templater.MenuSettings{ActiveYear: *ay, SignedIn: sin, UserPiesEnabled: wrpr.config.UserPiesEnabled},
			},
			PageData: map[string]interface{}{
				"Year":        year,
				"Breadcrumb":  templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/submit", year)),
				"Readonly":    false,
				"Makers":      getMakersForYear(c, wrpr.storage, year),
				"Errors":      errs,
				"DisplayName": displayName,
				"Maker":       maker,
				"Location":    location,
				"Fresh":       fresh,
				"LabelsStr":   labelsStr,
				"WebLink":     webLink,
				"PackCount":   packCountStr,
				"PackPrice":   packPriceStr,
			},
		})
		return
	}

	// Create the user pie
	uid := c.Keys["userid"].(string)
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

	_, err = wrpr.storage.CreateUserPie(c, pie)
	if err != nil {
		log.Error().Err(err).Msg("Error creating user pie")
		c.AbortWithStatus(500)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/years/%d/user-pies/%s", year, pie.Id))
}

// listUserPies implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) ListUserPies(c *gin.Context, year int64, params generated.ListUserPiesParams) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	if !wrpr.config.UserPiesEnabled {
		c.AbortWithStatus(404)
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

	pies, err := wrpr.storage.GetAllApprovedUserPies(c, int32(year), int(limit), int(pageZeroIdx*limit))
	if err != nil {
		log.Error().Err(err).Msg("Error getting approved user pies")
		c.AbortWithStatus(500)
		return
	}

	count, err := wrpr.storage.GetUserPieCountByYear(c, int32(year))
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie count")
		c.AbortWithStatus(500)
		return
	}

	pieCards := []templater.UserPieCardData{}
	for _, pie := range pies {
		rs, err := wrpr.storage.GetUserPieRankingSummary(c, pie.Year, pie.Id)
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
			Title:       fmt.Sprintf("Community Pies for %d", year),
			Description: fmt.Sprintf("Community-submitted mince pies for %d", year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d, community", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:        *ay,
				SignedIn:          c.Keys["signedin"].(bool),
				UserPiesEnabled:   wrpr.config.UserPiesEnabled,
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

	c.HTML(http.StatusOK, "page:userpielist", vals)
}

// getUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetUserPie(c *gin.Context, year int64, id string) {
	wrpr.getUserPiePage(c, year, id, nil)
}

// rateUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) RateUserPie(c *gin.Context, year int64, id string) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
	}

	if !wrpr.config.UserPiesEnabled {
		c.AbortWithStatus(404)
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

	pie := &pieWithOwner.Pie
	status := storage.UserPieStatus(pie.Status)
	if !status.IsRankable() {
		wrpr.getUserPiePage(c, year, id, []string{"This pie is not yet available for ranking"})
		return
	}

	if pie.Year != int32(*ay) || ro {
		wrpr.getUserPiePage(c, year, id, []string{"Rating is not open for this year or is read-only"})
		return
	}

	rnkingmap, validationErrs := parseRatings(c.Request.FormValue)
	if len(validationErrs) > 0 {
		wrpr.getUserPiePage(c, year, id, validationErrs)
		return
	}

	uid := c.Keys["userid"].(string)
	wrpr.storage.GetQuerier().UpsertUserPieRanking(c, sqlcgen.UpsertUserPieRankingParams{
		Year:    pie.Year,
		Pieid:   pie.Id,
		Userid:  uid,
		Pastry:  rnkingmap["pastry"],
		Filling: rnkingmap["filling"],
		Topping: rnkingmap["topping"],
		Looks:   rnkingmap["looks"],
		Value:   rnkingmap["value"],
	})

	wrpr.getUserPiePage(c, year, id, nil)
}

// getUserPiePage renders the user pie detail page with optional errors.
func (wrpr *WebsiteWrapper) getUserPiePage(c *gin.Context, year int64, id string, errs []string) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
	}

	if !wrpr.config.UserPiesEnabled {
		c.AbortWithStatus(404)
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

	pie := &pieWithOwner.Pie

	// Check if pie is approved/promoted (publicly visible)
	status := storage.UserPieStatus(pie.Status)
	if !status.IsPublic() && errs == nil {
		c.AbortWithStatus(404)
		return
	}

	rankingSummary, err := wrpr.storage.GetUserPieRankingSummary(c, pie.Year, pie.Id)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie ranking summary")
		c.AbortWithStatus(500)
		return
	}

	uid, _ := c.Keys["userid"].(string)
	sin, _ := c.Keys["signedin"].(bool)

	var userRanking *types.UserPieRankingYearly
	if sin {
		ur, err := wrpr.storage.GetUserPieRanking(c, pie.Year, pie.Id, uid)
		if err == nil && ur != nil {
			userRanking = ur
		}
	}

	isOwner := sin && uid == pie.OwnerUserid

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%s (Community) :: %.2f/5", pie.DisplayName, rankingSummary.Average),
			Description: fmt.Sprintf("%s from %s — community-submitted mince pie scored %.2f/5 in %d", pie.DisplayName, pie.Maker, rankingSummary.Average, pie.Year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d, %s, community", pie.Year, pie.DisplayName),
			MenuSettings: templater.MenuSettings{
				ActiveYear:        *ay,
				SignedIn:          sin,
				UserPiesEnabled:   wrpr.config.UserPiesEnabled,
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
		},
	}

	c.HTML(http.StatusOK, "page:userpie", vals)
}

// getMakersForYear is a helper to fetch makers for the dropdown.
func getMakersForYear(c *gin.Context, storage storage.Operations, year int64) *[]types.Maker {
	makers, err := storage.GetMakersForYear(c, year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting makers for year")
		return nil
	}
	return makers
}
