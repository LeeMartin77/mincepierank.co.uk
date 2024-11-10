package website

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
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

	ro, err := wrpr.storage.GetReadonly(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
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

	returnSinglePiePage(wrpr, ctx, year, ay, ro, pie)
}

// RateYearBrandPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) RateYearBrandPie(ctx *gin.Context, year int64, brand string, pieParam string) {
	ay, err := wrpr.storage.GetActiveYear(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		ctx.AbortWithStatus(500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
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

	if pie.Year != int32(*ay) || ro {
		ctx.AbortWithStatus(400)
		return
	}

	returnSinglePiePage(wrpr, ctx, year, ay, ro, pie)
}

func returnSinglePiePage(wrpr *WebsiteWrapper, ctx *gin.Context, year int64, ay *int64, ro bool, pie *types.MakerPieYearly) {
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

	uid, _ := ctx.Keys["userid"].(string)
	sin, _ := ctx.Keys["signedin"].(bool)

	urnking, _ := wrpr.storage.GetQuerier().GetUserMakerPieRanking(ctx, sqlcgen.GetUserMakerPieRankingParams{
		Userid:  uid,
		Year:    pie.Year,
		Makerid: pie.MakerId,
		Pieid:   pie.Id,
	})

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
				SignedIn:   sin,
			},
		},
		PageData: map[string]interface{}{
			"CategoryLinks":  pieCategoryLinks,
			"ImgprssrPrefix": wrpr.config.ImgprssrPrefix,
			"HasMaker":       true,
			"Pie":            pie,
			"RankingSummary": rankingSummary,
			"Maker":          maker,
			"Breadcrumb":     templater.BreadcrumbsFromUrl(fmt.Sprintf("/years/%d/brands/%s/%s", pie.Year, pie.MakerId, pie.Id)),
			"Readonly":       ro,
			"UserRanking":    urnking,
		},
	}

	ctx.HTML(http.StatusOK, "page:singlepie", vals)
}
