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
	"github.com/rs/zerolog/log"
)

// myUserPies implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) MyUserPies(w http.ResponseWriter, r *http.Request, year int64) {
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

	sin := signedInFrom(r)
	if !sin {
		requireSignIn(w, r)
		return
	}

	uid := userIDFrom(r)

	pies, err := wrpr.storage.GetOwnerUserPies(r.Context(), int32(year), uid)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pies")
		abortStatus(w, 500)
		return
	}

	pieCards := []map[string]interface{}{}
	for _, pie := range pies {
		status := storage.UserPieStatus(pie.Status)
		pieLink := fmt.Sprintf("/years/%d/user-pies/%s", pie.Year, pie.Id)
		if !status.IsPublic() {
			pieLink = fmt.Sprintf("/profile/pies/%d/%s", pie.Year, pie.Id)
		}
		rs, err := wrpr.storage.GetUserPieRankingSummary(r.Context(), pie.Year, pie.Id)
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
			Snow:        wrpr.snowing(r.Context()),
			Title:       fmt.Sprintf("My Community Pies for %d", year),
			Description: "Manage your community-submitted mince pies",
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d, my pies", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        sin,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Heading":    fmt.Sprintf("My Community Pies for %d", year),
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/profile/pies/%d", year)),
			"PieCards":   pieCards,
			"Year":       year,
		},
	}

	wrpr.renderPage(w, http.StatusOK, "profilepies", vals)
}

// editUserPieForm implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) EditUserPieForm(w http.ResponseWriter, r *http.Request, year int64, id string) {
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

	sin := signedInFrom(r)
	if !sin {
		requireSignIn(w, r)
		return
	}

	uid := userIDFrom(r)

	pie, err := wrpr.storage.GetUserPieByOwner(r.Context(), int32(year), id, uid)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie")
		abortStatus(w, 500)
		return
	}
	if pie.Id == "" {
		abortStatus(w, 404)
		return
	}

	// Only editable if pending
	status := storage.UserPieStatus(pie.Status)
	if !status.IsOwnerEditable() {
		// Not editable — redirect to the pie page or my pies
		if status.IsPublic() {
			http.Redirect(w, r, fmt.Sprintf("/years/%d/user-pies/%s", year, id), http.StatusSeeOther)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/profile/pies/%d", year), http.StatusSeeOther)
		}
		return
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       fmt.Sprintf("Edit: %s", pie.DisplayName),
			Description: "Edit your community-submitted mince pie",
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d, edit", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        sin,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"Pie":        pie,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/profile/pies/%d/%s", year, id)),
			"Makers":     getMakersForYear(w, r, wrpr.storage, year),
			"LabelsStr":  strings.Join(pie.Labels, ", "),
		},
	}

	wrpr.renderPage(w, http.StatusOK, "editpie", vals)
}

// editUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) EditUserPie(w http.ResponseWriter, r *http.Request, year int64, id string) {
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

	uid := userIDFrom(r)

	pie, err := wrpr.storage.GetUserPieByOwner(r.Context(), int32(year), id, uid)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie")
		abortStatus(w, 500)
		return
	}
	if pie.Id == "" {
		abortStatus(w, 404)
		return
	}

	// Only editable if pending
	status := storage.UserPieStatus(pie.Status)
	if !status.IsOwnerEditable() {
		if status.IsPublic() {
			http.Redirect(w, r, fmt.Sprintf("/years/%d/user-pies/%s", year, id), http.StatusSeeOther)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/profile/pies/%d", year), http.StatusSeeOther)
		}
		return
	}

	if ro {
		wrpr.renderPage(w, http.StatusOK, "editpie", templater.PageData{
			Head: templater.PageDataHead{
				Snow:         wrpr.snowing(r.Context()),
				Title:        "Edit pie",
				Description:  "Edit your community-submitted mince pie",
				Keywords:     fmt.Sprintf("Mince pies, UK, ranking, %d, edit", year),
				MenuSettings: templater.MenuSettings{ActiveYear: *ay, SignedIn: sin, UserPiesEnabled: piesEnabled},
			},
			PageData: map[string]interface{}{
				"Year":       year,
				"Pie":        pie,
				"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/profile/pies/%d/%s", year, id)),
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

	// Handle image upload (optional — keep existing if no new image)
	const maxImageSize = 500 * 1024
	imageFile := pie.ImageFile
	imgFile, imgHdr, imgErr := r.FormFile("imageFile")
	if imgErr == nil && imgFile != nil {
		defer imgFile.Close()
		if imgHdr.Size > maxImageSize {
			errs = append(errs, "Image must be 500KB or smaller")
		} else {
			// Delete old image if we're replacing it
			if pie.ImageFile != "" {
				oldPath := fmt.Sprintf("%s/%s", wrpr.config.ImgprssrDir, pie.ImageFile)
				_ = os.Remove(oldPath)
			}
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
		wrpr.renderPage(w, http.StatusOK, "editpie", templater.PageData{
			Head: templater.PageDataHead{
				Snow:         wrpr.snowing(r.Context()),
				Title:        "Edit pie",
				Description:  "Edit your community-submitted mince pie",
				Keywords:     fmt.Sprintf("Mince pies, UK, ranking, %d, edit", year),
				MenuSettings: templater.MenuSettings{ActiveYear: *ay, SignedIn: sin, UserPiesEnabled: piesEnabled},
			},
			PageData: map[string]interface{}{
				"Year":        year,
				"Pie":         pie,
				"Breadcrumb":  templater.BreadcrumbsFromUrl(fmt.Sprintf("/profile/pies/%d/%s", year, id)),
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

	err = wrpr.storage.UpdateUserPie(r.Context(), updatedPie)
	if err != nil {
		log.Error().Err(err).Msg("Error updating user pie")
		abortStatus(w, 500)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/years/%d/user-pies/%s", year, pie.Id), http.StatusSeeOther)
}

// withdrawUserPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) WithdrawUserPie(w http.ResponseWriter, r *http.Request, year int64, id string) {
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

	uid := userIDFrom(r)

	pie, err := wrpr.storage.GetUserPieByOwner(r.Context(), int32(year), id, uid)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user pie")
		abortStatus(w, 500)
		return
	}
	if pie.Id == "" {
		abortStatus(w, 404)
		return
	}

	// Only owner can withdraw, and only if pending
	status := storage.UserPieStatus(pie.Status)
	if !status.IsOwnerEditable() {
		if status.IsPublic() {
			http.Redirect(w, r, fmt.Sprintf("/years/%d/user-pies/%s", year, id), http.StatusSeeOther)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/profile/pies/%d", year), http.StatusSeeOther)
		}
		return
	}

	// Delete the pie and its rankings
	wrpr.storage.GetQuerier().DeleteUserPieRankings(r.Context(), sqlcgen.DeleteUserPieRankingsParams{
		Year:  pie.Year,
		Pieid: pie.Id,
	})
	wrpr.storage.DeleteUserPieByOwner(r.Context(), pie.Year, pie.Id, uid)

	// Delete the stored image file
	if pie.ImageFile != "" {
		imgPath := fmt.Sprintf("%s/%s", wrpr.config.ImgprssrDir, pie.ImageFile)
		if delErr := os.Remove(imgPath); delErr != nil {
			log.Error().Err(delErr).Str("imagePath", imgPath).Msg("Error deleting user pie image on withdraw")
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/profile/pies/%d", year), http.StatusSeeOther)
}
