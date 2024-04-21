package website

import (
	"context"
	"fmt"
	"net/url"
	"slices"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

func getFilterLinks(wrpr *WebsiteWrapper, c context.Context, year int64, activeFilters []string, rootPath string) (*templater.FilterLinks, error) {
	cats, err := wrpr.storage.GetMakerPieCategoriesForYear(c, year)
	if err != nil {
		return nil, err
	}
	flinks := templater.FilterLinks{
		ActiveFilters:    []templater.Link{},
		AvailableFilters: []templater.Link{},
	}
	rootUrl, err := url.Parse(rootPath)
	if err != nil {
		return nil, err
	}
	activeQueries := rootUrl.Query()
	for _, ct := range activeFilters {
		activeQueries.Add("categories", ct)
	}
	for _, ct := range *cats {
		if slices.Contains(activeFilters, url.QueryEscape(ct)) {
			this := rootUrl.Query()
			for _, af := range activeFilters {
				if af == url.QueryEscape(ct) {
					continue
				}
				this.Add("categories", url.QueryEscape(af))
			}
			flinks.ActiveFilters = append(flinks.ActiveFilters, templater.Link{URL: rootUrl.String() + "?" + this.Encode(), Label: ct})
		} else {
			qry := rootUrl.Query()
			for k, aq := range activeQueries {
				for _, vl := range aq {
					qry.Add(k, vl)
				}
			}
			qry.Add("categories", url.QueryEscape(ct))
			flinks.AvailableFilters = append(flinks.AvailableFilters, templater.Link{URL: rootUrl.String() + "?" + qry.Encode(), Label: ct})
		}
	}
	return &flinks, err
}

// YearAllPies implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllPies(c context.Context, request generated.YearAllPiesRequestObject) (generated.YearAllPiesResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return nil, err
	}

	limit := int64(20)
	pageZeroIdx := int64(0)
	if request.Params.Limit != nil && *request.Params.Limit > 0 && *request.Params.Limit < 101 {
		limit = *request.Params.Limit
	}
	if request.Params.Page != nil && *request.Params.Page > 0 {
		pageZeroIdx = *request.Params.Page - 1
	}

	catFilters := []string{}
	unescapedCatFilters := []string{}
	if request.Params.Categories != nil {
		catFilters = *request.Params.Categories
		for _, cf := range catFilters {
			ue, err := url.QueryUnescape(cf)
			if err != nil {
				continue
			}
			unescapedCatFilters = append(unescapedCatFilters, ue)
		}
	}

	flinks, err := getFilterLinks(wrpr, c, request.Year, catFilters, fmt.Sprintf("/years/%d/all-pies", request.Year))
	if err != nil {
		return nil, err
	}

	pies, err := wrpr.storage.GetFilterableMakerPies(c, request.Year, limit, pageZeroIdx, storage.PieFilters{
		Categories: unescapedCatFilters,
	})
	if err != nil {
		return nil, err
	}

	pieCards := []templater.PieCardData{}

	for i, pie := range *pies {

		pieCategoryLinks := []templater.Link{}
		for _, ct := range pie.Labels {
			pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", pie.Year, url.QueryEscape(ct)), Label: ct})
		}
		pieCards = append(pieCards, templater.PieCardData{
			Pie:            pie,
			CategoryLinks:  pieCategoryLinks,
			ImgprssrPrefix: imgprssrPrefix,
			HasDate:        false,
			IsGold:         i == 0 && pie.Count > 0,
			PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", pie.Year, pie.MakerId, pie.Id),
		})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("All pies for %d", request.Year),
			Description: fmt.Sprintf("All pies for %d", request.Year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d", request.Year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
			},
		},
		PageData: map[string]interface{}{
			"Heading":     fmt.Sprintf("All pies for %d", request.Year),
			"Breadcrumb":  templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/all-pies", request.Year)),
			"PieCards":    pieCards,
			"FilterLinks": *flinks,
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("pielist", vals)
	if err != nil {
		return nil, err
	}
	return generated.YearAllPies200TexthtmlResponse{
		Body: rdr,
	}, nil
}

