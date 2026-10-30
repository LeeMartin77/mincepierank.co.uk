package website

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

// validateCategoryForm checks the label and derives a slug from it (or takes
// an explicit slug) using the same convention as the brands: lowercase,
// spaces to dashes.
func validateCategoryForm(label, slug string) (string, string, []string) {
	errs := []string{}
	label = strings.TrimSpace(label)
	if label == "" {
		errs = append(errs, "Category label is required")
	}

	slug = strings.TrimSpace(slug)
	if slug == "" {
		slug = brandIDFromName(label)
		if slug == "" {
			errs = append(errs, "Category slug could not be derived from the label — supply one")
		}
	}
	if slug != "" && !brandIDPattern.MatchString(slug) {
		errs = append(errs, "Category slug must be lowercase letters, numbers and dashes (e.g. \"vegetarian\")")
	}
	return label, slug, errs
}

// categoryListPage renders the admin categories list with per-category pie
// usage counts.
func categoryListPage(wrpr *WebsiteWrapper, c *gin.Context, pageErrs ...string) {
	cats, err := wrpr.storage.GetAllCategoriesAdmin(c)
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
			Title:       "Admin Categories",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Categories": cats,
			"PageErrors": pageErrs,
		},
	}

	wrpr.html(c, http.StatusOK, "page:admin:categories", vals)
}

// categoryNewPage renders the category create form, keeping entered values on
// error.
func categoryNewPage(wrpr *WebsiteWrapper, c *gin.Context, status int, form map[string]string, pageErrs ...string) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	vals := templater.PageData{
		Head: templater.PageDataHead{
			IsAdminPage: true,
			Title:       "New Category :: Admin",
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

	wrpr.html(c, status, "page:admin:categorynew", vals)
}

// categoryEditPage renders the category edit form. The slug is shown but not
// editable: pies reference categories by slug, so changing it would break
// bulk imports and saved selections.
func categoryEditPage(wrpr *WebsiteWrapper, c *gin.Context, status int, cat *types.Category, pieCount int64, pageErrs ...string) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(c)
	vals := templater.PageData{
		Head: templater.PageDataHead{
			IsAdminPage: true,
			Title:       cat.Label + " :: Admin",
			Description: "",
			Keywords:    "",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      ay,
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Category":   cat,
			"PieCount":   pieCount,
			"PageErrors": pageErrs,
		},
	}

	wrpr.html(c, status, "page:admin:category", vals)
}

// GetAdminCategories implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminCategories(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}
	categoryListPage(wrpr, c)
}

// GetAdminCategoryNew implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminCategoryNew(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}
	categoryNewPage(wrpr, c, http.StatusOK, map[string]string{})
}

// CreateAdminCategory implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) CreateAdminCategory(c *gin.Context) {
	if !validateAdmin(wrpr, c) {
		return
	}

	label, slug, errs := validateCategoryForm(c.Request.FormValue("label"), c.Request.FormValue("slug"))
	if len(errs) > 0 {
		categoryNewPage(wrpr, c, http.StatusBadRequest, map[string]string{
			"label": label, "slug": c.Request.FormValue("slug"),
		}, errs...)
		return
	}

	created, err := wrpr.storage.CreateCategory(c, types.Category{Slug: slug, Label: label})
	if errors.Is(err, storage.ErrCategoryAlreadyExists) {
		categoryNewPage(wrpr, c, http.StatusBadRequest, map[string]string{
			"label": label, "slug": slug,
		}, "A category with the slug \""+slug+"\" already exists")
		return
	}
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	userId := c.Keys["userid"].(string)
	log.Info().Str("category", created.Slug).Str("by", userId).Msg("Created category")
	c.Redirect(http.StatusSeeOther, "/admin/categories/"+created.Id.String())
}

// GetAdminCategory implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminCategory(c *gin.Context, id string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	cat, pieCount, ok := loadCategory(wrpr, c, id)
	if !ok {
		return
	}
	categoryEditPage(wrpr, c, http.StatusOK, cat, pieCount)
}

// UpdateAdminCategory implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) UpdateAdminCategory(c *gin.Context, id string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	cat, pieCount, ok := loadCategory(wrpr, c, id)
	if !ok {
		return
	}

	label := strings.TrimSpace(c.Request.FormValue("label"))
	if label == "" {
		categoryEditPage(wrpr, c, http.StatusBadRequest, cat, pieCount, "Category label is required")
		return
	}

	// The slug is intentionally kept: it is the key pies reference.
	updated, err := wrpr.storage.UpdateCategory(c, types.Category{
		Id:    cat.Id,
		Slug:  cat.Slug,
		Label: label,
	})
	if errors.Is(err, storage.ErrCategoryNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	userId := c.Keys["userid"].(string)
	log.Info().Str("category", updated.Slug).Str("by", userId).Msg("Updated category")
	c.Redirect(http.StatusSeeOther, "/admin/categories/"+updated.Id.String())
}

// DeleteAdminCategory implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminCategory(c *gin.Context, id string) {
	if !validateAdmin(wrpr, c) {
		return
	}
	cat, pieCount, ok := loadCategory(wrpr, c, id)
	if !ok {
		return
	}

	// Guard: removing a in-use category silently un-categorises every pie
	// that used it, so admins must clear the usage first.
	if pieCount > 0 {
		categoryEditPage(wrpr, c, http.StatusBadRequest, cat, pieCount,
			"Cannot delete: "+countNoun(pieCount, "pie")+" still use this category. Remove the category from them first.")
		return
	}

	if err := wrpr.storage.DeleteCategory(c, cat.Id); err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(500, err)
		return
	}

	userId := c.Keys["userid"].(string)
	log.Info().Str("category", cat.Slug).Str("by", userId).Msg("Deleted category")
	c.Redirect(http.StatusSeeOther, "/admin/categories")
}

// loadCategory fetches a category plus its pie usage count, aborting with 400
// on a malformed id and 404 when absent.
func loadCategory(wrpr *WebsiteWrapper, c *gin.Context, id string) (*types.Category, int64, bool) {
	uid, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return nil, 0, false
	}
	cat, err := wrpr.storage.GetCategoryById(c, uid)
	if err != nil {
		c.AbortWithError(500, err)
		return nil, 0, false
	}
	if cat == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return nil, 0, false
	}
	pieCount, err := wrpr.storage.GetCategoryPieCount(c, uid)
	if err != nil {
		c.AbortWithError(500, err)
		return nil, 0, false
	}
	return cat, pieCount, true
}
