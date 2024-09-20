package website

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

// YearPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearPage(c *gin.Context, year int64) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}
	topPie, err := wrpr.storage.GetTopMakerPie(c, year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting top maker pie")
		c.AbortWithStatus(500)
		return
	}

	makers, err := wrpr.storage.GetMakersForYear(c, year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting makers")
		c.AbortWithStatus(500)
		return
	}

	cats, err := wrpr.storage.GetMakerPieCategoriesForYear(c, year, []string{}, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error getting pie categories")
		c.AbortWithStatus(500)
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
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%d Pies", year),
			Description: fmt.Sprintf("Mince pie categories and brands we have data for %d", year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
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

	c.HTML(200, "page:year", vals)
}

// YearAllBrands implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllBrands(c *gin.Context, year int64) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	makers, err := wrpr.storage.GetMakersForYear(c, year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting makers")
		c.AbortWithStatus(500)
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

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%d Brands", year),
			Description: fmt.Sprintf("All brands we loaded data for in %d", year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"MakerCards": mrks,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/brands", year)),
		},
	}

	c.HTML(200, "page:yearbrands", vals)
}

// YearAllCategories implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllCategories(c *gin.Context, year int64) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	cats, err := wrpr.storage.GetMakerPieCategoriesForYear(c, year, []string{}, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error getting GetMakerPieCategoriesForYear")
		c.AbortWithStatus(500)
		return
	}
	categoryLinks := []templater.Link{}
	for _, ct := range *cats {
		categoryLinks = append(categoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", year, ct.Slug), Label: ct.Label})
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%d Categories", year),
			Description: fmt.Sprintf("All categories we loaded data for in %d", year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Year":       year,
			"Categories": categoryLinks,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/categories", year)),
		},
	}

	c.HTML(200, "page:yearcategories", vals)
}

// YearsPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearsPage(c *gin.Context) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}
	yrs, err := wrpr.storage.GetQuerier().GetMakerPieYears(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting years")
		c.AbortWithStatus(500)
		return
	}
	slices.Sort(yrs)
	slices.Reverse(yrs)
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Available Years",
			Description: "All the years Mince Pie Rank has data for",
			Keywords:    "Mince Pies, UK, Ranking",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Years":      yrs,
			"Breadcrumb": templater.BreadcrumbsFromUrl("/years"),
		},
	}

	c.HTML(200, "page:years", vals)
}
