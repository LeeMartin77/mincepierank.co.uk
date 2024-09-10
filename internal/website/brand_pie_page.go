package website

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

func (wrpr *WebsiteWrapper) YearBrandPie(ctx *gin.Context, year int64, brand string, pieParam string) {
	ay, err := wrpr.storage.GetActiveYear(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		ctx.AbortWithStatus(500)
		return
	}

	pie, err := wrpr.storage.GetMakerPieYearly(ctx, int32(year), ctx.Param("brand"), ctx.Param("pie"))
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker year pie")
		ctx.AbortWithStatus(500)
		return
	}
	if pie == nil {
		ctx.AbortWithStatus(404)
		return
	}

	maker, err := wrpr.storage.GetMaker(ctx, ctx.Param("brand"))
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker year pie")
		ctx.AbortWithStatus(500)
		return
	}

	rankingSummary, err := wrpr.storage.GetMakerPieYearlyRankingSummary(ctx, int32(year), ctx.Param("brand"), ctx.Param("pie"))
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker year pie rank summary")
		ctx.AbortWithStatus(500)
		return
	}

	cats, err := wrpr.storage.GetMakerPieCategoriesForMakerPieOid(ctx, pie.OId)
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker year pie rank categories")
		ctx.AbortWithStatus(500)
		return
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

	ctx.HTML(http.StatusOK, "page:singlepie", vals)
}
