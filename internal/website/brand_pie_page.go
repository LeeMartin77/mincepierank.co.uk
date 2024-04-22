package website

import (
	"context"
	"fmt"

	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

func returnPiePageUnexpectedError(err error) (generated.YearBrandPieResponseObject, error) {
	return generated.YearBrandPiedefaultJSONResponse{
		Body: generated.Error{
			Code:    500,
			Message: err.Error(),
		},
		StatusCode: 500,
	}, nil
}

func (wrpr *WebsiteWrapper) YearBrandPie(ctx context.Context, request generated.YearBrandPieRequestObject) (generated.YearBrandPieResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(ctx)
	if err != nil {
		return returnPiePageUnexpectedError(err)
	}

	pie, err := wrpr.storage.GetMakerPieYearly(ctx, int32(request.Year), request.Brand, request.Pie)
	if err != nil {
		return returnPiePageUnexpectedError(err)
	}
	if pie == nil {
		return generated.YearBrandPie404TexthtmlResponse{}, nil
	}

	maker, err := wrpr.storage.GetMaker(ctx, request.Brand)
	if err != nil {
		return returnPiePageUnexpectedError(err)
	}

	rankingSummary, err := wrpr.storage.GetMakerPieYearlyRankingSummary(ctx, int32(request.Year), request.Brand, request.Pie)
	if err != nil {
		return returnPiePageUnexpectedError(err)
	}

	cats, err := wrpr.storage.GetMakerPieCategoriesForMakerPieOid(ctx, pie.OId)
	if err != nil {
		return returnPiePageUnexpectedError(err)
	}

	pieCategoryLinks := []templater.Link{}
	for _, ct := range *cats {
		pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", pie.Year, ct.Slug), Label: ct.Label})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%s from %s :: %.2f/5", pie.DisplayName, maker.Name, rankingSummary.Average),
			Description: fmt.Sprintf("%s from %s :: scored  %.2f/5 on average in %d", pie.DisplayName, maker.Name, rankingSummary.Average, pie.Year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d, %s, %s", pie.Year, maker.Name, pie.DisplayName),
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
			},
		},
		PageData: map[string]interface{}{
			"CategoryLinks":  pieCategoryLinks,
			"ImgprssrPrefix": imgprssrPrefix,
			"HasMaker":       true,
			"Pie":            pie,
			"RankingSummary": rankingSummary,
			"Maker":          maker,
			"Breadcrumb":     templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/brands/%s/%s", pie.Year, pie.MakerId, pie.Id)),
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("singlepie", vals)
	if err != nil {
		return nil, err
	}
	return generated.YearBrandPie200TexthtmlResponse{
		Body: rdr,
	}, nil
}
