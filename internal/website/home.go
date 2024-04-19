package website

import (
	"context"
	"fmt"

	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

func returnUnexpectedError(err error) (generated.HomePageResponseObject, error) {
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
	// urlExample := "postgres://username:password@localhost:5432/mincepierank"
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return returnUnexpectedError(err)
	}

	topPie, err := wrpr.storage.GetTopMakerPie(c, *ay)
	if err != nil {
		return returnUnexpectedError(err)
	}

	makers, err := wrpr.storage.GetMakersForYear(c, *ay)
	if err != nil {
		return returnUnexpectedError(err)
	}
	mrks := []templater.MakerCardData{}
	for _, mkr := range *makers {
		mrks = append(mrks, templater.MakerCardData{
			ImgprssrPrefix: imgprssrPrefix,
			Year:           *ay,
			Maker:          mkr,
		})
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title: "Home Page",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
			},
		},
		PageData: map[string]interface{}{
			"ActiveYear": *ay,
			"TopPie": templater.PieCardData{
				Pie:            *topPie,
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
