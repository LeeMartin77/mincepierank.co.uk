package website

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

func returnHomepageUnexpectedError(err error) (generated.HomePageResponseObject, error) {
	return generated.HomePagedefaultJSONResponse{
		Body: generated.Error{
			Code:    500,
			Message: err.Error(),
		},
		StatusCode: 500,
	}, nil
}

var imgprssrPrefix = "https://static.mincepierank.co.uk"

func (wrpr *WebsiteWrapper) HomePage(c context.Context, req generated.HomePageRequestObject) (generated.HomePageResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return returnHomepageUnexpectedError(err)
	}

	topPie, err := wrpr.storage.GetTopMakerPie(c, *ay)
	if err != nil {
		return returnHomepageUnexpectedError(err)
	}

	makers, err := wrpr.storage.GetMakersForYear(c, *ay)
	if err != nil {
		return returnHomepageUnexpectedError(err)
	}
	mrks := []templater.MakerCardData{}
	for _, mkr := range *makers {
		mrks = append(mrks, templater.MakerCardData{
			ImgprssrPrefix: imgprssrPrefix,
			Year:           *ay,
			Maker:          mkr,
		})
	}

	pieCategoryLinks := []templater.Link{}
	for _, ct := range topPie.Categories {
		pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", topPie.Year, ct.Slug), Label: ct.Label})
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Home Page",
			Description: "Rank mince pies from major UK brands",
			Keywords:    "Mince Pies, UK, Ranking",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.(*gin.Context).Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{
			"TopPie": templater.PieCardData{
				Pie:            *topPie,
				CategoryLinks:  pieCategoryLinks,
				ImgprssrPrefix: imgprssrPrefix,
				HasDate:        false,
				PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", topPie.Year, topPie.MakerId, topPie.Id),
			},
			"MakerCards": mrks,
		},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("index", vals)
	if err != nil {
		return nil, err
	}
	return generated.HomePage200TexthtmlResponse{
		Body: rdr,
	}, nil
}
