package website

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	generated "github.com/leemartin77/mincepierank.co.uk/internal/website/generated"
)

// MapPage implements generated.StrictServerInterface.
func (wrpr *WebsiteWrapper) MapPage(c context.Context, request generated.MapPageRequestObject) (generated.MapPageResponseObject, error) {
	ay, err := wrpr.storage.GetActiveYear(c)
	if err != nil {
		return returnMapPageUnexpectedError(err)
	}
	vals := templater.PageData{
		Head: templater.PageDataHead{
			Title:       "Map Page",
			Description: "Mince Pie Sightings Across the UK",
			Keywords:    "Mince Pies, UK, Ranking, Sightings",
			MenuSettings: templater.MenuSettings{
				ActiveYear: *ay,
				SignedIn:   c.(*gin.Context).Keys["signedin"].(bool),
			},
		},
		PageData: map[string]interface{}{},
	}

	rdr, err := wrpr.htmlTemplater.GeneratePage("map", vals)
	if err != nil {
		return nil, err
	}
	return generated.MapPage200TexthtmlResponse{
		Body: rdr,
	}, nil
}

func returnMapPageUnexpectedError(err error) (generated.MapPageResponseObject, error) {
	return generated.MapPagedefaultJSONResponse{
		Body: generated.Error{
			Code:    500,
			Message: err.Error(),
		},
		StatusCode: 500,
	}, nil
}
