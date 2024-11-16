package website

import (
	"github.com/gin-gonic/gin"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

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
