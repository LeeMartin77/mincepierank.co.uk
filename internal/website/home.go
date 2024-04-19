package website

import (
	"context"
	"fmt"

	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

type PieCardData struct {
	ImgprssrPrefix string
	HasDate        bool
	PieLink        string
	Pie            interface{}
}

func (wrpr *WebsiteWrapper) HomePage(c context.Context, req generated.HomePageRequestObject) (generated.HomePageResponseObject, error) {
	// urlExample := "postgres://username:password@localhost:5432/mincepierank"
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return generated.HomePagedefaultJSONResponse{
			Body: generated.Error{
				Code:    500,
				Message: err.Error(),
			},
			StatusCode: 500,
		}, nil
	}

	topPie, err := wrpr.storage.GetTopMakerPie(c, *ay)
	if err != nil {
		return generated.HomePagedefaultJSONResponse{
			Body: generated.Error{
				Code:    500,
				Message: err.Error(),
			},
			StatusCode: 500,
		}, nil
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
			"TopPie": PieCardData{
				Pie:            *topPie,
				ImgprssrPrefix: "https://static.mincepierank.co.uk",
				HasDate:        false,
				PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", topPie.Year, topPie.MakerId, topPie.Id),
			},
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
