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
}
