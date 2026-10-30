package website

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

// brandIDPattern guards brand ids: they end up in URLs and in filesystem
// paths under ImgprssrDir (/brands/{id}/{filename}), so keep them path-safe.
var brandIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// validateBrandForm checks the name and id for a brand create, returning
// the id and any validation errors. The id is entirely user controlled —
// it is used in URLs and filesystem paths, so the user must supply a
// suitable one instead of having one derived from the name.
func validateBrandForm(name, id string) (string, []string) {
	errs := []string{}
	name = strings.TrimSpace(name)
	if name == "" {
		errs = append(errs, "Brand name is required")
	}

	id = strings.TrimSpace(id)
	if id == "" {
		errs = append(errs, "Brand id is required — pick a short lowercase slug such as \"bramwell\"")
	} else if !brandIDPattern.MatchString(id) {
		errs = append(errs, "Brand id must be lowercase letters, numbers and dashes (e.g. \"bramwell\")")
	}
	return id, errs
}

// brandListPage renders the admin brands list.
func brandListPage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, pageErrs ...string) {
	brands, err := wrpr.storage.GetAllMakers(r.Context())
	if err != nil {
		abortErr(w, 500, err)
		return
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
			Title:       "Admin Brands",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Brands":     brands,
			"PageErrors": pageErrs,
		},
	}

	wrpr.renderPage(w, http.StatusOK, "admin:brands", vals)
}

// brandNewPage renders the brand create form, keeping entered values on error.
func brandNewPage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, status int, form map[string]string, pageErrs ...string) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		abortErr(w, 500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
			IsAdminPage: true,
			Title:       "New Brand :: Admin",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Form":       form,
			"PageErrors": pageErrs,
		},
	}

	wrpr.renderPage(w, status, "admin:brandnew", vals)
}

// brandEditPage renders the brand edit form.
func brandEditPage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, status int, maker *types.Maker, pieCount int64, pageErrs ...string) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		abortErr(w, 500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
			IsAdminPage: true,
			Title:       maker.Name + " :: Admin",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Maker":      maker,
			"PieCount":   pieCount,
			"PageErrors": pageErrs,
		},
	}

	wrpr.renderPage(w, status, "admin:brand", vals)
}

// GetAdminBrands implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminBrands(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	brandListPage(wrpr, w, r)
}

// GetAdminBrandNew implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminBrandNew(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	brandNewPage(wrpr, w, r, http.StatusOK, map[string]string{})
}

// CreateAdminBrand implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) CreateAdminBrand(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	website := strings.TrimSpace(r.FormValue("website"))
	id, errs := validateBrandForm(name, r.FormValue("id"))
	if len(errs) > 0 {
		brandNewPage(wrpr, w, r, http.StatusBadRequest, map[string]string{
			"id": r.FormValue("id"), "name": name, "website": website,
		}, errs...)
		return
	}

	// logo upload is optional; stored at /brands/{id}/{filename}
	logo := ""
	if publicPath, uploaded, err := saveUploadedFile(wrpr.config.ImgprssrDir, "logo", r, "brands", id); err != nil {
		abortErr(w, 500, err)
		return
	} else if uploaded {
		logo = publicPath
	}

	_, err := wrpr.storage.CreateMaker(r.Context(), types.Maker{
		Id:      id,
		Name:    name,
		Logo:    logo,
		Website: website,
	})
	if errors.Is(err, storage.ErrBrandAlreadyExists) {
		brandNewPage(wrpr, w, r, http.StatusBadRequest, map[string]string{
			"id": id, "name": name, "website": website,
		}, "A brand with the id \""+id+"\" already exists")
		return
	}
	if err != nil {
		abortErr(w, 500, err)
		return
	}

	userId := userIDFrom(r)
	log.Info().Str("brand", id).Str("by", userId).Msg("Created brand")
	http.Redirect(w, r, "/admin/brands/"+id, http.StatusSeeOther)
}

