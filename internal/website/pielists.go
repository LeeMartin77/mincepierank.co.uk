package website

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
	"github.com/rs/zerolog/log"
)

func getFilterLinks(wrpr *WebsiteWrapper, c context.Context, year int64, activeFilters []string, rootPath string, omitSlug *string, brandSlug *string) (*templater.FilterLinks, *url.Values, error) {
	reqslgs := activeFilters
	if omitSlug != nil {
		reqslgs = append(reqslgs, *omitSlug)
	}
	cats, err := wrpr.storage.GetMakerPieCategoriesForYear(c, year, reqslgs, brandSlug)
	if err != nil {
		return nil, nil, err
	}
	flinks := templater.FilterLinks{
		ActiveFilters:    []templater.Link{},
		AvailableFilters: []templater.Link{},
	}
	rootUrl, err := url.Parse(rootPath)
	if err != nil {
		return nil, nil, err
	}
	activeQueries := rootUrl.Query()
	for _, ct := range activeFilters {
		activeQueries.Add("categories", ct)
	}
	for _, ct := range *cats {
		if omitSlug != nil && ct.Slug == *omitSlug {
			continue
		}
		if slices.Contains(activeFilters, ct.Slug) {
			this := rootUrl.Query()
			for _, af := range activeFilters {
				if af == ct.Slug {
					continue
				}
				this.Add("categories", af)
			}
			flinks.ActiveFilters = append(flinks.ActiveFilters, templater.Link{URL: rootUrl.String() + "?" + this.Encode(), Label: ct.Label})
		} else {
			qry := rootUrl.Query()
			for k, aq := range activeQueries {
				for _, vl := range aq {
					qry.Add(k, vl)
				}
			}
			qry.Add("categories", ct.Slug)
			flinks.AvailableFilters = append(flinks.AvailableFilters, templater.Link{URL: rootUrl.String() + "?" + qry.Encode(), Label: ct.Label})
		}
	}
	return &flinks, &activeQueries, err
}

func getLimitAndIdx(c *gin.Context) (int64, int64) {
	limit := int64(20)
	pageZeroIdx := int64(0)

	prsdLimit, err := strconv.ParseInt(c.Query("limit"), 10, 64)

	if err == nil && prsdLimit > 0 && prsdLimit < 101 {
		limit = prsdLimit
	}
	prsdIdx, err := strconv.ParseInt(c.Query("page"), 10, 64)

	if err == nil && prsdIdx > 0 {
		pageZeroIdx = prsdIdx - 1
	}
	return limit, pageZeroIdx
}

// YearAllPies implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllPies(c *gin.Context, year int64, params generated.YearAllPiesParams) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	limit, pageZeroIdx := getLimitAndIdx(c)

	catFilters := c.QueryArray("categories")

	flinks, activequery, err := getFilterLinks(wrpr, c, year, catFilters, fmt.Sprintf("/years/%d/all-pies", year), nil, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error getting filter links")
		c.AbortWithStatus(500)
		return
	}

	pies, err := wrpr.storage.GetFilterableMakerPies(c, year, limit, pageZeroIdx, storage.PieFilters{
		CategorySlugs: catFilters,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker pies")
		c.AbortWithStatus(500)
		return
	}

	mmap, err := wrpr.getMakerMap(c, year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker pies")
		c.AbortWithStatus(500)
		return
	}

	pieCards := []templater.PieCardData{}

	for i, pie := range *pies {

		pieCategoryLinks := []templater.Link{}
		for _, ct := range pie.Categories {
			pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", pie.Year, ct.Slug), Label: ct.Label})
		}
		pieCards = append(pieCards, templater.PieCardData{
			Pie:            pie,
			CategoryLinks:  pieCategoryLinks,
			ImgprssrPrefix: wrpr.config.ImgprssrPrefix,
			HasDate:        false,
			IsGold:         i == 0 && pie.Count > 0,
			PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", pie.Year, pie.MakerId, pie.Id),
			Maker:          mmap[pie.MakerId],
		})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("All pies for %d", year),
			Description: fmt.Sprintf("All pies for %d", year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Heading":            fmt.Sprintf("All pies for %d", year),
			"Breadcrumb":         templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/all-pies", year)),
			"PieCards":           pieCards,
			"FilterLinks":        *flinks,
			"NextPage":           pageZeroIdx + 2,
			"ShouldLoadNextPage": len(pieCards) == int(limit),
			"URL":                fmt.Sprintf("/years/%d/all-pies", year),
			"Query":              activequery.Encode(),
		},
	}

	c.HTML(200, "page:pielist", vals)
}

