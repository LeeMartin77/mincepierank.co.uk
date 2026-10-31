package website

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templater"
	"github.com/rs/zerolog/log"
)

func (wrpr *WebsiteWrapper) YearBrandPie(w http.ResponseWriter, r *http.Request, year int64, brand string, pieParam string) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
	}

	pie, err := wrpr.storage.GetMakerPieYearly(r.Context(), int32(year), brand, pieParam)
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker year pie")
		abortStatus(w, 500)
		return
	}
	if pie == nil {
		abortStatus(w, 404)
		return
	}

	returnSinglePiePage(wrpr, w, r, year, brand, pieParam, ay, ro, pie, nil)
}

var rankings []string = []string{"filling", "pastry", "topping", "looks", "value"}

var rankingLabels = map[string]string{
	"filling": "Filling",
	"pastry":  "Pastry",
	"topping": "Topping",
	"looks":   "Looks",
	"value":   "Value for Money",
}

// parseRatings reads each ranking field via get (e.g. Request.FormValue) and
// returns the parsed 1-5 values, plus a user-facing message for each field
// that is missing or invalid.
func parseRatings(get func(string) string) (map[string]int32, []string) {
	parsed := map[string]int32{}
	var errs []string
	for _, k := range rankings {
		label := rankingLabels[k]
		rv := get(k)
		if rv == "" {
			errs = append(errs, fmt.Sprintf("%s is required", label))
			continue
		}
		prsd, err := strconv.ParseInt(rv, 10, 32)
		if err != nil || prsd < 1 || prsd > 5 {
			errs = append(errs, fmt.Sprintf("%s must be a whole number from 1 to 5", label))
			continue
		}
		parsed[k] = int32(prsd)
	}
	return parsed, errs
}

// RateYearBrandPie implements generated.ServerInterface.
func (wrpr *WebsiteWrapper) RateYearBrandPie(w http.ResponseWriter, r *http.Request, year int64, brand string, pieParam string) {
	ay, err := wrpr.storage.GetActiveYear(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting active year")
		abortStatus(w, 500)
		return
	}

	ro, err := wrpr.storage.GetReadonly(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly - fails closed")
	}

	pie, err := wrpr.storage.GetMakerPieYearly(r.Context(), int32(year), brand, pieParam)
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker year pie")
		abortStatus(w, 500)
		return
	}
	if pie == nil {
		abortStatus(w, 404)
		return
	}

	if pie.Year != int32(*ay) || ro {
		returnSinglePiePage(wrpr, w, r, year, brand, pieParam, ay, ro, pie, []string{"Rating is not open for this year or is read-only"})
		return
	}

	// rate

	rnkingmap, validationErrs := parseRatings(r.FormValue)
	if len(validationErrs) > 0 {
		returnSinglePiePage(wrpr, w, r, year, brand, pieParam, ay, ro, pie, validationErrs)
		return
	}

	uid := userIDFrom(r)
	wrpr.storage.GetQuerier().UpsertUserMakerPieRanking(r.Context(), sqlcgen.UpsertUserMakerPieRankingParams{
		Year:    int32(year),
		Makerid: brand,
		Pieid:   pieParam,
		Userid:  uid,
		Pastry:  rnkingmap["pastry"],
		Filling: rnkingmap["filling"],
		Topping: rnkingmap["topping"],
		Looks:   rnkingmap["looks"],
		Value:   rnkingmap["value"],
	})

	returnSinglePiePage(wrpr, w, r, year, brand, pieParam, ay, ro, pie, nil)
}

func returnSinglePiePage(wrpr *WebsiteWrapper, w http.ResponseWriter, r *http.Request, year int64, brand string, pieParam string, ay *int64, ro bool, pie *types.MakerPieYearly, errs []string) {
	maker, err := wrpr.storage.GetMaker(r.Context(), brand)
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker year pie")
		abortStatus(w, 500)
		return
	}

	rankingSummary, err := wrpr.storage.GetMakerPieYearlyRankingSummary(r.Context(), int32(year), brand, pieParam)
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker year pie rank summary")
		abortStatus(w, 500)
		return
	}

	cats, err := wrpr.storage.GetMakerPieCategoriesForMakerPieOid(r.Context(), pie.OId)
	if err != nil {
		log.Error().Err(err).Msg("Error getting maker year pie rank categories")
		abortStatus(w, 500)
		return
	}

	uid := userIDFrom(r)
	sin := signedInFrom(r)
	piesEnabled, _ := wrpr.storage.GetUserPiesEnabled(r.Context())

	urnking, _ := wrpr.storage.GetQuerier().GetUserMakerPieRanking(r.Context(), sqlcgen.GetUserMakerPieRankingParams{
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
			Snow:        wrpr.snowing(r.Context()),
			Title:       fmt.Sprintf("%s from %s :: %.2f/5", pie.DisplayName, maker.Name, rankingSummary.Average),
			Description: fmt.Sprintf("%s from %s :: scored  %.2f/5 on average in %d", pie.DisplayName, maker.Name, rankingSummary.Average, pie.Year),
			Keywords:    fmt.Sprintf("Mince Pies, UK, Ranking, %d, %s, %s", pie.Year, maker.Name, pie.DisplayName),
			MenuSettings: templater.MenuSettings{
				ActiveYear:      *ay,
				SignedIn:        sin,
				UserPiesEnabled: piesEnabled,
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
			"Errors":         errs,
			"SignedIn":       sin,
			"SignInRedirect": fmt.Sprintf("/years/%d/brands/%s/%s", pie.Year, pie.MakerId, pie.Id),
		},
	}

	wrpr.htmlTemplater.WritePage(w, http.StatusOK, "singlepie", vals)
}
