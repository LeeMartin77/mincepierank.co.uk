package website

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/rs/zerolog/log"
)

// snowing reports whether the public site should show snowfall. The snow is
// the signal that the year is active, so it shows exactly while the site is
// not in read-only mode. On a config read failure it fails closed: no snow.
func (wrpr *WebsiteWrapper) snowing(c context.Context) bool {
	ro, err := wrpr.storage.GetReadonly(c)
	if err != nil {
		log.Error().Err(err).Msg("Error getting readonly for snow state")
		return false
	}
	return !ro
}

func (wrpr *WebsiteWrapper) getMakerMap(c *gin.Context, year int64) (map[string]types.Maker, error) {
	mkrs, err := wrpr.storage.GetMakersForYear(c, year)
	if err != nil {
		return nil, err
	}
	mkrMap := map[string]types.Maker{}
	for _, mkr := range *mkrs {
		mkrMap[mkr.Id] = mkr
	}
	return mkrMap, nil
}