func (wrpr *WebsiteWrapper) YearPersonalRanking(c *gin.Context, year int64, params generated.YearPersonalRankingParams) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	limit, pageZeroIdx := getLimitAndIdx(c)

	catFilters := c.QueryArray("categories")

	flinks, activequery, err := getFilterLinks(wrpr, c, year, catFilters, fmt.Sprintf("/years/%d/all-pies", year), nil, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error getting filter links")
		c.AbortWithStatus(500)
		return
	}

	pies, err := wrpr.storage.GetFilterableMakerPiesForUser(c, year, c.Keys["userid"].(string), limit, pageZeroIdx, storage.PieFilters{
		CategorySlugs: catFilters,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error getting user rankings")
		c.AbortWithStatus(500)
		return
	}

	mmap, err := wrpr.getMakerMap(c, year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker pies")
		c.AbortWithStatus(500)
		return
	}

	pieCards := []templater.PieCardData{}

	for i, pie := range *pies {

		pieCategoryLinks := []templater.Link{}
		for _, ct := range pie.Categories {
			pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", pie.Year, ct.Slug), Label: ct.Label})
		}
		pieCards = append(pieCards, templater.PieCardData{
			Pie:            pie,
			CategoryLinks:  pieCategoryLinks,
			ImgprssrPrefix: wrpr.config.ImgprssrPrefix,
			HasDate:        false,
			IsGold:         i == 0 && pie.Count > 0,
			PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", pie.Year, pie.MakerId, pie.Id),
			Maker:          mmap[pie.MakerId],
		})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("Your rankings for %d", year),
			Description: fmt.Sprintf("Your rankings for %d", year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Heading":            fmt.Sprintf("Your rankings for %d", year),
			"Breadcrumb":         []templater.Link{},
			"PieCards":           pieCards,
			"FilterLinks":        *flinks,
			"NextPage":           pageZeroIdx + 2,
			"ShouldLoadNextPage": len(pieCards) == int(limit),
			"URL":                fmt.Sprintf("/profile/rankings/%d", year),
			"Query":              activequery.Encode(),
		},
	}

	c.HTML(200, "page:pielist", vals)
}

