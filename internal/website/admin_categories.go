package website

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

// validateCategoryForm checks the label and slug. The slug is entirely user
// controlled — it is used in public category URLs, so the user must supply a
// suitable one instead of having one derived from the label.
func validateCategoryForm(label, slug string) (string, string, []string) {
	errs := []string{}
	label = strings.TrimSpace(label)
	if label == "" {
		errs = append(errs, "Category label is required")
	}

	slug = strings.TrimSpace(slug)
	if slug == "" {
		errs = append(errs, "Category slug is required — pick a short lowercase slug such as \"vegetarian\"")
	} else if !brandIDPattern.MatchString(slug) {
		errs = append(errs, "Category slug must be lowercase letters, numbers and dashes (e.g. \"vegetarian\")")
	}
	return label, slug, errs
}

// categoryListPage renders the admin categories list with per-category pie
// usage counts.
func categoryListPage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, pageErrs ...string) {
	cats, err := wrpr.storage.GetAllCategoriesAdmin(r.Context())
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

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "admin:categories", vals)
}

// categoryNewPage renders the category create form, keeping entered values on
// error.
func categoryNewPage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, status int, form map[string]string, pageErrs ...string) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		abortErr(w, 500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
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

	wrpr.htmlTemplater.WritePage(w, status, "admin:categorynew", vals)
}

// categoryEditPage renders the category edit form. The slug is shown but not
// editable: pies reference categories by slug, so changing it would break
// bulk imports and saved selections.
func categoryEditPage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, status int, cat *types.Category, pieCount int64, pageErrs ...string) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(r.Context())
	if err != nil {
		abortErr(w, 500, err)
		return
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
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

	wrpr.htmlTemplater.WritePage(w, status, "admin:category", vals)
}

// GetAdminCategories implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminCategories(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	categoryListPage(wrpr, w, r)
}

// GetAdminCategoryNew implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminCategoryNew(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	categoryNewPage(wrpr, w, r, http.StatusOK, map[string]string{})
}

// CreateAdminCategory implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) CreateAdminCategory(w http.ResponseWriter, r *http.Request) {
	if !validateAdmin(wrpr, w, r) {
		return
	}

	label, slug, errs := validateCategoryForm(r.FormValue("label"), r.FormValue("slug"))
	if len(errs) > 0 {
		categoryNewPage(wrpr, w, r, http.StatusBadRequest, map[string]string{
			"label": label, "slug": r.FormValue("slug"),
		}, errs...)
		return
	}

	created, err := wrpr.storage.CreateCategory(r.Context(), types.Category{Slug: slug, Label: label})
	if errors.Is(err, storage.ErrCategoryAlreadyExists) {
		categoryNewPage(wrpr, w, r, http.StatusBadRequest, map[string]string{
			"label": label, "slug": slug,
		}, "A category with the slug \""+slug+"\" already exists")
		return
	}
	if err != nil {
		abortErr(w, 500, err)
		return
	}

	userId := userIDFrom(r)
	log.Info().Str("category", created.Slug).Str("by", userId).Msg("Created category")
	http.Redirect(w, r, "/admin/categories/"+created.Id.String(), http.StatusSeeOther)
}

// GetAdminCategory implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) GetAdminCategory(w http.ResponseWriter, r *http.Request, id string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	cat, pieCount, ok := loadCategory(wrpr, w, r, id)
	if !ok {
		return
	}
	categoryEditPage(wrpr, w, r, http.StatusOK, cat, pieCount)
}

// UpdateAdminCategory implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) UpdateAdminCategory(w http.ResponseWriter, r *http.Request, id string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	cat, pieCount, ok := loadCategory(wrpr, w, r, id)
	if !ok {
		return
	}

	label := strings.TrimSpace(r.FormValue("label"))
	if label == "" {
		categoryEditPage(wrpr, w, r, http.StatusBadRequest, cat, pieCount, "Category label is required")
		return
	}

	// The slug is intentionally kept: it is the key pies reference.
	updated, err := wrpr.storage.UpdateCategory(r.Context(), types.Category{
		Id:    cat.Id,
		Slug:  cat.Slug,
		Label: label,
	})
	if errors.Is(err, storage.ErrCategoryNotFound) {
		abortStatus(w, http.StatusNotFound)
		return
	}
	if err != nil {
		abortErr(w, 500, err)
		return
	}

	userId := userIDFrom(r)
	log.Info().Str("category", updated.Slug).Str("by", userId).Msg("Updated category")
	http.Redirect(w, r, "/admin/categories/"+updated.Id.String(), http.StatusSeeOther)
}

// DeleteAdminCategory implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) DeleteAdminCategory(w http.ResponseWriter, r *http.Request, id string) {
	if !validateAdmin(wrpr, w, r) {
		return
	}
	cat, pieCount, ok := loadCategory(wrpr, w, r, id)
	if !ok {
		return
	}

	// Guard: removing a in-use category silently un-categorises every pie
	// that used it, so admins must clear the usage first.
	if pieCount > 0 {
		categoryEditPage(wrpr, w, r, http.StatusBadRequest, cat, pieCount,
			"Cannot delete: "+countNoun(pieCount, "pie")+" still use this category. Remove the category from them first.")
		return
	}

	if err := wrpr.storage.DeleteCategory(r.Context(), cat.Id); err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			abortStatus(w, http.StatusNotFound)
			return
		}
		abortErr(w, 500, err)
		return
	}

	userId := userIDFrom(r)
	log.Info().Str("category", cat.Slug).Str("by", userId).Msg("Deleted category")
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

// loadCategory fetches a category plus its pie usage count, aborting with 400
// on a malformed id and 404 when absent.
func loadCategory(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, id string) (*types.Category, int64, bool) {
	uid, err := uuid.Parse(id)
	if err != nil {
		abortStatus(w, http.StatusBadRequest)
		return nil, 0, false
	}
	cat, err := wrpr.storage.GetCategoryById(r.Context(), uid)
	if err != nil {
		abortErr(w, 500, err)
		return nil, 0, false
	}
	if cat == nil {
		abortStatus(w, http.StatusNotFound)
		return nil, 0, false
	}
	pieCount, err := wrpr.storage.GetCategoryPieCount(r.Context(), uid)
	if err != nil {
		abortErr(w, 500, err)
		return nil, 0, false
	}
	return cat, pieCount, true
}
