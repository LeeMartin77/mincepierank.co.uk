package website

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

func (wrpr *WebsiteWrapper) HomePage(c *gin.Context) {
	ay, err := wrpr.storage.GetQuerier().GetActiveYear(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		c.AbortWithStatus(500)
		return
	}

	topPie, err := wrpr.storage.GetTopMakerPie(c, ay)
	if err != nil {
		log.Error().Err(err).Msg("Error getting top maker pie")
		c.AbortWithStatus(500)
		return
	}

	makers, err := wrpr.storage.GetMakersForYear(c, ay)
	if err != nil {
		log.Error().Err(err).Msg("Error getting year makers")
		c.AbortWithStatus(500)
		return
	}
	mrks := []templater.MakerCardData{}
	for _, mkr := range *makers {
		mrks = append(mrks, templater.MakerCardData{
			ImgprssrPrefix: imgprssrPrefix,
			Year:           ay,
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
				ActiveYear: ay,
				SignedIn:   c.Keys["signedin"].(bool),
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

	c.HTML(http.StatusOK, "page:index", vals)
}
