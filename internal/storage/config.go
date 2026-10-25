package storage

import (
	"context"
	"strconv"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage/generated"
	"github.com/rs/zerolog/log"
)

func (o *OperationWrapper) GetActiveYear(c context.Context) (*int64, error) {
	cfg, err := generated.ConfigRead(c, o.db, "activeYear")

	if err != nil || cfg == nil {
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

// GetReadonly implements Operations.
func (o *OperationWrapper) GetReadonly(c context.Context) (bool, error) {
	cfg, err := generated.ConfigRead(c, o.db, "readonly")
	if err != nil || cfg == nil {
		// Missing key means not read-only — fail closed, not an error.
		return false, nil
	}
	if cfg == nil {
		return false, nil
	}
	val, err := strconv.ParseBool(cfg.Value)
	if err != nil {
		log.Error().Err(err).Msg("Parse failed")
		return false, nil
	}
	return val, nil
}

// GetUserPiesEnabled implements Operations.
func (o *OperationWrapper) GetUserPiesEnabled(c context.Context) (bool, error) {
	cfg, err := generated.ConfigRead(c, o.db, "userPiesEnabled")
	if err != nil || cfg == nil {
		// Missing key means disabled — fail closed, not an error.
		return false, nil
	}
	val, err := strconv.ParseBool(cfg.Value)
	if err != nil {
		// Unparseable value means disabled.
		return false, nil
	}
	return val, nil
}
