package storage

import (
	"context"

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
	GetMakersForYear(c context.Context, activeYear int64) (*[]types.Maker, error)
	GetMakerPieCategoriesForYear(c context.Context, year int64) (*[]string, error)

	GetMaker(ctx context.Context, Id string) (*types.Maker, error)
	GetMakerPieYearly(ctx context.Context, year int32, makerid string, id string) (*types.MakerPieYearly, error)
	GetMakerPieYearlyRankingSummary(c context.Context) (*RankingSummary, error)
}

func (o *OperationWrapper) GetMaker(ctx context.Context, Id string) (*types.Maker, error) {
	return generated.MakerRead(ctx, o.db, Id)
}

func (o *OperationWrapper) GetMakerPieYearly(ctx context.Context, year int32, makerid string, id string) (*types.MakerPieYearly, error) {
	return generated.MakerPieYearlyRead(ctx, o.db, year, makerid, id)
}

func (o *OperationWrapper) GetMakerPieYearlyRankingSummary(c context.Context) (*RankingSummary, error) {
	return &RankingSummary{}, nil
}