// YearBrandPies implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearBrandPies(c context.Context, request generated.YearBrandPiesRequestObject) (generated.YearBrandPiesResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return nil, err
	}

	maker, err := wrpr.storage.GetMaker(c, request.Brand)
	if err != nil {
		return nil, err
	}
	if maker == nil {
		return generated.YearBrandPies404TexthtmlResponse{}, nil
	}
	limit := int64(20)
	pageZeroIdx := int64(0)

	if request.Params.Limit != nil && *request.Params.Limit > 0 && *request.Params.Limit < 101 {
		limit = *request.Params.Limit
	}
	if request.Params.Page != nil && *request.Params.Page > 0 {
		pageZeroIdx = *request.Params.Page - 1
	}

	catFilters := []string{}
	unescapedCatFilters := []string{}
	if request.Params.Categories != nil {
		catFilters = *request.Params.Categories
		for _, cf := range catFilters {
			ue, err := url.QueryUnescape(cf)
			if err != nil {
				continue
			}
			unescapedCatFilters = append(unescapedCatFilters, ue)
		}
	}

	flinks, err := getFilterLinks(wrpr, c, request.Year, catFilters, fmt.Sprintf("/years/%d/brands/%s", request.Year, request.Brand))
	if err != nil {
		return nil, err
	}

	pies, err := wrpr.storage.GetFilterableMakerPies(c, request.Year, limit, pageZeroIdx, storage.PieFilters{
		BrandIds:   []string{request.Brand},
		Categories: unescapedCatFilters,
	})
	if err != nil {
		return nil, err
	}

	pieCards := []templater.PieCardData{}

	for i, pie := range *pies {

		pieCategoryLinks := []templater.Link{}
		for _, ct := range pie.Labels {
			pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", pie.Year, url.QueryEscape(ct)), Label: ct})
		}
		pieCards = append(pieCards, templater.PieCardData{
			Pie:            pie,
			CategoryLinks:  pieCategoryLinks,
			ImgprssrPrefix: imgprssrPrefix,
			HasDate:        false,
			IsGold:         i == 0 && pie.Count > 0,
			PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", pie.Year, pie.MakerId, pie.Id),
		})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%s pies for %d", maker.Name, request.Year),
			Description: fmt.Sprintf("%s for %d", maker.Name, request.Year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d", request.Year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
			},
		},
		PageData: map[string]interface{}{
			"Heading":     fmt.Sprintf("%s pies for %d", maker.Name, request.Year),
			"Breadcrumb":  templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/brands/%s", request.Year, request.Brand)),
			"PieCards":    pieCards,
			"FilterLinks": *flinks,
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("pielist", vals)
	if err != nil {
		return nil, err
	}
	return generated.YearBrandPies200TexthtmlResponse{
		Body: rdr,
	}, nil
}

// YearCategoryPies implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearCategoryPies(c context.Context, request generated.YearCategoryPiesRequestObject) (generated.YearCategoryPiesResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return nil, err
	}
	cat, err := url.QueryUnescape(request.Category)
	if err != nil {
		return nil, err
	}

	limit := int64(20)
	pageZeroIdx := int64(0)

	if request.Params.Limit != nil && *request.Params.Limit > 0 && *request.Params.Limit < 101 {
		limit = *request.Params.Limit
	}
	if request.Params.Page != nil && *request.Params.Page > 0 {
		pageZeroIdx = *request.Params.Page - 1
	}

	catFilters := []string{}
	unescapedCatFilters := []string{}
	if request.Params.Categories != nil {
		catFilters = *request.Params.Categories
		for _, cf := range catFilters {
			ue, err := url.QueryUnescape(cf)
			if err != nil {
				continue
			}
			unescapedCatFilters = append(unescapedCatFilters, ue)
		}
	}

	flinks, err := getFilterLinks(wrpr, c, request.Year, catFilters, fmt.Sprintf("/years/%d/categories/%s", request.Year, url.QueryEscape(request.Category)))
	if err != nil {
		return nil, err
	}

	unescapedCatFilters = append(unescapedCatFilters, cat)
	pies, err := wrpr.storage.GetFilterableMakerPies(c, request.Year, limit, pageZeroIdx, storage.PieFilters{
		Categories: unescapedCatFilters,
	})
	if err != nil {
		return nil, err
	}

	pieCards := []templater.PieCardData{}

	for i, pie := range *pies {

		pieCategoryLinks := []templater.Link{}
		for _, ct := range pie.Labels {
			pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", pie.Year, url.QueryEscape(ct)), Label: ct})
		}
		pieCards = append(pieCards, templater.PieCardData{
			Pie:            pie,
			CategoryLinks:  pieCategoryLinks,
			ImgprssrPrefix: imgprssrPrefix,
			HasDate:        false,
			IsGold:         i == 0 && pie.Count > 0,
			PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", pie.Year, pie.MakerId, pie.Id),
		})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%s pies for %d", cat, request.Year),
			Description: fmt.Sprintf("%s pies for %d", cat, request.Year),
			Keywords:    fmt.Sprintf("Mince pies, UK, ranking, %d", request.Year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
			},
		},
		PageData: map[string]interface{}{
			"Heading":     fmt.Sprintf("%s pies for %d", cat, request.Year),
			"Breadcrumb":  templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/categories/%s", request.Year, request.Category)),
			"PieCards":    pieCards,
			"FilterLinks": *flinks,
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("pielist", vals)
	if err != nil {
		return nil, err
	}
	return generated.YearCategoryPies200TexthtmlResponse{
		Body: rdr,
	}, nil
}
