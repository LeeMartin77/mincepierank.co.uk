package website

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

// YearPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearPage(w http.ResponseWriter, r *http.Request, year int64) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}
	topPie, err := wrpr.storage.GetTopMakerPie(r.Context(), year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting top maker pie")
		abortStatus(w, 500)
		return
	}

	makers, err := wrpr.storage.GetMakersForYear(r.Context(), year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting makers")
		abortStatus(w, 500)
		return
	}

	cats, err := wrpr.storage.GetMakerPieCategoriesForYear(r.Context(), year, []string{}, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error getting pie categories")
		abortStatus(w, 500)
		return
	}
	categoryLinks := []templater.Link{}
	for _, ct := range *cats {
		categoryLinks = append(categoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", year, ct.Slug), Label: ct.Label})
	}
	mrks := []templater.MakerCardData{}
	for _, mkr := range *makers {
		mrks = append(mrks, templater.MakerCardData{
			ImgprssrPrefix: wrpr.config.ImgprssrPrefix,
			Year:           year,
			Maker:          mkr,
		})
	}
	pageData := map[string]interface{}{
		"Year":       year,
		"MakerCards": mrks,
		"Categories": categoryLinks,
		"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d", year)),
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       fmt.Sprintf("%d Pies", year),
			Description: fmt.Sprintf("Mince pie categories and brands we have data for %d", year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        SignedInFrom(r),
				UserPiesEnabled: piesEnabled,
			},
		},
	}
	if topPie != nil {

		pieCategoryLinks := []templater.Link{}
		for _, ct := range topPie.Categories {
			pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", year, ct.Slug), Label: ct.Label})
		}
		topPieData := templater.PieCardData{
			Pie:            *topPie,
			ImgprssrPrefix: wrpr.config.ImgprssrPrefix,
			CategoryLinks:  pieCategoryLinks,
			HasDate:        false,
			PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", topPie.Year, topPie.MakerId, topPie.Id),
		}
		pageData["TopPie"] = topPieData
	}
	vals.PageData = pageData

	wrpr.htmlTemplater.WritePage(w, 200, "year", vals)
}

// YearAllBrands implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllBrands(w http.ResponseWriter, r *http.Request, year int64) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	makers, err := wrpr.storage.GetMakersForYear(r.Context(), year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting makers")
		abortStatus(w, 500)
		return
	}

	mrks := []templater.MakerCardData{}
	for _, mkr := range *makers {
		mrks = append(mrks, templater.MakerCardData{
			ImgprssrPrefix: wrpr.config.ImgprssrPrefix,
			Year:           year,
			Maker:          mkr,
		})
	}

	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       fmt.Sprintf("%d Brands", year),
			Description: fmt.Sprintf("All brands we loaded data for in %d", year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        SignedInFrom(r),
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"MakerCards": mrks,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/brands", year)),
		},
	}

	wrpr.htmlTemplater.WritePage(w, 200, "yearbrands", vals)
}

// YearAllCategories implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllCategories(w http.ResponseWriter, r *http.Request, year int64) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	cats, err := wrpr.storage.GetMakerPieCategoriesForYear(r.Context(), year, []string{}, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error getting GetMakerPieCategoriesForYear")
		abortStatus(w, 500)
		return
	}
	categoryLinks := []templater.Link{}
	for _, ct := range *cats {
		categoryLinks = append(categoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", year, ct.Slug), Label: ct.Label})
	}
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       fmt.Sprintf("%d Categories", year),
			Description: fmt.Sprintf("All categories we loaded data for in %d", year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        SignedInFrom(r),
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"Categories": categoryLinks,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/categories", year)),
		},
	}

	wrpr.htmlTemplater.WritePage(w, 200, "yearcategories", vals)
}

// YearsPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearsPage(w http.ResponseWriter, r *http.Request) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}
	yrs, err := wrpr.storage.GetQuerier().GetMakerPieYears(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting years")
		abortStatus(w, 500)
		return
	}
	slices.Sort(yrs)
	slices.Reverse(yrs)
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Snow:        wrpr.snowing(r.Context()),
			Title:       "Available Years",
			Description: "All the years Mince Pie Rank has data for",
			Keywords:    "Mince Pies, UK, Ranking",
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        SignedInFrom(r),
				UserPiesEnabled: piesEnabled,
			},
		},
		PageData: map[string]interface{}{
			"Years":      yrs,
			"Breadcrumb": templater.BreadcrumbsFromUrl("/years"),
		},
	}

	wrpr.htmlTemplater.WritePage(w, 200, "years", vals)
}
