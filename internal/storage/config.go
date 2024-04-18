package storage

import (
	"context"
	"strconv"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage/generated"
	"github.com/rs/zerolog/log"
)

func (o *OperationWrapper) GetActiveYear(c context.Context) (*int64, error) {
	cfg, err := generated.ConfigRead(c, o.db, "activeYear")

	if err != nil {
		log.Error().Err(err).Msg("Query failed")
		return nil, err
	}
	ay, err := strconv.ParseInt(cfg.Value, 0, 64)
	if err != nil {
		log.Error().Err(err).Msg("Parse failed")
		return nil, err
	}
	return &ay, nil
}
