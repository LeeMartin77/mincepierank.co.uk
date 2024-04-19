package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/rs/zerolog/log"
)

func (o *OperationWrapper) GetMakersForYear(c context.Context, activeYear int64) (*[]types.Maker, error) {
	r := []types.Maker{}
	sql := `
	with ordered_brand_ids as (
		select makerid, count(1) as cnt
		  from maker_pie_yearly
		  where year = $1
		  group by makerid
		)
		select distinct mkr.id, mkr.name, mkr.logo, mkr.website, mpy.cnt as pieCount
		from maker mkr
		inner join ordered_brand_ids mpy
				on mpy.makerid = mkr.id
		 order by mpy.cnt desc
	`
	rows, err := o.db.Query(c, sql, activeYear)
	if err == pgx.ErrNoRows {
		return &r, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var itm types.Maker
		var cnt int64
		if err := rows.Scan(&itm.Id, &itm.Name, &itm.Logo, &itm.Website, &cnt); err != nil {
			log.Error().Err(err).Msg("Scan failed in getmakers for year")
			continue
		}
		r = append(r, itm)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error during makers iteration")
	}
	return &r, nil
}
