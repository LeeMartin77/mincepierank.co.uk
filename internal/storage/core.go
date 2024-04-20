package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/generated"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

type Operations interface {
	Shutdown()
	TestData(c context.Context) (*string, error)
	Migrate(c context.Context) error
	GetDatabase() *pgxpool.Pool

	GetActiveYear(c context.Context) (*int64, error)

	GetTopMakerPie(c context.Context, activeYear int64) (*MakerPieYearlyWithRankings, error)
	GetFilterableMakerPies(c context.Context, year int64, pageSize int64, zeroIdxPage int64) (*[]MakerPieYearlyWithRankings, error)

	GetMakersForYear(c context.Context, activeYear int64) (*[]types.Maker, error)
	GetMakerPieCategoriesForYear(c context.Context, year int64) (*[]string, error)

	GetMaker(ctx context.Context, Id string) (*types.Maker, error)
	GetMakerPieYearly(ctx context.Context, year int32, makerid string, id string) (*types.MakerPieYearly, error)
	GetMakerPieYearlyRankingSummary(c context.Context, year int32, makerid string, id string) (*RankingSummary, error)
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