// YearBrandPies implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearBrandPies(c *gin.Context, year int64, brand string, params generated.YearBrandPiesParams) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	maker, err := wrpr.storage.GetMaker(c, c.Param("brand"))
	if err != nil {
		log.Error().Err(err).Msg("Error getting brand")
		c.AbortWithStatus(500)
		return
	}
	if maker == nil {
		c.AbortWithStatus(404)
		return
	}

	limit, pageZeroIdx := getLimitAndIdx(c)

	catFilters := c.QueryArray("categories")
	unescapedCatFilters := []string{}
	for _, cf := range catFilters {
		ue, err := url.QueryUnescape(cf)
		if err != nil {
			continue
		}
		unescapedCatFilters = append(unescapedCatFilters, ue)
	}

	flinks, activequery, err := getFilterLinks(wrpr, c, year, catFilters, fmt.Sprintf("/years/%d/brands/%s", year, maker.Id), nil, &maker.Id)
	if err != nil {
		log.Error().Err(err).Msg("Error getting filter links")
		c.AbortWithStatus(500)
		return
	}

	pies, err := wrpr.storage.GetFilterableMakerPies(c, year, limit, pageZeroIdx, storage.PieFilters{
		BrandIds:      []string{maker.Id},
		CategorySlugs: unescapedCatFilters,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error getting pies")
		c.AbortWithStatus(500)
		return
	}

	pieCards := []templater.PieCardData{}

	for i, pie := range *pies {

		pieCategoryLinks := []templater.Link{}
		for _, ct := range pie.Categories {
			pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", pie.Year, ct.Slug), Label: ct.Label})
		}
		pieCards = append(pieCards, templater.PieCardData{
			Pie:            pie,
			CategoryLinks:  pieCategoryLinks,
			ImgprssrPrefix: wrpr.config.ImgprssrPrefix,
			HasDate:        false,
			IsGold:         i == 0 && pie.Count > 0,
			PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", pie.Year, pie.MakerId, pie.Id),
			Maker:          *maker,
		})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%s pies for %d", maker.Name, year),
			Description: fmt.Sprintf("%s for %d", maker.Name, year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Heading":            fmt.Sprintf("%s pies for %d", maker.Name, year),
			"Breadcrumb":         templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/brands/%s", year, maker.Id)),
			"PieCards":           pieCards,
			"FilterLinks":        *flinks,
			"NextPage":           pageZeroIdx + 2,
			"ShouldLoadNextPage": len(pieCards) == int(limit),
			"URL":                fmt.Sprintf("/years/%d/brands/%s", year, maker.Id),
			"Query":              activequery.Encode(),
		},
	}

	c.HTML(200, "page:pielist", vals)
}

// YearCategoryPies implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearCategoryPies(c *gin.Context, year int64, brand string, params generated.YearCategoryPiesParams) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}
	cat, err := url.QueryUnescape(c.Param("category"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	limit, pageZeroIdx := getLimitAndIdx(c)

	catFilters := c.QueryArray("categories")

	flinks, activequery, err := getFilterLinks(wrpr, c, year, catFilters, fmt.Sprintf("/years/%d/categories/%s", year, url.QueryEscape(c.Param("category"))), &cat, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error getting filter links")
		c.AbortWithStatus(500)
		return
	}

	catFilters = append(catFilters, cat)
	pies, err := wrpr.storage.GetFilterableMakerPies(c, year, limit, pageZeroIdx, storage.PieFilters{
		CategorySlugs: catFilters,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker pies")
		c.AbortWithStatus(500)
		return
	}

	mmap, err := wrpr.getMakerMap(c, year)
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker pies")
		c.AbortWithStatus(500)
		return
	}

	pieCards := []templater.PieCardData{}

	for i, pie := range *pies {

		pieCategoryLinks := []templater.Link{}
		for _, ct := range pie.Categories {
			pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", pie.Year, ct.Slug), Label: ct.Label})
		}
		pieCards = append(pieCards, templater.PieCardData{
			Pie:            pie,
			CategoryLinks:  pieCategoryLinks,
			ImgprssrPrefix: wrpr.config.ImgprssrPrefix,
			HasDate:        false,
			IsGold:         i == 0 && pie.Count > 0,
			PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", pie.Year, pie.MakerId, pie.Id),
			Maker:          mmap[pie.MakerId],
		})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%s pies for %d", cat, year),
			Description: fmt.Sprintf("%s pies for %d", cat, year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d", year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Heading":            fmt.Sprintf("%s pies for %d", cat, year),
			"Breadcrumb":         templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/categories/%s", year, c.Param("category"))),
			"PieCards":           pieCards,
			"FilterLinks":        *flinks,
			"NextPage":           pageZeroIdx + 2,
			"ShouldLoadNextPage": len(pieCards) == int(limit),
			"URL":                fmt.Sprintf("/years/%d/categories/%s", year, c.Param("category")),
			"Query":              activequery.Encode(),
		},
	}

	c.HTML(200, "page:pielist", vals)
}
