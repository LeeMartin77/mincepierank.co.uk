package storage

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/rs/zerolog/log"
)

type OperationWrapper struct {
	db      *pgxpool.Pool
	querier sqlcgen.Querier
}

// GetQuerier implements Operations.
func (o *OperationWrapper) GetQuerier() sqlcgen.Querier {
	return o.querier
}

// GetPool implements Storage.
func (o *OperationWrapper) GetTransaction(ctx context.Context) (pgx.Tx, error) {
	return o.db.Begin(ctx)
}

// GetQuerierWithTx implements Storage.
func (pgqc *OperationWrapper) GetQuerierWithTx(tx pgx.Tx) sqlcgen.Querier {
	return sqlcgen.New(tx)
}

// TestData implements Operations.
func (o *OperationWrapper) TestData(c context.Context) (*string, error) {

	var data string
	err := o.db.QueryRow(c, "select test_data from migration_test_table where id=$1", 1).Scan(&data)
	if err != nil {
		log.Error().Err(err).Msg("QueryRow failed")
		return nil, err
	}
	return &data, nil
}

// Shutdown implements Operations.
func (o *OperationWrapper) Shutdown() {
	o.db.Close()
}

func NewOperations(connectionString string) (Operations, error) {
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
		os.Exit(1)
	}
	tmt, clc := context.WithTimeout(context.Background(), time.Second*5)
	defer clc()
	err = pool.Ping(tmt)
	if err != nil {
		log.Error().Err(err).Msg("Unable to ping to database")
		os.Exit(1)
	}

	return &OperationWrapper{
		db:      pool,
		querier: sqlcgen.New(pool),
	}, nil
}

func (o *OperationWrapper) GetDatabase() *pgxpool.Pool {
	return o.db
}
