package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Operations interface {
	Shutdown()
	TestData(c context.Context) (*string, error)
	Migrate(c context.Context) error
	GetDatabase() *pgxpool.Pool

	GetActiveYear(c context.Context) (*int64, error)

	GetTopMakerPie(c context.Context, activeYear int64) (*MakerPieYearlyWithRankings, error)
}
