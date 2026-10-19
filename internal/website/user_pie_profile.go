package website

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

// myUserPies implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) MyUserPies(c *gin.Context, year int64) {
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

	sin, _ := c.Keys["signedin"].(bool)
	if !sin {
		c.Redirect(http.StatusSeeOther, "/signin")
		return
	}

	uid := c.Keys["userid"].(string)

	pies, err := wrpr.storage.GetOwnerUserPies(c, int32(year), uid)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pies")
		c.AbortWithStatus(500)
		return
	}

	pieCards := []map[string]interface{}{}
	for _, pie := range pies {
		status := storage.UserPieStatus(pie.Status)
		pieLink := fmt.Sprintf("/years/%d/user-pies/%s", pie.Year, pie.Id)
		if !status.IsPublic() {
			pieLink = fmt.Sprintf("/profile/pies/%d/%s", pie.Year, pie.Id)
		}
		rs, err := wrpr.storage.GetUserPieRankingSummary(c, pie.Year, pie.Id)
		if err != nil {
			log.Error().Err(err).Msg("Error getting user pie ranking summary")
			rs = &types.RankingSummary{}
		}
		pieCards = append(pieCards, map[string]interface{}{
			"ImgprssrPrefix": wrpr.config.ImgprssrPrefix,
			"PieLink":        pieLink,
			"Pie":            pie,
			"RankingSummary": *rs,
			"Status":         string(status),
			"ModerationNote": pie.ModerationNote,
		})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("My Community Pies for %d", year),
			Description: "Manage your community-submitted mince pies",
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d, my pies", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:        *ay,
				SignedIn:          sin,
				UserPiesEnabled:   wrpr.config.UserPiesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Heading":    fmt.Sprintf("My Community Pies for %d", year),
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/profile/pies/%d", year)),
			"PieCards":   pieCards,
			"Year":       year,
		},
	}

	c.HTML(http.StatusOK, "page:profilepies", vals)
}

// editUserPieForm implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) EditUserPieForm(c *gin.Context, year int64, id string) {
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

	sin, _ := c.Keys["signedin"].(bool)
	if !sin {
		c.Redirect(http.StatusSeeOther, "/signin")
		return
	}

	uid := c.Keys["userid"].(string)

	pie, err := wrpr.storage.GetUserPieByOwner(c, int32(year), id, uid)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie")
		c.AbortWithStatus(500)
		return
	}
	if pie.Id == "" {
		c.AbortWithStatus(404)
		return
	}

	// Only editable if pending
	status := storage.UserPieStatus(pie.Status)
	if !status.IsOwnerEditable() {
		// Not editable — redirect to the pie page or my pies
		if status.IsPublic() {
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/years/%d/user-pies/%s", year, id))
		} else {
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/profile/pies/%d", year))
		}
		return
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("Edit: %s", pie.DisplayName),
			Description: "Edit your community-submitted mince pie",
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d, edit", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:        *ay,
				SignedIn:          sin,
				UserPiesEnabled:   wrpr.config.UserPiesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"Pie":        pie,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/profile/pies/%d/%s", year, id)),
			"Makers":     getMakersForYear(c, wrpr.storage, year),
			"LabelsStr":  strings.Join(pie.Labels, ", "),
		},
	}

	c.HTML(http.StatusOK, "page:editpie", vals)
}

// editUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) EditUserPie(c *gin.Context, year int64, id string) {
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

	uid := c.Keys["userid"].(string)

	pie, err := wrpr.storage.GetUserPieByOwner(c, int32(year), id, uid)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie")
		c.AbortWithStatus(500)
		return
	}
	if pie.Id == "" {
		c.AbortWithStatus(404)
		return
	}

	// Only editable if pending
	status := storage.UserPieStatus(pie.Status)
	if !status.IsOwnerEditable() {
		if status.IsPublic() {
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/years/%d/user-pies/%s", year, id))
		} else {
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/profile/pies/%d", year))
		}
		return
	}

	if ro {
		c.HTML(http.StatusOK, "page:editpie", templater.PageData{
			Head: templater.PageDataHead{
				Title:        "Edit pie",
				Description:  "Edit your community-submitted mince pie",
				Keywords:     fmt.Sprintf("Mince pies, UK, ranking, %d, edit", year),
				MenuSettings: templater.MenuSettings{ActiveYear: *ay, SignedIn: sin, UserPiesEnabled: wrpr.config.UserPiesEnabled},
			},
			PageData: map[string]interface{}{
				"Year":       year,
				"Pie":        pie,
				"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/profile/pies/%d/%s", year, id)),
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

	// Handle image upload (optional — keep existing if no new image)
	imageFile := pie.ImageFile
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
		c.HTML(http.StatusOK, "page:editpie", templater.PageData{
			Head: templater.PageDataHead{
				Title:        "Edit pie",
				Description:  "Edit your community-submitted mince pie",
				Keywords:     fmt.Sprintf("Mince pies, UK, ranking, %d, edit", year),
				MenuSettings: templater.MenuSettings{ActiveYear: *ay, SignedIn: sin, UserPiesEnabled: wrpr.config.UserPiesEnabled},
			},
			PageData: map[string]interface{}{
				"Year":        year,
				"Pie":         pie,
				"Breadcrumb":  templater.BreadcrumbsFromUrl(fmt.Sprintf("/profile/pies/%d/%s", year, id)),
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

	// Update the pie
	updatedPie := *pie
	updatedPie.DisplayName = displayName
	updatedPie.Maker = maker
	updatedPie.Location = location
	updatedPie.Fresh = fresh
	updatedPie.Labels = labels
	updatedPie.ImageFile = imageFile
	updatedPie.WebLink = webLink
	updatedPie.PackCount = packCount
	updatedPie.PackPriceInPence = packPriceInPence

	err = wrpr.storage.UpdateUserPie(c, updatedPie)
	if err != nil {
		log.Error().Err(err).Msg("Error updating user pie")
		c.AbortWithStatus(500)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/years/%d/user-pies/%s", year, pie.Id))
}

// withdrawUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) WithdrawUserPie(c *gin.Context, year int64, id string) {
	if !wrpr.config.UserPiesEnabled {
		c.AbortWithStatus(404)
		return
	}

	sin, _ := c.Keys["signedin"].(bool)
	if !sin {
		c.Redirect(http.StatusSeeOther, "/signin")
		return
	}

	uid := c.Keys["userid"].(string)

	pie, err := wrpr.storage.GetUserPieByOwner(c, int32(year), id, uid)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie")
		c.AbortWithStatus(500)
		return
	}
	if pie.Id == "" {
		c.AbortWithStatus(404)
		return
	}

	// Only owner can withdraw, and only if pending
	status := storage.UserPieStatus(pie.Status)
	if !status.IsOwnerEditable() {
		if status.IsPublic() {
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/years/%d/user-pies/%s", year, id))
		} else {
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/profile/pies/%d", year))
		}
		return
	}

	// Delete the pie and its rankings
	wrpr.storage.GetQuerier().DeleteUserPieRankings(c, sqlcgen.DeleteUserPieRankingsParams{
		Year: pie.Year,
		Pieid: pie.Id,
	})
	wrpr.storage.DeleteUserPieByOwner(c, pie.Year, pie.Id, uid)

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/profile/pies/%d", year))
}
