package website

import (
	"context"
	"fmt"
	"net/url"

	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

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

	pies, err := wrpr.storage.GetFilterableMakerPies(c, request.Year, limit, pageZeroIdx)
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
			"Heading":    fmt.Sprintf("All pies for %d", request.Year),
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/all-pies", request.Year)),
			"PieCards":   pieCards,
			"FilterLinks": map[string]interface{}{
				"ActiveFilters":    []templater.Link{},
				"AvailableFilters": []templater.Link{},
			},
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
func (wrpr *WebsiteWrapper) YearBrandPies(ctx context.Context, request generated.YearBrandPiesRequestObject) (generated.YearBrandPiesResponseObject, error) {
	panic("unimplemented")
}

// YearCategoryPies implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearCategoryPies(ctx context.Context, request generated.YearCategoryPiesRequestObject) (generated.YearCategoryPiesResponseObject, error) {
	panic("unimplemented")
}
