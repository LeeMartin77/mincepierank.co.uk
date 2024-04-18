package website

import (
	"context"

	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/server.cfg.yaml ../../api/website.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/types.cfg.yaml ../../api/website.yaml

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
		},
		PageData: map[string]interface{}{
			"ActiveYear": *ay,
			"TopPie":     *topPie,
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
