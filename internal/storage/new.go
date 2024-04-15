package storage

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type OperationWrapper struct {
	db *pgxpool.Pool
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

	return &OperationWrapper{
		db: pool,
	}, nil
}

func (o *OperationWrapper) GetDatabase() *pgxpool.Pool {
	return o.db
}
