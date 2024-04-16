package website

import (
	"context"
	"fmt"

	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/server.cfg.yaml ../../api/website.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=../../tools/types.cfg.yaml ../../api/website.yaml

var (
	ErrShutdown error = fmt.Errorf("error shutting down server")
)

func (wrpr *WebsiteWrapper) HomePage(c context.Context, req generated.HomePageRequestObject) (generated.HomePageResponseObject, error) {
	// urlExample := "postgres://username:password@localhost:5432/mincepierank"
	str, err := wrpr.storage.TestData(c)
	if err != nil {
		return generated.HomePagedefaultJSONResponse{
			Body: generated.Error{
				Code:    500,
				Message: err.Error(),
			},
			StatusCode: 500,
		}, nil
	}

	vals := map[string]interface{}{
		"Head": map[string]interface{}{
			"Title": "Index Page",
		},
		"PageData": map[string]interface{}{
			"TestData": *str,
		},
		"Foot": map[string]interface{}{},
	}

	rdr, err := wrpr.htmlTemplater.GenerateTemplate("page:index", vals)
	if err != nil {
		return nil, err
	}
	return generated.HomePage200TexthtmlResponse{
		Body: rdr,
	}, nil
}
