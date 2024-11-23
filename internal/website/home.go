package website

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
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

	ro, err := wrpr.storage.GetReadonly(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly")
		c.AbortWithStatus(500)
		return
	}

	topPie, err := wrpr.storage.GetTopMakerPie(c, ay)
	if err != nil {
		log.Error().Err(err).Msg("Error getting top maker pie")
		c.AbortWithStatus(500)
		return
	}

	ltstpie, ltstpiewhen, err := wrpr.storage.GetLatestPieRanking(c, ay)
	if err != nil {
		log.Error().Err(err).Msg("Error getting latest ranking")
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
			ImgprssrPrefix: wrpr.config.ImgprssrPrefix,
			Year:           ay,
			Maker:          mkr,
		})
	}
	mkrMap := map[string]types.Maker{}
	for _, mkr := range *makers {
		mkrMap[mkr.Id] = mkr
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
			"MakerCards": mrks,
		},
	}
	if topPie != nil {
		pieCategoryLinks := []templater.Link{}
		for _, ct := range topPie.Categories {
			pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", topPie.Year, ct.Slug), Label: ct.Label})
		}
		tp := templater.PieCardData{
			Pie:            *topPie,
			CategoryLinks:  pieCategoryLinks,
			ImgprssrPrefix: wrpr.config.ImgprssrPrefix,
			HasDate:        false,
			PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", topPie.Year, topPie.MakerId, topPie.Id),
			Maker:          mkrMap[topPie.MakerId],
		}

		vals.PageData = map[string]interface{}{
			"MakerCards": mrks,
			"TopPie":     tp,
		}
	}

	if ltstpie != nil && !ro {
		og := vals.PageData.(map[string]interface{})

		pieCategoryLinks := []templater.Link{}
		for _, ct := range ltstpie.Categories {
			pieCategoryLinks = append(pieCategoryLinks, templater.Link{URL: fmt.Sprintf("/years/%d/categories/%s", topPie.Year, ct.Slug), Label: ct.Label})
		}
		lp := templater.PieCardData{
			Pie:            *ltstpie,
			CategoryLinks:  pieCategoryLinks,
			ImgprssrPrefix: wrpr.config.ImgprssrPrefix,
			HasDate:        false,
			PieLink:        fmt.Sprintf("/years/%d/brands/%s/%s", topPie.Year, topPie.MakerId, topPie.Id),
			Maker:          mkrMap[topPie.MakerId],
			RankingTime:    *ltstpiewhen,
		}
		fmt.Println(ltstpiewhen)
		og["LatestRanking"] = lp

		vals.PageData = og
	}

	c.HTML(http.StatusOK, "page:index", vals)
}