// countNoun renders a count with a naive-English pluralised noun.
func countNoun(n int64, singular string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %ss", n, singular)
}

// GetAdminBrand implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminBrand(w http.ResponseWriter, r *http.Request, id string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	maker, pieCount, ok := loadBrand(wrpr, w, r, id)
	if !ok {
		return
	}
	brandEditPage(wrpr, w, r, http.StatusOK, maker, pieCount)
}

// UpdateAdminBrand implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) UpdateAdminBrand(w http.ResponseWriter, r *http.Request, id string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	maker, pieCount, ok := loadBrand(wrpr, w, r, id)
	if !ok {
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	website := strings.TrimSpace(r.FormValue("website"))
	if name == "" {
		brandEditPage(wrpr, w, r, http.StatusBadRequest, maker, pieCount, "Brand name is required")
		return
	}

	// optional logo upload: a new file is added under /brands/{id}/ and
	// replaces the recorded path; the old file is removed best-effort.
	if publicPath, uploaded, err := saveUploadedFile(wrpr.config.ImgprssrDir, "logo", r, "brands", id); err != nil {
		abortErr(w, 500, err)
		return
	} else if uploaded {
		if maker.Logo != "" && maker.Logo != publicPath {
			if err := os.Remove(wrpr.config.ImgprssrDir + maker.Logo); err != nil && !os.IsNotExist(err) {
				log.Error().Err(err).Str("path", maker.Logo).Msg("Error removing old brand logo")
			}
		}
		maker.Logo = publicPath
	}

	maker.Name = name
	maker.Website = website
	updated, err := wrpr.storage.UpdateMaker(r.Context(), *maker)
	if errors.Is(err, storage.ErrBrandNotFound) {
		abortStatus(w, http.StatusNotFound)
		return
	}
	if err != nil {
		abortErr(w, 500, err)
		return
	}

	userId := userIDFrom(r)
	log.Info().Str("brand", id).Str("by", userId).Msg("Updated brand")
	http.Redirect(w, r, "/admin/brands/"+updated.Id, http.StatusSeeOther)
}

// DeleteAdminBrand implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminBrand(w http.ResponseWriter, r *http.Request, id string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	maker, pieCount, ok := loadBrand(wrpr, w, r, id)
	if !ok {
		return
	}

	// Guard: deleting a referenced brand silently orphans its pies, because
	// every brand/pie page inner joins maker.
	if pieCount > 0 {
		brandEditPage(wrpr, w, r, http.StatusBadRequest, maker, pieCount,
			"Cannot delete: "+countNoun(pieCount, "pie")+" still reference this brand. Reassign or delete them first.")
		return
	}

	if err := wrpr.storage.DeleteMaker(r.Context(), id); err != nil {
		abortErr(w, 500, err)
		return
	}
	if maker.Logo != "" {
		if err := os.Remove(wrpr.config.ImgprssrDir + maker.Logo); err != nil && !os.IsNotExist(err) {
			log.Error().Err(err).Str("path", maker.Logo).Msg("Error removing brand logo file")
		}
	}

	userId := userIDFrom(r)
	log.Info().Str("brand", id).Str("by", userId).Msg("Deleted brand")
	http.Redirect(w, r, "/admin/brands", http.StatusSeeOther)
}

// loadBrand fetches a brand plus its pie count, aborting with 404 when absent.
func loadBrand(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, id string) (*types.Maker, int64, bool) {
	maker, err := wrpr.storage.GetMaker(r.Context(), id)
	if err != nil {
		abortErr(w, 500, err)
		return nil, 0, false
	}
	if maker == nil {
		abortStatus(w, http.StatusNotFound)
		return nil, 0, false
	}
	pieCount, err := wrpr.storage.GetMakerPieCount(r.Context(), id)
	if err != nil {
		abortErr(w, 500, err)
		return nil, 0, false
	}
	return maker, pieCount, true
}
