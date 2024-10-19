package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/generated"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

type PieFilters struct {
	BrandIds      []string
	CategorySlugs []string
}

type Operations interface {
	Shutdown()
	TestData(c context.Context) (*string, error)
	Migrate(c context.Context) error
	GetDatabase() *pgxpool.Pool

	GetQuerier() sqlcgen.Querier

	GetActiveYear(c context.Context) (*int64, error)

	GetReadonly(c context.Context) (bool, error)

	GetTopMakerPie(c context.Context, activeYear int64) (*MakerPieYearlyWithRankings, error)
	GetFilterableMakerPies(c context.Context, year int64, pageSize int64, zeroIdxPage int64, filters PieFilters) (*[]MakerPieYearlyWithRankings, error)
	GetFilterableMakerPiesForUser(c context.Context, year int64, userid string, pageSize int64, zeroIdxPage int64, filters PieFilters) (*[]MakerPieYearlyWithRankings, error)

	GetMakersForYear(c context.Context, activeYear int64) (*[]types.Maker, error)
	GetMakerPieCategoriesForYear(c context.Context, year int64, requiredSlugs []string, requiredBrand *string) (*[]types.Category, error)

	GetMaker(ctx context.Context, Id string) (*types.Maker, error)
	GetMakerPieYearly(ctx context.Context, year int32, makerid string, id string) (*types.MakerPieYearly, error)
	GetMakerPieCategoriesForMakerPieOid(ctx context.Context, oid uuid.UUID) (*[]types.Category, error)
	GetMakerPieYearlyRankingSummary(c context.Context, year int32, makerid string, id string) (*RankingSummary, error)

	GetMakerPieCategories(ctx context.Context) (*[]types.Category, error)
	SetPieCategories(ctx context.Context, oid uuid.UUID, category_keys []string) error
}

func (o *OperationWrapper) GetMaker(ctx context.Context, Id string) (*types.Maker, error) {
	return generated.MakerRead(ctx, o.db, Id)
}

func (o *OperationWrapper) GetMakerPieYearly(ctx context.Context, year int32, makerid string, id string) (*types.MakerPieYearly, error) {
	return generated.MakerPieYearlyRead(ctx, o.db, year, makerid, id)
}

func (o *OperationWrapper) GetMakerPieYearlyRankingSummary(c context.Context, year int32, makerid string, id string) (*RankingSummary, error) {
	r := RankingSummary{}
	sql := `
	 SELECT
			SUM(pastry::float )/COUNT(1) as pastry,
	    	SUM(filling::float )/COUNT(1)as filling,
	     SUM(topping::float )/COUNT(1) as topping,
	     SUM(looks::float )/COUNT(1) as looks,
	     SUM(value::float )/COUNT(1) as value,
		(
	     SUM(pastry::float )/COUNT(1) +
	    	SUM(filling::float )/COUNT(1) +
	     SUM(topping::float )/COUNT(1) +
	     SUM(looks::float )/COUNT(1) +
	     SUM(value::float )/COUNT(1)
	   )/5 as avg,
		COUNT(1) as count
	 FROM maker_pie_ranking_yearly
		WHERE year = $1 AND makerid = $2 AND pieid = $3
	 GROUP BY year, makerid, pieid
	`
	res := o.db.QueryRow(c, sql, year, makerid, id)
	err := res.Scan(
		&r.Pastry,
		&r.Filling,
		&r.Topping,
		&r.Looks,
		&r.Value,
		&r.Average,
		&r.Count,
	)
	if err == pgx.ErrNoRows {
		return &r, nil
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}
