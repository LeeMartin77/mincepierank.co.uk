package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/generated"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/rs/zerolog/log"
)

var (
	ErrBrandAlreadyExists = errors.New("a brand with that id already exists")
	ErrBrandNotFound      = errors.New("brand not found")
)

// AdminMakerRow is a brand plus how many brand pies reference it, for the
// admin brands list and the delete guard.
type AdminMakerRow struct {
	types.Maker
	PieCount int64
}

// GetAllMakers returns every brand (not just those with pies this year),
// ordered by name, with a total pie count across all years.
func (o *OperationWrapper) GetAllMakers(c context.Context) (*[]AdminMakerRow, error) {
	r := []AdminMakerRow{}
	sql := `
	select mkr.id, mkr.name, mkr.logo, mkr.website,
	       coalesce(pies.cnt, 0) as piecount
	from maker mkr
	left join (
		select makerid, count(1) as cnt
		from maker_pie_yearly
		group by makerid
	) pies on pies.makerid = mkr.id
	order by lower(mkr.name), mkr.id
	`
	rows, err := o.db.Query(c, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var itm AdminMakerRow
		if err := rows.Scan(&itm.Id, &itm.Name, &itm.Logo, &itm.Website, &itm.PieCount); err != nil {
			log.Error().Err(err).Msg("Scan failed in get all makers")
			continue
		}
		r = append(r, itm)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error during all makers iteration")
	}
	return &r, nil
}

// GetMakerPieCount counts brand pies (all years) referencing the maker id.
// Used to guard brand deletes against orphaning pies.
func (o *OperationWrapper) GetMakerPieCount(c context.Context, makerid string) (int64, error) {
	var cnt int64
	err := o.db.QueryRow(c,
		`select count(1) from maker_pie_yearly where makerid = $1`, makerid).Scan(&cnt)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// CreateMaker inserts a new brand. Returns ErrBrandAlreadyExists when the id
// is taken (the generated insert is ON CONFLICT DO NOTHING, so the existence
// check happens here first).
func (o *OperationWrapper) CreateMaker(c context.Context, m types.Maker) (*types.Maker, error) {
	existing, err := generated.MakerRead(c, o.db, m.Id)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrBrandAlreadyExists
	}
	return generated.MakerCreate(c, o.db, m)
}

// UpdateMaker updates name/logo/website for an existing brand.
// Returns ErrBrandNotFound when the brand does not exist.
func (o *OperationWrapper) UpdateMaker(c context.Context, m types.Maker) (*types.Maker, error) {
	existing, err := generated.MakerRead(c, o.db, m.Id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrBrandNotFound
	}
	return generated.MakerUpdate(c, o.db, m)
}

// DeleteMaker removes a brand. Callers should check GetMakerPieCount first so
// pies are not silently orphaned (brand pages inner join maker).
func (o *OperationWrapper) DeleteMaker(c context.Context, makerid string) error {
	return generated.MakerDelete(c, o.db, makerid)
}

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
