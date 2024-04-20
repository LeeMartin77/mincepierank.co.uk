package website

import (
	"context"
	"fmt"
	"net/url"

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

	pieCategoryLinks := []templater.Link{}
	for _, ct := range pie.Labels {
		pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", pie.Year, url.QueryEscape(ct)), Label: ct})
	}

	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       fmt.Sprintf("%s from %s :: %d/5", pie.DisplayName, maker.Name, 0),
			Description: fmt.Sprintf("%s from %s :: scored  %d/5 on average in %d", pie.DisplayName, maker.Name, 0, pie.Year),
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
