package website

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

// YearPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearPage(c context.Context, request generated.YearPageRequestObject) (generated.YearPageResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return nil, err
	}
	topPie, err := wrpr.storage.GetTopMakerPie(c, request.Year)
	if err != nil {
		return nil, err
	}

	makers, err := wrpr.storage.GetMakersForYear(c, request.Year)
	if err != nil {
		return nil, err
	}

	cats, err := wrpr.storage.GetMakerPieCategoriesForYear(c, request.Year, []string{}, nil)
	if err != nil {
		return nil, err
	}
	categoryLinks := []templater.Link{}
	for _, ct := range *cats {
		categoryLinks = append(categoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", request.Year, ct.Slug), Label: ct.Label})
	}
	pieCategoryLinks := []templater.Link{}
	for _, ct := range *&topPie.Categories {
		pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", request.Year, ct.Slug), Label: ct.Label})
	}
	mrks := []templater.MakerCardData{}
	for _, mkr := range *makers {
		mrks = append(mrks, templater.MakerCardData{
			ImgprssrPrefix: imgprssrPrefix,
			Year:           request.Year,
			Maker:          mkr,
		})
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%d Pies", request.Year),
			Description: fmt.Sprintf("Mince pie categories and brands we have data for %d", request.Year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d", request.Year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.(*gin.Context).Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Year": request.Year,
			"TopPie": templater.PieCardData{
				Pie:            *topPie,
				ImgprssrPrefix: imgprssrPrefix,
				CategoryLinks:  pieCategoryLinks,
				HasDate:        false,
				PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", topPie.Year, topPie.MakerId, topPie.Id),
			},
			"MakerCards": mrks,
			"Categories": categoryLinks,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d", request.Year)),
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("year", vals)
	if err != nil {
		return nil, err
	}
	return generated.YearPage200TexthtmlResponse{
		Body: rdr,
	}, nil
}

// YearAllBrands implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllBrands(c context.Context, request generated.YearAllBrandsRequestObject) (generated.YearAllBrandsResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return nil, err
	}

	makers, err := wrpr.storage.GetMakersForYear(c, request.Year)
	if err != nil {
		return nil, err
	}

	mrks := []templater.MakerCardData{}
	for _, mkr := range *makers {
		mrks = append(mrks, templater.MakerCardData{
			ImgprssrPrefix: imgprssrPrefix,
			Year:           request.Year,
			Maker:          mkr,
		})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%d Brands", request.Year),
			Description: fmt.Sprintf("All brands we loaded data for in %d", request.Year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d", request.Year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.(*gin.Context).Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Year":       request.Year,
			"MakerCards": mrks,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/brands", request.Year)),
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("yearbrands", vals)
	if err != nil {
		return nil, err
	}
	return generated.YearAllBrands200TexthtmlResponse{
		Body: rdr,
	}, nil
}

// YearAllCategories implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearAllCategories(c context.Context, request generated.YearAllCategoriesRequestObject) (generated.YearAllCategoriesResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return nil, err
	}

	cats, err := wrpr.storage.GetMakerPieCategoriesForYear(c, request.Year, []string{}, nil)
	if err != nil {
		return nil, err
	}
	categoryLinks := []templater.Link{}
	for _, ct := range *cats {
		categoryLinks = append(categoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", request.Year, ct.Slug), Label: ct.Label})
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%d Categories", request.Year),
			Description: fmt.Sprintf("All categories we loaded data for in %d", request.Year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d", request.Year),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.(*gin.Context).Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Year":       request.Year,
			"Categories": categoryLinks,
			"Breadcrumb": templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/categories", request.Year)),
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("yearcategories", vals)
	if err != nil {
		return nil, err
	}
	return generated.YearAllCategories200TexthtmlResponse{
		Body: rdr,
	}, nil
}

// YearsPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) YearsPage(c context.Context, request generated.YearsPageRequestObject) (generated.YearsPageResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return nil, err
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Available Years",
			Description: "All the years Mince Pie Rank has data for",
			Keywords:    "Mince Pies, UK, Ranking",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.(*gin.Context).Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"Years":      []int64{2023, 2022},
			"Breadcrumb": templater.BreadcrumbsFromUrl("/years"),
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("years", vals)
	if err != nil {
		return nil, err
	}
	return generated.YearsPage200TexthtmlResponse{
		Body: rdr,
	}, nil
}
