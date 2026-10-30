package website

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

// brandIDPattern guards brand ids: they end up in URLs and in filesystem
// paths under ImgprssrDir (/brands/{id}/{filename}), so keep them path-safe.
var brandIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// brandIDFromName derives a brand id from a display name using the same
// convention as the category seed migration: lowercase, spaces to dashes.
func brandIDFromName(name string) string {
	id := strings.ToLower(strings.TrimSpace(name))
	id = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(id, "-")
	id = strings.Trim(id, "-")
	return id
}

// validateBrandForm resolves and checks the id for a brand create, returning
// the id and any validation errors.
func validateBrandForm(name, id string) (string, []string) {
	errs := []string{}
	name = strings.TrimSpace(name)
	if name == "" {
		errs = append(errs, "Brand name is required")
	}

	id = strings.TrimSpace(id)
	if id == "" {
		id = brandIDFromName(name)
		if id == "" {
			errs = append(errs, "Brand id could not be derived from the name — supply one")
		}
	}
	if id != "" && !brandIDPattern.MatchString(id) {
		errs = append(errs, "Brand id must be lowercase letters, numbers and dashes (e.g. \"bramwell\")")
	}
	return id, errs
}

// brandListPage renders the admin brands list.
func brandListPage(wrpr *WebsiteWrapper, c *gin.Context, pageErrs ...string) {
	brands, err := wrpr.storage.GetAllMakers(c)
	if err != nil {
		c.AbortWithError(500, err)
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

	wrpr.html(c, http.StatusOK, "page:admin:brands", vals)
}

// brandNewPage renders the brand create form, keeping entered values on error.
func brandNewPage(wrpr *WebsiteWrapper, c *gin.Context, status int, form map[string]string, pageErrs ...string) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
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

	wrpr.html(c, status, "page:admin:brandnew", vals)
}

// brandEditPage renders the brand edit form.
func brandEditPage(wrpr *WebsiteWrapper, c *gin.Context, status int, maker *types.Maker, pieCount int64, pageErrs ...string) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
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

	wrpr.html(c, status, "page:admin:brand", vals)
}

// GetAdminBrands implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminBrands(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}
	brandListPage(wrpr, c)
}

// GetAdminBrandNew implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminBrandNew(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}
	brandNewPage(wrpr, c, http.StatusOK, map[string]string{})
}

// CreateAdminBrand implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) CreateAdminBrand(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}

	name := strings.TrimSpace(c.Request.FormValue("name"))
	website := strings.TrimSpace(c.Request.FormValue("website"))
	id, errs := validateBrandForm(name, c.Request.FormValue("id"))
	if len(errs) > 0 {
		brandNewPage(wrpr, c, http.StatusBadRequest, map[string]string{
			"id": c.Request.FormValue("id"), "name": name, "website": website,
		}, errs...)
		return
	}

	// logo upload is optional; stored at /brands/{id}/{filename}
	logo := ""
	if publicPath, uploaded, err := saveUploadedFile(wrpr.config.ImgprssrDir, "logo", c, "brands", id); err != nil {
		c.AbortWithError(500, err)
		return
	} else if uploaded {
		logo = publicPath
	}

	_, err := wrpr.storage.CreateMaker(c, types.Maker{
		Id:      id,
		Name:    name,
		Logo:    logo,
		Website: website,
	})
	if errors.Is(err, storage.ErrBrandAlreadyExists) {
		brandNewPage(wrpr, c, http.StatusBadRequest, map[string]string{
			"id": id, "name": name, "website": website,
		}, "A brand with the id \""+id+"\" already exists")
		return
	}
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	userId := c.Keys["userid"].(string)
	log.Info().Str("brand", id).Str("by", userId).Msg("Created brand")
	c.Redirect(http.StatusSeeOther, "/admin/brands/"+id)
}

// countNoun renders a count with a naive-English pluralised noun.
func countNoun(n int64, singular string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %ss", n, singular)
}

// GetAdminBrand implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminBrand(c *gin.Context, id string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	maker, pieCount, ok := loadBrand(wrpr, c, id)
	if !ok {
		return
	}
	brandEditPage(wrpr, c, http.StatusOK, maker, pieCount)
}

// UpdateAdminBrand implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) UpdateAdminBrand(c *gin.Context, id string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	maker, pieCount, ok := loadBrand(wrpr, c, id)
	if !ok {
		return
	}

	name := strings.TrimSpace(c.Request.FormValue("name"))
	website := strings.TrimSpace(c.Request.FormValue("website"))
	if name == "" {
		brandEditPage(wrpr, c, http.StatusBadRequest, maker, pieCount, "Brand name is required")
		return
	}

	// optional logo upload: a new file is added under /brands/{id}/ and
	// replaces the recorded path; the old file is removed best-effort.
	if publicPath, uploaded, err := saveUploadedFile(wrpr.config.ImgprssrDir, "logo", c, "brands", id); err != nil {
		c.AbortWithError(500, err)
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
	updated, err := wrpr.storage.UpdateMaker(c, *maker)
	if errors.Is(err, storage.ErrBrandNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	userId := c.Keys["userid"].(string)
	log.Info().Str("brand", id).Str("by", userId).Msg("Updated brand")
	c.Redirect(http.StatusSeeOther, "/admin/brands/"+updated.Id)
}

// DeleteAdminBrand implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminBrand(c *gin.Context, id string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	maker, pieCount, ok := loadBrand(wrpr, c, id)
	if !ok {
		return
	}

	// Guard: deleting a referenced brand silently orphans its pies, because
	// every brand/pie page inner joins maker.
	if pieCount > 0 {
		brandEditPage(wrpr, c, http.StatusBadRequest, maker, pieCount,
			"Cannot delete: "+countNoun(pieCount, "pie")+" still reference this brand. Reassign or delete them first.")
		return
	}

	if err := wrpr.storage.DeleteMaker(c, id); err != nil {
		c.AbortWithError(500, err)
		return
	}
	if maker.Logo != "" {
		if err := os.Remove(wrpr.config.ImgprssrDir + maker.Logo); err != nil && !os.IsNotExist(err) {
			log.Error().Err(err).Str("path", maker.Logo).Msg("Error removing brand logo file")
		}
	}

	userId := c.Keys["userid"].(string)
	log.Info().Str("brand", id).Str("by", userId).Msg("Deleted brand")
	c.Redirect(http.StatusSeeOther, "/admin/brands")
}

// loadBrand fetches a brand plus its pie count, aborting with 404 when absent.
func loadBrand(wrpr *WebsiteWrapper, c *gin.Context, id string) (*types.Maker, int64, bool) {
	maker, err := wrpr.storage.GetMaker(c, id)
	if err != nil {
		c.AbortWithError(500, err)
		return nil, 0, false
	}
	if maker == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return nil, 0, false
	}
	pieCount, err := wrpr.storage.GetMakerPieCount(c, id)
	if err != nil {
		c.AbortWithError(500, err)
		return nil, 0, false
	}
	return maker, pieCount, true
}
